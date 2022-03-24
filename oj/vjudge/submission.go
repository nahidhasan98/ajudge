package vjudge

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/nlogger"
	"go.mongodb.org/mongo-driver/bson"
)

//Submit function for submitting provlem solution to VJ
func Submit(w http.ResponseWriter, r *http.Request, contestID int, serialIndex string) {
	defer errorhandling.Recovery() //for panic() error Recovery

	//getting form data
	OJ := r.FormValue("OJ")
	pNum := strings.TrimSpace(r.FormValue("pNum"))
	language := r.FormValue("language")
	source := strings.TrimSpace(r.FormValue("source"))

	//for submission first login to VJudge
	if Login() != "success" { //if login unsuccessful
		w.WriteHeader(http.StatusInternalServerError) //status code such as: 400, 404 etc.
		model.Info["StatusCode"] = http.StatusInternalServerError
		model.Tpl.ExecuteTemplate(w, "pageNotFound.gohtml", model.Info)
		return
	}
	//VJudge login success
	logger := nlogger.NewLogger()
	logger.Warn("40: vjudge submission: login done", time.Now())

	//preparing data for POST Request
	postData := url.Values{
		"language": {language},
		"open":     {"0"},
		"source":   {source},
		"captcha":  {""},
		"oj":       {OJ},
		"probNum":  {pNum},
	}

	//submitting to Vjudge
	apiURL := "https://vjudge.net/problem/submit"
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
	errorhandling.Check(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

	response, err := model.Client.Do(req)
	errorhandling.Check(err)
	defer response.Body.Close()
	//subbmission done
	logger.Warn("63: vjudge submission: submission done", time.Now())

	//getting submission ID
	body, err := ioutil.ReadAll(response.Body)
	errorhandling.Check(err)
	type result struct { //json reply gives either error or runID
		RunID int64  `json:"runId"`
		Error string `json:"error"`
	}
	var res result
	json.Unmarshal(body, &res)
	logger.Warn("74: vjudge submission: submission res: "+fmt.Sprintf("%d: ", res.RunID)+res.Error, time.Now())

	if res.Error != "" {
		model.ErrorType = res.Error
		model.PopUpCause = "submissionError"
		//http.Redirect(w, r, model.LastPage, http.StatusSeeOther)
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(res)
		w.Write(b)
		logger.Warn("80: vjudge submission: submission res error", time.Now())
		return
	}
	// else if res.RunID == 0 {
	// 	model.ErrorType = res.Error
	// 	model.PopUpCause = "submissionErrorCustom"
	// 	//http.Redirect(w, r, model.LastPage, http.StatusSeeOther)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	b, _ := json.Marshal("Something went Wrong. Try again.")
	// 	w.Write(b)
	// 	logger.Warn("86: vjudge submission: submission res runid 0", time.Now())
	// 	return
	// }
	//fmt.Println(res.RunID)
	//got submission ID

	//getting language name against value
	language = model.LanguagePack[language]
	//fmt.Println(language)

	//inserting submission records to DB
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	logger.Warn("100: vjudge submission: before DB", time.Now())

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	submissionCollection := DB.Collection("submission")
	counterCollection := DB.Collection("counter")

	//getting data from DB
	var lastSubmissionID model.LastUsedID
	err = counterCollection.FindOne(ctx, bson.M{}).Decode(&lastSubmissionID)
	errorhandling.Check(err)

	currentTime := time.Now().Unix() //this is for DB insertion

	//formating currentTime time to display on frontend
	timeDotTime := time.Unix(currentTime, 0)
	submittedAt := timeDotTime.Format("02-Jan-2006 (15:04:05)")

	//preparing data for inserting to DB
	submissionData := model.SubmissionData{
		SubID:       lastSubmissionID.LastSubmissionID + 1,
		Username:    session.Values["username"].(string),
		OJ:          OJ,
		PNum:        pNum,
		Language:    language,
		SubmittedAt: currentTime,
		VID:         strconv.FormatInt(res.RunID, 10),
		SourceCode:  source,
		Verdict:     "Queueing",
		ContestID:   contestID,
		SerialIndex: serialIndex,
	}
	_, err = submissionCollection.InsertOne(ctx, submissionData)
	errorhandling.Check(err)

	//updating LastSubmissionID to DB for later use/next submission
	updateField := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "lastSubmissionID", Value: 1},
		}},
	}
	_, err = counterCollection.UpdateOne(ctx, bson.M{}, updateField)
	errorhandling.Check(err)

	//preparing data for response back
	respData := struct {
		SubID       int
		OJ          string
		PNum        string
		Language    string
		SourceCode  string
		SubmittedAt string
		ContestID   int
		SerialIndex string
		Error       string `json:"error"`
	}{
		SubID:       submissionData.SubID, //sending submit id to frontend for getting the verdict with ajax call
		OJ:          OJ,
		PNum:        pNum,
		Language:    language,
		SourceCode:  html.EscapeString(source),
		SubmittedAt: submittedAt,
		ContestID:   contestID,
		SerialIndex: serialIndex,
		Error:       "",
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(respData)
	w.Write(b)

	// notofy to discord
	discord := discord.Init()
	discord.SendMessage(submissionData, "submission")
}
