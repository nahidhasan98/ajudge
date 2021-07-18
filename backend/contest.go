package backend

import (
	"encoding/json"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/oj/dimik"
	"github.com/nahidhasan98/ajudge/oj/toph"
	"github.com/nahidhasan98/ajudge/oj/uri"
	"github.com/nahidhasan98/ajudge/oj/vjudge"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Contest function for retrieving contest list
func Contest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	model.Info["Username"] = session.Values["username"]
	model.Info["Password"] = session.Values["password"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "ContestList"
	model.Info["PageTitle"] = "Contest List | AJudge"
	model.Info["LastPage"] = model.LastPage
	model.Info["PopUpCause"] = model.PopUpCause

	model.Tpl.ExecuteTemplate(w, "contestList.gohtml", model.Info)

	//clearing up values (because it may be used in wrong place unintentionally)
	model.PopUpCause = ""
	model.Info["PopUpCause"] = model.PopUpCause
}

//GetContestList function for retrieving contest list
func GetContestList(w http.ResponseWriter, r *http.Request) {
	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	contestCollection := DB.Collection("contest")

	//getting data for this user from DB
	var contestList []model.ContestData

	//setting uo options for retrieving data from DB
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "contestID", Value: -1}}) //sorting by contestID
	cursor, err := contestCollection.Find(ctx, bson.D{}, opts)
	errorhandling.Check(err)

	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.ContestData
		err := cursor.Decode(&temp)
		errorhandling.Check(err)

		contestList = append(contestList, temp)
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(contestList)
	w.Write(b)
}

//CreateContest function for creating new contest
func CreateContest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := model.Store.Get(r, "mysession")
	model.LastPage = r.URL.Path

	if r.Method != "POST" {
		if session.Values["isLogin"] == true {
			if model.IsAccVerifed(r) {
				model.Info["Username"] = session.Values["username"]
				model.Info["Password"] = session.Values["password"]
				model.Info["IsLogged"] = session.Values["isLogin"]
				model.Info["PageName"] = "CreateContest"
				model.Info["PageTitle"] = "Create Contest | AJudge"

				model.Tpl.ExecuteTemplate(w, "createContest.gohtml", model.Info)
			} else {
				model.PopUpCause = "verifyRequired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		} else {
			model.PopUpCause = "loginRequired"
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	} else if r.Method == "POST" {
		//getting form data
		contestTitle := r.FormValue("contestTitle")
		contestDate := r.FormValue("contestDate") //this format: 2020-12-31
		contestTime := r.FormValue("contestTime") //this format: 23:59
		contestDuration := r.FormValue("contestDuration")
		clientTZOffset, _ := strconv.Atoi(r.FormValue("timeZoneOffset"))

		//currentTime := time.Now()
		//zone, _ := currentTime.Zone()
		contestDateTime := contestDate + "T" + contestTime + "Z"
		contestDT, err := time.Parse(time.RFC3339, contestDateTime) //format RFC3339 = "2006-01-02T15:04:05Z07:00"
		errorhandling.Check(err)

		contestStartAt := contestDT.Unix() //converting start time to unix for storing to DB
		contestStartAt += int64(clientTZOffset)

		var probSetData []model.ProblemSet

		for i := 0; i < 26; i++ { //max 26 problem
			index := strconv.Itoa(i + 65)

			if r.FormValue("OJ"+index) == "" {
				break
			} else {
				var temp model.ProblemSet

				temp.SerialIndex = string(rune(i + 65)) //converting ascii to char: 65 -> A
				temp.OJ = r.FormValue("OJ" + index)
				temp.PNum = r.FormValue("pNum" + index)
				temp.PName = r.FormValue("pName" + index)
				temp.CustomName = r.FormValue("customName" + index)

				probSetData = append(probSetData, temp)
			}
		}

		//connecting to DB
		DB, ctx, cancel := db.Connect()
		defer cancel()
		defer DB.Client().Disconnect(ctx)

		//taking DB collection/table to a variable
		contestCollection := DB.Collection("contest")
		counterCollection := DB.Collection("counter")

		//getting data from DB
		var dbQuery model.LastUsedID
		err = counterCollection.FindOne(ctx, bson.M{}).Decode(&dbQuery)
		errorhandling.Check(err)

		//preparing data for inserting to DB
		contestData := model.ContestData{
			ContestID:  dbQuery.LastContestID + 1,
			Title:      contestTitle,
			Date:       contestDate,
			Time:       contestTime,
			StartAt:    contestStartAt,
			Duration:   contestDuration,
			Author:     session.Values["username"].(string),
			ProblemSet: probSetData,
		}
		_, err = contestCollection.InsertOne(ctx, contestData) //inserting contest details to contest table
		errorhandling.Check(err)

		//updating lastContestID to DB for later use/next contest
		updateField := bson.D{
			{Key: "$inc", Value: bson.D{
				{Key: "lastContestID", Value: 1},
			}},
		}
		_, err = counterCollection.UpdateOne(ctx, bson.M{}, updateField)
		errorhandling.Check(err)

		model.PopUpCause = "contestCreated"
		http.Redirect(w, r, "/contest", http.StatusSeeOther)
		return
	}
}

//ContestUpadte function for updating contest data if author want to
func ContestUpadte(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := model.Store.Get(r, "mysession")
	model.LastPage = r.URL.Path

	path := r.URL.Path
	conID := strings.TrimPrefix(path, "/contestUpdate/")
	contestID, _ := strconv.Atoi(conID)

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	contestCollection := DB.Collection("contest")

	//getting data from DB
	var dbQuery model.ContestData
	res := contestCollection.FindOne(ctx, bson.M{"contestID": contestID}).Decode(&dbQuery)
	if res == mongo.ErrNoDocuments {
		errorPage(w, http.StatusBadRequest)
		return
	}

	if dbQuery.Author != session.Values["username"] { //update can be done by author
		errorPage(w, http.StatusBadRequest)
		return
	}
	//taking care of time
	timeHHText := dbQuery.Time[0:2]
	timeMMText := dbQuery.Time[3:5]
	timeHH, _ := strconv.Atoi(timeHHText)

	if timeHH == 00 {
		dbQuery.Time = "12:" + timeMMText + " AM"
	} else if timeHH >= 1 && timeHH <= 11 {
		dbQuery.Time += " AM"
	} else if timeHH == 12 {
		dbQuery.Time += " PM"
	} else if timeHH >= 13 && timeHH <= 23 {
		dbQuery.Time = strconv.Itoa(timeHH-12) + ":" + timeMMText + " PM"
	}

	if r.Method != "POST" {
		if session.Values["isLogin"] == true {
			if model.IsAccVerifed(r) {
				model.Info["Username"] = session.Values["username"]
				model.Info["Password"] = session.Values["password"]
				model.Info["IsLogged"] = session.Values["isLogin"]
				model.Info["PageName"] = "UpdateContest"
				model.Info["PageTitle"] = "Update Contest | AJudge"
				model.Info["ContestData"] = dbQuery
				model.Info["Addition"] = model.Addition{}

				model.Tpl.ExecuteTemplate(w, "createContest.gohtml", model.Info)
			} else {
				model.PopUpCause = "verifyRequired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		} else {
			model.PopUpCause = "loginRequired"
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	} else if r.Method == "POST" {
		//getting form data
		contestTitle := r.FormValue("contestTitle")
		contestDuration := r.FormValue("contestDuration")

		var probSetData []model.ProblemSet

		for i := 0; i < 26; i++ { //max 26 problem
			index := strconv.Itoa(i + 65)

			if r.FormValue("OJ"+index) == "" {
				break
			} else {
				var temp model.ProblemSet

				temp.SerialIndex = string(rune(i + 65)) //converting ascii to char: 65 -> A
				temp.OJ = r.FormValue("OJ" + index)
				temp.PNum = r.FormValue("pNum" + index)
				temp.PName = r.FormValue("pName" + index)
				temp.CustomName = r.FormValue("customName" + index)

				probSetData = append(probSetData, temp)
			}
		}

		//updating new data to DB
		updateField := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "title", Value: contestTitle},
				{Key: "duration", Value: contestDuration},
				{Key: "problemSet", Value: probSetData},
			}},
		}
		_, err := contestCollection.UpdateOne(ctx, bson.M{"contestID": contestID}, updateField)
		errorhandling.Check(err)

		model.PopUpCause = "contestUpdated"
		http.Redirect(w, r, "/contest", http.StatusSeeOther)
		return
	}
}

//ContestGround function for contest arena
func ContestGround(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	path := r.URL.Path
	conID := strings.TrimPrefix(path, "/contest/")
	contestID, _ := strconv.Atoi(conID)

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	contestCollection := DB.Collection("contest")

	//getting data from DB
	var dbQuery model.ContestData

	res := contestCollection.FindOne(ctx, bson.M{"contestID": contestID}).Decode(&dbQuery)

	if res == mongo.ErrNoDocuments {
		errorPage(w, http.StatusBadRequest)
		return
	}

	var runningStatus string

	cMM := dbQuery.Duration[len(dbQuery.Duration)-2:]
	cHH := ""
	for i := 0; i < len(dbQuery.Duration); i++ {
		if string(dbQuery.Duration[i]) == ":" {
			break
		} else {
			cHH += string(dbQuery.Duration[i])
		}
	}
	cMin, _ := strconv.Atoi(cMM)
	cHour, _ := strconv.Atoi(cHH)

	cDuration := (cHour * 60 * 60) + (cMin * 60)
	currentTime := time.Now().Unix()

	if currentTime < dbQuery.StartAt {
		runningStatus = "BeforeContest"
	} else if currentTime >= dbQuery.StartAt && currentTime < dbQuery.StartAt+int64(cDuration) {
		runningStatus = "RunningContest"
	} else {
		runningStatus = "AfterContest"
	}

	if session.Values["username"] == nil {
		model.Info["Username"] = ""
	} else {
		model.Info["Username"] = session.Values["username"]
	}
	model.Info["Password"] = session.Values["password"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "ContestGround"
	model.Info["PageTitle"] = "Contest Ground | AJudge"
	model.Info["LastPage"] = model.LastPage
	model.Info["PopUpCause"] = model.PopUpCause
	model.Info["ContestData"] = dbQuery
	model.Info["RunningStatus"] = runningStatus

	model.Tpl.ExecuteTemplate(w, "contestGround.gohtml", model.Info)
}

//GetProblemSet funtion for getting a single problem description of a contest
func GetProblemSet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	pathSeg := strings.TrimPrefix(path, "/problemSet/") //url like this: /problemSet/1/A/Toph-coycat
	segments := strings.Split(pathSeg, "/")

	contestID := segments[0]
	serialIndex := segments[1]
	OJpNum := segments[2]

	need := "-"
	index := strings.Index(OJpNum, need)

	if index == -1 { //url is not like this "/problemview/OJ-pNum"
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
		return
	}

	OJ, pNum := "", ""
	if len(OJpNum) >= 4 {
		OJ = OJpNum[0:3]    //For chinese OJ, it's 3 -> 计蒜客
		pNum = OJpNum[3+1:] //For chinese OJ
	}

	if OJ != "计蒜客" { //Other than Chinese OJ, In normal situation
		//checking for 2nd chinese OJ
		if len(OJpNum) >= 5 {
			OJ = OJpNum[0:4]    //For chinese OJ, it's 4 -> 黑暗爆炸
			pNum = OJpNum[4+1:] //For chinese OJ

			if OJ != "黑暗爆炸" { //Other than Chinese OJ, In normal situation
				OJ = OJpNum[0:index]
				pNum = OJpNum[index+1:]
			}
		}
	}

	if !model.OJSet[OJ] || pNum == "" { //bad url, not OJ & pNum specified
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
		return
	}
	// got something in OJ and pNum
	allowSubmit := false  //just for declaring, need later
	var uvaSegment string //declaring, will be used later

	DimikOJProblem := map[string]interface{}{}
	TophProblem := map[string]interface{}{}
	URIProblem := map[string]interface{}{}
	VJProblem := map[string]interface{}{}

	if OJ == "DimikOJ" {
		DimikOJProblem["Des"] = template.HTML(dimik.ProbDes(pNum))
		allowSubmit = true

		if model.PTitle == "" {
			errorPage(w, http.StatusBadRequest)
			return
		}
	} else if OJ == "Toph" {
		TophProblem["Des"] = template.HTML(toph.ProbDes(pNum))
		allowSubmit = true

		if model.PTitle == "" {
			errorPage(w, http.StatusBadRequest)
			return
		}
	} else if OJ == "URI" {
		URIProblem["Des"] = template.HTML(uri.ProbDes(pNum))
		allowSubmit = true

		if model.PTitle == "" {
			errorPage(w, http.StatusBadRequest)
			return
		}
	} else {
		var tempDes string
		var status int
		tempDes, allowSubmit, status = vjudge.ProbDes(OJ, pNum)
		VJProblem["Des"] = template.HTML(tempDes)

		if model.PTitle == "" {
			errorPage(w, http.StatusBadRequest)
			return
		}
		//checking whether problem submission allowed or not
		if allowSubmit && status == 0 {
			allowSubmit = true
		}

		//for UVA pdf description
		if OJ == "UVA" {
			temp, _ := strconv.Atoi(pNum)
			IntSegment := temp / 100
			uvaSegment = strconv.Itoa(IntSegment)
		}
	}

	model.LastPage = "/contest"
	session, _ := model.Store.Get(r, "mysession")

	model.Info["Username"] = session.Values["username"]
	model.Info["Password"] = session.Values["password"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "ProblemSet"
	model.Info["PageTitle"] = model.PTitle + " | AJudge"
	model.Info["OJ"] = OJ
	model.Info["PNum"] = pNum
	model.Info["AllowSubmit"] = allowSubmit
	model.Info["UvaSegment"] = uvaSegment
	model.Info["PName"] = model.PTitle
	model.Info["TimeLimit"] = model.PTimeLimit
	model.Info["MemoryLimit"] = model.PMemoryLimit
	model.Info["SourceLimit"] = model.PSourceLimit
	model.Info["DimikOJProblem"] = DimikOJProblem
	model.Info["TophProblem"] = TophProblem
	model.Info["URIProblem"] = URIProblem
	model.Info["VJProblem"] = VJProblem
	model.Info["ContestID"] = contestID
	model.Info["SerialIndex"] = serialIndex

	model.Tpl.ExecuteTemplate(w, "problemSet.gohtml", model.Info)
}

//SubmitC function for submitting a problem solution for contest
func SubmitC(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		contestIDSerialIndex := r.FormValue("conIDSerial")
		mixedStr := strings.Split(contestIDSerialIndex, " - ") //string is like: "21 - A"

		//url is something like this "/submit/OJ-pNum"
		contestID, _ := strconv.Atoi(mixedStr[0])
		serialIndex := mixedStr[1]

		//connecting to DB
		DB, ctx, cancel := db.Connect()
		defer cancel()
		defer DB.Client().Disconnect(ctx)

		//taking DB collection/table to a variable
		contestCollection := DB.Collection("contest")

		var dbQuery model.ContestData

		err := contestCollection.FindOne(ctx, bson.M{"contestID": contestID}).Decode(&dbQuery)
		if err == mongo.ErrNoDocuments { //found no rows (username available)
			errorPage(w, http.StatusInternalServerError) //http.StatusBadRequest = 400
			return
		}

		cMM := dbQuery.Duration[len(dbQuery.Duration)-2:]
		cHH := ""
		for i := 0; i < len(dbQuery.Duration); i++ {
			if string(dbQuery.Duration[i]) == ":" {
				break
			} else {
				cHH += string(dbQuery.Duration[i])
			}
		}
		cMin, _ := strconv.Atoi(cMM)
		cHour, _ := strconv.Atoi(cHH)

		cDuration := (cHour * 60 * 60) + (cMin * 60)
		currentTime := time.Now().Unix()

		if currentTime >= dbQuery.StartAt+int64(cDuration) { //if contest end, submission details won't be stored for contest
			contestID = 0
			serialIndex = ""
		}

		//getting form data
		OJ := r.FormValue("OJ")

		if OJ == "DimikOJ" {
			pNum := strings.TrimSpace(r.FormValue("pNum"))

			//checking if the pNum is a int or not
			pNumInt, _ := strconv.Atoi(pNum) //if pNum contains only digit, it will remain same("12"->12), otherwise become 0("12abc"->0)
			pNum = strconv.Itoa(pNumInt)

			//checking for problem exist or not
			apiURL := "https://dimikoj.com/problems/" + pNum
			response := dimik.GETRequest(apiURL)
			defer response.Body.Close()
			document, err := goquery.NewDocumentFromReader(response.Body)
			errorhandling.Check(err)

			title := document.Find("h2[class='card-title']").Text()

			if title == "" { //no such problem
				model.PopUpCause = "NoSuchProblem"
				model.Info["PopUpCause"] = model.PopUpCause
				http.Redirect(w, r, "/submit", http.StatusSeeOther)
				return
			}
			dimik.Submit(w, r, contestID, serialIndex)
			return
		} else if OJ == "Toph" {
			//checking for problem exist or not
			response := toph.GETRequest("https://toph.co/p/" + strings.TrimSpace(r.FormValue("pNum")))
			defer response.Body.Close()
			document, err := goquery.NewDocumentFromReader(response.Body)
			errorhandling.Check(err)

			title := document.Find("span[class='artifact__caption']").Find("h1").Text()

			if title == "" { //no such problem
				model.PopUpCause = "NoSuchProblem"
				model.Info["PopUpCause"] = model.PopUpCause
				http.Redirect(w, r, "/submit", http.StatusSeeOther)
				return
			}
			toph.Submit(w, r, contestID, serialIndex)
			return
		} else if OJ == "URI" {
			pNum := strings.TrimSpace(r.FormValue("pNum"))

			//checking if the pNum is a int or not
			pNumInt, _ := strconv.Atoi(pNum) //if pNum contains only digit, it will remain same("12"->12), otherwise become 0("12abc"->0)
			pNum = strconv.Itoa(pNumInt)

			//checking for problem exist or not
			apiURL := "https://www.urionlinejudge.com.br/judge/en/problems/view/" + pNum
			response := uri.GETRequest(apiURL)
			defer response.Body.Close()

			document, err := goquery.NewDocumentFromReader(response.Body)
			errorhandling.Check(err)
			pDesSrcURI, _ := document.Find("iframe").Attr("src") //this is prob description source

			if pDesSrcURI == "" { //no such problem
				model.PopUpCause = "NoSuchProblem"
				model.Info["PopUpCause"] = model.PopUpCause
				http.Redirect(w, r, "/submit", http.StatusSeeOther)
				return
			}
			uri.Submit(w, r, contestID, serialIndex)
			return
		} else {
			//checking for problem exist or not
			tempOJ := r.FormValue("OJ")
			tempPNum := strings.TrimSpace(r.FormValue("pNum"))
			if tempOJ == "计蒜客" || tempOJ == "黑暗爆炸" {
				tempOJ = url.QueryEscape(tempOJ)
			}

			apiURL := "https://vjudge.net/problem/" + tempOJ + "-" + tempPNum
			response := vjudge.GETRequest(apiURL)
			defer response.Body.Close()
			document, err := goquery.NewDocumentFromReader(response.Body)
			errorhandling.Check(err)
			title := document.Find("div[id='prob-title']").Find("h2").Text()

			if title == "" { //no such problem
				model.PopUpCause = "NoSuchProblem"
				model.Info["PopUpCause"] = model.PopUpCause
				http.Redirect(w, r, "/submit", http.StatusSeeOther)
				return
			}
			vjudge.Submit(w, r, contestID, serialIndex)
			return
		}
	}
}

//GetContestData function for getting contest data link submissions information
func GetContestData(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	contestIDStr := strings.TrimPrefix(path, "/dataContest/")
	contestID, _ := strconv.Atoi(contestIDStr)

	session, _ := model.Store.Get(r, "mysession")

	type setData struct {
		Username    string
		SerialIndex string
	}
	setSolvedStatus := make(map[string]bool)     // New empty set
	setAttempedStatus := make(map[string]bool)   // New empty set
	setTotalSolved := make(map[setData]bool)     // New empty set
	setTotalSubmission := make(map[setData]bool) // New empty set
	setContestant := make(map[string]bool)       // New empty set

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	submissionCollection := DB.Collection("submission")
	contestCollection := DB.Collection("contest")

	//taking contest data/start time from DB
	var dbQuery2 model.ContestData

	res := contestCollection.FindOne(ctx, bson.M{"contestID": contestID}).Decode(&dbQuery2)
	if res == mongo.ErrNoDocuments {
		errorPage(w, http.StatusBadRequest)
	}

	//preparing data for retrieving to DB
	var cSubmissionList []model.SubmissionData

	type subDetails struct {
		SerialIndex      string
		Penalty          int
		CompilationError int
		Verdict          string
		AcceptedAt       int64
	}
	type contestant struct {
		Username    string
		TotalSolved int
		TotalTime   int64
		SubDetails  map[string]subDetails
	}
	var contestantData []contestant

	//setting uo options for retrieving data from DB
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "submittedAt", Value: -1}}) //descending order
	cursor, err := submissionCollection.Find(ctx, bson.M{"contestID": contestID}, opts)
	errorhandling.Check(err)

	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.SubmissionData
		err := cursor.Decode(&temp)
		errorhandling.Check(err)
		temp.SourceCode = html.EscapeString(temp.SourceCode) //specially for reserving newline

		//checking if this user is appropriate to see the source code or not
		//source code will be provided to the correct owner and the contest author
		if (temp.Username != session.Values["username"]) && (dbQuery2.Author != session.Values["username"]) {
			temp.SourceCode = ""
		}

		//no need to send OJ & pNum to frontend of a contest
		temp.OJ, temp.PNum = "", ""

		//for submission list
		cSubmissionList = append(cSubmissionList, temp)

		//for dashboard solved status for logged in user
		if session.Values["username"] != nil {
			currentUser := session.Values["username"].(string)
			if currentUser == temp.Username && temp.Verdict == "Accepted" {
				setSolvedStatus[temp.SerialIndex] = true //adding to set
			} else if currentUser == temp.Username && temp.Verdict != "Accepted" {
				setAttempedStatus[temp.SerialIndex] = true
			}
		}

		//for overall solved/submission - 1st part
		var temp2 setData
		temp2.Username = temp.Username
		temp2.SerialIndex = temp.SerialIndex
		if temp.Verdict == "Accepted" {
			setTotalSolved[temp2] = true //adding to set
		}
		setTotalSubmission[temp2] = true

		//for standings section - 1st part
		tempUser := temp.Username
		setContestant[tempUser] = true
	}
	//for overall solved/submission - 2nd part
	var totalSolved = make(map[string]int)
	var totalSubmission = make(map[string]int)
	for key := range setTotalSolved {
		totalSolved[key.SerialIndex]++
	}
	for key := range setTotalSubmission {
		totalSubmission[key.SerialIndex]++
	}

	//for standings section - 2nd part
	for key := range setContestant {
		var temp contestant
		temp.Username = key
		contestantData = append(contestantData, temp)
	}
	for i := 0; i < len(contestantData); i++ {
		contestantData[i].SubDetails = make(map[string]subDetails)
		var problemSet = make(map[string]subDetails)

		//all submitted problem info of a user - will be used in next loop
		for j := 65; j < 91; j++ { //resetting problem info for a new user
			index := string(rune(j)) //converting 65 to "A"

			var temp subDetails
			temp.SerialIndex = index
			temp.Penalty = 0
			temp.CompilationError = 0
			temp.Verdict = ""
			temp.AcceptedAt = 0

			problemSet[index] = temp
		}
		for j := len(cSubmissionList) - 1; j >= 0; j-- { //reverse loop for getting first submission first
			if contestantData[i].Username == cSubmissionList[j].Username {
				tempSerial := cSubmissionList[j].SerialIndex
				tempVerdict := cSubmissionList[j].Verdict

				if problemSet[tempSerial].Verdict != "Accepted" { //for AC submission, only count the first accepted submission
					if cSubmissionList[j].TerminalVerdict { //if got final verdict - not in queue/judging
						if tempVerdict == "Accepted" {
							//setting up temp problemSet
							var temp subDetails
							temp.SerialIndex = tempSerial
							temp.Penalty = problemSet[tempSerial].Penalty
							temp.CompilationError = problemSet[tempSerial].CompilationError
							temp.Verdict = cSubmissionList[j].Verdict        //verdict will be set to "Accepted"
							temp.AcceptedAt = cSubmissionList[j].SubmittedAt //submitted time will be set
							problemSet[tempSerial] = temp

							//setting up original variable contestantData
							contestantData[i].SubDetails[tempSerial] = temp //totalSolved & time will be set
							contestantData[i].TotalSolved++
							contestantData[i].TotalTime += currentElapsedTime(cSubmissionList[j].SubmittedAt, dbQuery2.StartAt) + int64(problemSet[tempSerial].Penalty*(20*60))
						} else if tempVerdict == "Compilation Error" || tempVerdict == "Compilation error" || tempVerdict == "Compile Error" || tempVerdict == "Compile error" {
							//setting up temp problemSet
							var temp subDetails
							temp.SerialIndex = tempSerial
							temp.Penalty = problemSet[tempSerial].Penalty
							temp.CompilationError = problemSet[tempSerial].CompilationError + 1 //only CompilationError counter will inc by 1
							temp.Verdict = cSubmissionList[j].Verdict
							temp.AcceptedAt = cSubmissionList[j].SubmittedAt
							problemSet[tempSerial] = temp

							//setting up original variable contestantData
							contestantData[i].SubDetails[tempSerial] = temp
						} else { //if wrong answer or vice versa
							//setting up temp problemSet
							var temp subDetails
							temp.SerialIndex = tempSerial
							temp.Penalty = problemSet[tempSerial].Penalty + 1 //only penalty will inc by 1
							temp.CompilationError = problemSet[tempSerial].CompilationError
							temp.Verdict = cSubmissionList[j].Verdict
							temp.AcceptedAt = cSubmissionList[j].SubmittedAt
							problemSet[tempSerial] = temp

							//setting up original variable contestantData
							contestantData[i].SubDetails[tempSerial] = temp
						}
					}
				}
			}
		}
	}

	//now sorting the standings
	sort.SliceStable(contestantData, func(i, j int) bool {
		a, b := contestantData[i], contestantData[j]
		if a.TotalSolved != b.TotalSolved {
			return a.TotalSolved > b.TotalSolved
		}
		return a.TotalTime < b.TotalTime
	})

	//fmt.Println(contestantData)

	//preparing for returning data
	mapD := map[string]interface{}{
		"CSubmissionList":  cSubmissionList,
		"CSolvedStatus":    setSolvedStatus,
		"CAttempedStatus":  setAttempedStatus,
		"CTotalSolved":     totalSolved,
		"CTotalSubmission": totalSubmission,
		"CContestantData":  contestantData,
	}
	mapB, _ := json.Marshal(mapD)
	returnData := []byte(mapB)

	w.Header().Set("Content-Type", "application/json")
	w.Write(returnData)
}

func currentElapsedTime(currSubAt, startAt int64) int64 {
	res := currSubAt - startAt
	return res
}
