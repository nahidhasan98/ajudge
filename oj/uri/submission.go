package uri

import (
	"encoding/json"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
)

//Submit function for submitting provlem solution to URI
func Submit(w http.ResponseWriter, r *http.Request, contestID int, serialIndex string) {
	defer errorhandling.Recovery() //for panic() error Recovery

	//gettimng form data
	pNum := strings.TrimSpace(r.FormValue("pNum"))
	language := r.FormValue("language")
	source := strings.TrimSpace(r.FormValue("source"))

	//preparing source code for later comparing (ignoring carriage return /r) (for verdict)
	sourceEscape := html.EscapeString(source)
	var sourceMod string

	for i := 0; i < len(sourceEscape); i++ {
		if sourceEscape[i] != 13 { //ignoring (carriage return-/r) for comparing
			sourceMod += string(sourceEscape[i])
		}
	}

	//for submission first login to URI
	if Login() != "success" { //if login unsuccessful
		w.WriteHeader(http.StatusInternalServerError) //status code such as: 400, 404 etc.
		model.Info["StatusCode"] = http.StatusInternalServerError
		model.Tpl.ExecuteTemplate(w, "pageNotFound.gohtml", model.Info)
		return
	}
	//URI login success

	//preparing data for POST Request
	postData := url.Values{
		"_method":          {method},
		"_csrfToken":       {csrfToken},
		"problem_id":       {pNum},
		"language_id":      {language},
		"template":         {"1"},
		"source_code":      {source},
		"_Token[fields]":   {tokenFields},
		"_Token[unlocked]": {tokenUnlocked},
	}

	//submitting to URI
	apiURL := "https://www.beecrowd.com.br/judge/en/runs/add"
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

	//setting up requset to prevent auto redirect
	model.Client = &http.Client{
		Jar: model.CookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := model.Client.Do(req)
	errorhandling.Check(err)
	defer response.Body.Close()
	//subbmission done

	//getting submission ID
	apiURL = "https://www.beecrowd.com.br/judge/en/runs?problem_id=" + pNum + "&language_id=" + language
	response = GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	//getting all submission ID from latest submission page
	var subIDList []string
	var tempID string
	document.Find("td[class='id']").Each(func(index int, mixedStr *goquery.Selection) {
		tempID = mixedStr.Find("a").Text()
		subIDList = append(subIDList, tempID)
	})

	//getting submission ID by matching original source code & submitted source code
	var actualSubID string
	for i := 0; i < len(subIDList); i++ {
		//getting submitted code one by one with collected subIDList
		apiURL := "https://www.beecrowd.com.br/judge/en/runs/code/" + subIDList[i]
		response := GETRequest(apiURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		tempSubCode, _ := document.Find("pre[id='code']").Html()
		tempSubCode = strings.TrimSpace(tempSubCode)

		if tempSubCode == sourceMod { //if original source code & submitted source code matched
			actualSubID = subIDList[i] //got actual/this submission ID
			break
		}
	}
	//got submission ID

	//getting language name against value
	language = model.LanguagePack[language]
	index := strings.Index(language, "(") //taking only language name for storing in DB, avoiding comliler name
	if index != -1 {                      //C++17 (g++ 7.3.0, -std=c++17 -O2 -lm) [+0s for timelimit]
		language = language[:index-1] //C++17 (taking only this)
	}
	//fmt.Println(language)

	//inserting submission records to DB
	session, _ := model.Store.Get(r, "mysession")

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
		OJ:          "URI",
		PNum:        pNum,
		Language:    language,
		SubmittedAt: currentTime,
		VID:         actualSubID,
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
		Error       string `json:"error"` //for vj submit error
	}{
		SubID:       submissionData.SubID, //sending submit id to frontend for getting the verdict with ajax call
		OJ:          "URI",
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
