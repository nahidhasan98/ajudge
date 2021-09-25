package toph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
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

//Submit function for submitting provlem solution to Toph
func Submit(w http.ResponseWriter, r *http.Request, contestID int, serialIndex string) {
	defer errorhandling.Recovery() //for panic() error Recovery

	//getting form data
	pNum := strings.TrimSpace(r.FormValue("pNum"))
	language := r.FormValue("language")
	source := strings.TrimSpace(r.FormValue("source"))

	//for submission first login to Toph
	if Login() != "success" { //if login unsuccessful
		w.WriteHeader(http.StatusInternalServerError) //status code such as: 400, 404 etc.
		model.Info["StatusCode"] = http.StatusInternalServerError
		model.Tpl.ExecuteTemplate(w, "pageNotFound.gohtml", model.Info)
		return
	}
	//Toph login success

	//getting submission link
	apiURL := "https://toph.co/p/" + pNum
	response := GETRequest(apiURL)
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	probID, _ := document.Find("aside").Attr("id")
	probID = probID[9:]
	//fmt.Println(probID)
	apiURL = "https://toph.co/api/problems/" + probID + "/submissions"
	//got submission link

	//submitting to Toph
	//preparing data for POST Request
	postData, contentType, err := createMultipart(language, source) //Toph receives multipart/form-data for code submission
	errorhandling.Check(err)

	req, err := http.NewRequest("POST", apiURL, postData)
	errorhandling.Check(err)
	req.Header.Add("authorization", tokenID)
	req.Header.Set("Content-Type", contentType)

	response, err = model.Client.Do(req)
	errorhandling.Check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	errorhandling.Check(err)
	//submission done

	//now getting the submission ID
	type subResponse struct {
		SubID int64 `json:"prettyId"`
	}
	var subRes subResponse
	json.Unmarshal(body, &subRes)
	//fmt.Println(subRes.SubID)
	//got submission ID

	//getting language name against lang-value for inserting to DB
	language = model.LanguagePack[language]
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
		OJ:          "Toph",
		PNum:        pNum,
		Language:    language,
		SubmittedAt: currentTime,
		VID:         strconv.FormatInt(subRes.SubID, 10),
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
		OJ:          "Toph",
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
	discord.SendMessage(submissionData)
}

func createMultipart(language, source string) (io.Reader, string, error) {
	//defining variable for returning data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// without filename parameter (Toph form-data input field, name="languageId")
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="languageId"`)
	fileWrite, err := w.CreatePart(header)
	errorhandling.Check(err)
	_, err = io.Copy(fileWrite, strings.NewReader(language))
	errorhandling.Check(err)

	// with filename (Toph form-data input field, name="source"; filename="solution.cpp")
	header = make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="source"; filename="`+getSubmitableFilenameAndExtention(language)+`"`))
	header.Set("Content-Type", "application/octet-stream")
	fileWrite, err = w.CreatePart(header)
	errorhandling.Check(err)
	_, err = io.Copy(fileWrite, strings.NewReader(source))
	errorhandling.Check(err)

	// Close the multipart writer so that the request conatians the terminating boundary.
	w.Close()

	return &b, w.FormDataContentType(), nil
}

func getSubmitableFilenameAndExtention(language string) string {
	ext := ""

	if language == "5d8289da9d55050001e97eae" {
		ext = "solution.sh"
	} else if language == "5d8211eb728b11000151faf5" {
		ext = "solution.bf"
	} else if language == "5d8280551335cb000138ba63" {
		ext = "solution.cs"
	} else if language == "5d84f038f10beb00010af77c" {
		ext = "solution.cpp"
	} else if language == "5d84ef3ef10beb00010af742" {
		ext = "solution.cpp"
	} else if language == "5d828f1e9d55050001e97ee4" {
		ext = "solution.cpp"
	} else if language == "5d832d5e1335cb000138bd1f" {
		ext = "solution.c"
	} else if language == "5eae739b36ac0000016688bd" {
		ext = "solution.lisp"
	} else if language == "5e85f3d7e2613b000165fc33" {
		ext = "solution.erl"
	} else if language == "5d821a8bf2eba50001686581" {
		ext = "solution.pas"
	} else if language == "55c9ab8c421aa961d1000007" {
		ext = "solution.go"
	} else if language == "5d832a4f1335cb000138bd08" {
		ext = "solution.hs"
	} else if language == "58483d7504469e2585024395" {
		ext = "Solution.java"
	} else if language == "59ca12114ad24000017dcaf9" {
		ext = "solution.kt"
	} else if language == "5d832ccd1335cb000138bd19" {
		ext = "solution.js"
	} else if language == "5d822463728b11000151fbc9" {
		ext = "solution.pl"
	} else if language == "5d8334c31335cb000138bd5c" {
		ext = "solution.php"
	} else if language == "55c9a240421aa9479c000010" {
		ext = "solution.py"
	} else if language == "55c9a6a6421aa961d1000003" {
		ext = "solution.py"
	} else if language == "58482b5504469e2585024320" {
		ext = "solution.py"
	} else if language == "58482c1804469e2585024324" {
		ext = "solution.py"
	} else if language == "5f4793f146e836000119165f" {
		ext = "solution.py"
	} else if language == "5848505704469e258502445b" {
		ext = "solution.rb"
	} else if language == "5f6aeccae1863f0001267dc8" {
		ext = "main.swift"
	} else if language == "5eeb7c1c67d6530001de8e47" {
		ext = "solution.ws"
	}

	return ext
}
