package backend

import (
	"encoding/json"
	"html"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/oj/dimik"
	"github.com/nahidhasan98/ajudge/oj/toph"
	"github.com/nahidhasan98/ajudge/oj/uri"
	"github.com/nahidhasan98/ajudge/oj/vjudge"
	"github.com/nahidhasan98/ajudge/vault"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//CheckDB function for checking username/email already exist or not in DB
func CheckDB(w http.ResponseWriter, r *http.Request) {
	var username, email string
	usernameList := r.URL.Query()["username"]
	if len(usernameList) > 0 {
		username = usernameList[0]
	}
	emailList := r.URL.Query()["email"]
	if len(emailList) > 0 {
		email = emailList[0]
	}

	//declaring variable for returning data
	type Data struct {
		IsUsernameExist, IsEmailExist bool
	}
	var data Data

	//checking for username/email already exist or not
	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	userCollection := DB.Collection("user")

	//getting data for this user from DB
	var userData model.UserData
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userData)

	if err == mongo.ErrNoDocuments { //found no rows (username available)
		data.IsUsernameExist = false
	} else { //found a row
		data.IsUsernameExist = true
	}

	err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&userData)

	if err == mongo.ErrNoDocuments { //found no rows (email available)
		data.IsEmailExist = false
	} else { //found a row
		data.IsEmailExist = true
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(data)
	w.Write(b)
}

//ProblemList function for searching problem from OJ
func ProblemList(w http.ResponseWriter, r *http.Request) {
	OJList := r.URL.Query()["OJ"]
	pNumList := r.URL.Query()["pNum"]
	pNameList := r.URL.Query()["pName"]

	var OJ, pNum, pName string
	if len(OJList) > 0 {
		OJ = OJList[0]
	}
	if len(pNumList) > 0 {
		pNum = pNumList[0]
	}
	if len(pNameList) > 0 {
		pName = pNameList[0]
	}
	//fmt.Println(OJ,pNum,pName)

	var pListFinal []model.ProblemList

	if pNum == "" && pName == "" {
		pListFinal = multipleSearch(OJ)
	} else {
		pListFinal = singleSearch(OJ, pNum, pName)
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(pListFinal)
	w.Write(b)
}

func multipleSearch(OJ string) []model.ProblemList {
	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	problemCollection := DB.Collection("problem")

	//getting data from DB
	//setting up options for retrieving data from DB
	var pipe []primitive.M

	if OJ != "" && OJ != "All" {
		pipe = []bson.M{
			{"$match": bson.M{
				"OJ": OJ,
			}},
		}
	}
	pipe = append(pipe, bson.M{
		"$sample": bson.M{
			"size": 2000,
		},
	})

	cursor, err := problemCollection.Aggregate(ctx, pipe)
	errorhandling.Check(err)

	var pList []model.ProblemList

	err = cursor.All(ctx, &pList)
	errorhandling.Check(err)

	err = cursor.Close(ctx)
	errorhandling.Check(err)

	return pList
}

func singleSearch(OJ, pNum, pName string) []model.ProblemList {
	var pListFinal, pList []model.ProblemList

	if OJ == "" || OJ == "All" || OJ == "DimikOJ" {
		if pNum != "" {
			pList = dimik.Search(pNum)
		} else if pName != "" {
			pList = dimik.Search(pName)
		}
		pListFinal = append(pListFinal, pList...)
	}
	if OJ == "" || OJ == "All" || OJ == "Toph" {
		if pNum != "" {
			pList = toph.Search(pNum)
		} else if pName != "" {
			pList = toph.Search(pName)
		}
		pListFinal = append(pListFinal, pList...)
	}
	if OJ == "" || OJ == "All" || OJ == "URI" {
		if pNum != "" {
			pList = uri.Search(pNum)
		} else if pName != "" {
			pList = uri.Search(pName)
		}
		pListFinal = append(pListFinal, pList...)
	}
	if OJ != "DimikOJ" && OJ != "Toph" && OJ != "URI" {
		pList := vjudge.Search(OJ, pNum, pName)

		pListFinal = append(pListFinal, pList...)
	}

	//Shuffling the problem list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pListFinal), func(i, j int) { pListFinal[i], pListFinal[j] = pListFinal[j], pListFinal[i] })

	return pListFinal
}

//GetUserSubmission function for grabbing all submissions of a user
func GetUserSubmission(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	user := strings.TrimPrefix(path, "/userSubmission/")

	type setData struct {
		OJ   string
		PNum string
	}
	setSolved := make(map[setData]bool)    // New empty set
	setAttempted := make(map[setData]bool) // New empty set

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	submissionCollection := DB.Collection("submission")

	//getting data from DB
	type subData struct {
		SubList          []model.SubmissionData
		ProblemSolved    int
		ProblemAttempted int
	}
	var subDataFinal subData

	//setting uo options for retrieving data from DB
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "subID", Value: -1}}) //descending order
	cursor, err := submissionCollection.Find(ctx, bson.D{{Key: "username", Value: user}}, opts)
	errorhandling.Check(err)

	// Finding multiple documents(rows) returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.SubmissionData
		err := cursor.Decode(&temp)
		errorhandling.Check(err)
		temp.SourceCode = html.EscapeString(temp.SourceCode) //specially for reserving newline

		//checking if this user is appropriate to see the source code or not
		session, _ := model.Store.Get(r, "mysession")
		if temp.Username != session.Values["username"] { //source code will be provided to the correct owner only
			temp.SourceCode = ""
		}

		if temp.ContestID == 0 {
			subDataFinal.SubList = append(subDataFinal.SubList, temp)
		} else {
			//If this submission is a contest submission and that contest in currently running,
			// then this submission should not be sent to fronted right now

			//taking DB collection/table to a variable
			contestCollection := DB.Collection("contest")

			var dbQuery model.ContestData
			err := contestCollection.FindOne(ctx, bson.M{"contestID": temp.ContestID}).Decode(&dbQuery)
			if err == mongo.ErrNoDocuments {
				errorPage(w, http.StatusInternalServerError) //http.StatusBadRequest = 400
				return
			}

			//converting duration, string to int64
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

			cDuration := int64((cHour * 60 * 60) + (cMin * 60))

			contestEndAt := dbQuery.StartAt + cDuration
			timeNow := time.Now().Unix()

			if timeNow > contestEndAt {
				subDataFinal.SubList = append(subDataFinal.SubList, temp)
			}
		}
		var tempSet setData
		tempSet.OJ = temp.OJ
		tempSet.PNum = temp.PNum

		if temp.Verdict == "Accepted" {
			setSolved[tempSet] = true //adding to set
		}
		setAttempted[tempSet] = true //adding to set
	}
	subDataFinal.ProblemSolved = len(setSolved)
	subDataFinal.ProblemAttempted = len(setAttempted)

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(subDataFinal)
	w.Write(b)
}

//Verdict function for collecting verdict from OJ
func Verdict(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	subIDText := strings.TrimPrefix(path, "/verdict/subID=")
	subID, _ := strconv.Atoi(subIDText) //In DB, subID stored as int32

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	submissionCollection := DB.Collection("submission")
	userCollection := DB.Collection("user")
	rankOJCollection := DB.Collection("rankOJ")

	//getting data from DB
	var submissionData model.SubmissionData
	err := submissionCollection.FindOne(ctx, bson.M{"subID": subID}).Decode(&submissionData)
	errorhandling.Check(err)

	//initializing variables
	rVerdict := "Waiting"
	rRuntime := "N/A"
	rMemory := "N/A"
	rTerminalVerdict := false

	if submissionData.OJ == "DimikOJ" {
		//first login to DimikOJ
		if dimik.Login() != "success" { //if login unsuccessful
			errorPage(w, http.StatusInternalServerError)
			return
		}
		//DimikOJ login success

		apiURL := "https://dimikoj.com/submissions/" + submissionData.VID
		req, err := http.NewRequest("GET", apiURL, nil)
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		//getting verdict
		rVerdict = document.Find("div[class='col-xl-4 col-lg-5 col-md-8']").Find("span").Text()

		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], rVerdict) { //got terminal/final verdict
			rTerminalVerdict = true

			//now getting runtime & memory
			var tempValue string
			var tempValueArray []string
			document.Find("div[class='col-xl-4 col-lg-5 col-md-8']").Find("p").Each(func(index int, mixedStr *goquery.Selection) {
				tempValue = mixedStr.Text()
				tempValue = strings.TrimSpace(tempValue)
				tempValueArray = append(tempValueArray, tempValue)
			})

			if len(tempValueArray) >= 2 {
				rRuntime = tempValueArray[0]
				//removing extra text from runtime
				need := " "
				index := strings.Index(rRuntime, need)
				rRuntime = rRuntime[0:index]
				rRuntime += " s"
			}
		}
	} else if submissionData.OJ == "Toph" {
		//first login to Toph
		if toph.Login() != "success" { //if login unsuccessful
			errorPage(w, http.StatusInternalServerError)
			return
		}
		//Toph login success

		apiURL := "https://toph.co/s/" + submissionData.VID
		req, err := http.NewRequest("GET", apiURL, nil)
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		var tempArray []string
		var temp string
		document.Find("td").Each(func(index int, IDSeg *goquery.Selection) {
			temp = IDSeg.Text()
			temp = strings.TrimSpace(temp)
			tempArray = append(tempArray, temp)
		})

		rVerdict = tempArray[5]
		rRuntime = tempArray[6]
		rMemory = tempArray[7]

		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], rVerdict) {
			rTerminalVerdict = true
		}
	} else if submissionData.OJ == "URI" {
		//first login to URI
		if uri.Login() != "success" { //if login unsuccessful
			errorPage(w, http.StatusInternalServerError)
			return
		}
		//URI login success

		apiURL := "https://www.beecrowd.com.br/judge/en/runs/code/" + submissionData.VID
		response := uri.GETRequest(apiURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		rVerdict = document.Find("dd").Find("span").Text()
		rVerdict = strings.TrimSpace(rVerdict)

		//checking if this is final verdict or not. This checking is for wheather this function should be called again or not
		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], rVerdict) {
			rTerminalVerdict = true

			//Finding runtime...
			var runtimeList []string
			var tempTime string
			document.Find("dd").Each(func(index int, mixedStr *goquery.Selection) {
				tempTime = mixedStr.Text()
				tempTime = strings.TrimSpace(tempTime)
				runtimeList = append(runtimeList, tempTime)
			})
			rRuntime = runtimeList[3]
			//got runtime
		}
	} else {
		//doesn't require to login to see verdict in VJudge
		apiURL := "https://vjudge.net/solution/data/" + submissionData.VID
		response := vjudge.GETRequest(apiURL)
		body, _ := ioutil.ReadAll(response.Body)

		type Res struct {
			Status  string `json:"status"`
			Runtime int    `json:"runtime"`
			Memory  int    `json:"memory"`
		}
		var res Res
		json.Unmarshal(body, &res)

		rVerdict = res.Status
		rRuntime = strconv.Itoa(res.Runtime) + " ms"
		rMemory = strconv.Itoa(res.Memory) + " kB"

		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], res.Status) {
			rTerminalVerdict = true
		}
	}

	//Now Updating totalSolved for this user if solution is accepted - before that checking if this user already solved this problem or not
	var tempSubData model.SubmissionData
	err = submissionCollection.FindOne(ctx, bson.M{"username": submissionData.Username, "OJ": submissionData.OJ, "pNum": submissionData.PNum, "verdict": "Accepted"}).Decode(&tempSubData)

	if err == mongo.ErrNoDocuments && rVerdict == "Accepted" { //not solved before & solved right now
		//updating user totalSolved to DB
		updateField := bson.D{
			{Key: "$inc", Value: bson.D{ //incrementing totalSolved by 1
				{Key: "totalSolved", Value: 1},
			}},
		}
		_, err := userCollection.UpdateOne(ctx, bson.M{"username": submissionData.Username}, updateField)
		errorhandling.Check(err)
	}

	//Now Updating totalSolved for rankOJ if solution is accepted - before that checking if this problem already solved or not
	var tempSubData2 model.SubmissionData
	err = submissionCollection.FindOne(ctx, bson.M{"OJ": submissionData.OJ, "pNum": submissionData.PNum, "verdict": "Accepted"}).Decode(&tempSubData2)

	if err == mongo.ErrNoDocuments && rVerdict == "Accepted" { //not solved before & solved right now
		//updating rankOJ totalSolved to DB
		updateField := bson.D{
			{Key: "$inc", Value: bson.D{ //incrementing totalSolved by 1
				{Key: "totalSolved", Value: 1},
			}},
		}
		_, err = rankOJCollection.UpdateOne(ctx, bson.M{"OJ": submissionData.OJ}, updateField)
		errorhandling.Check(err)
	}

	//finally updating submission result fields in DB
	updateField := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "verdict", Value: rVerdict},
			{Key: "timeExec", Value: rRuntime},
			{Key: "memoryExec", Value: rMemory},
			{Key: "terminalVerdict", Value: rTerminalVerdict},
		}},
	}
	_, err = submissionCollection.UpdateOne(ctx, bson.M{"subID": subID}, updateField)
	errorhandling.Check(err)

	//preparing for returning data
	mapD := map[string]interface{}{
		"Status":          rVerdict,
		"Runtime":         rRuntime,
		"Memory":          rMemory,
		"TerminalVerdict": rTerminalVerdict,
	}
	mapB, _ := json.Marshal(mapD)
	returnData := []byte(mapB)

	w.Header().Set("Content-Type", "application/json")
	w.Write(returnData)

	// notofy to discord
	disData := model.SubmissionData{
		SubID:      subID,
		TimeExec:   rRuntime,
		MemoryExec: rMemory,
		Verdict:    rVerdict,
	}
	discord := discord.Init()
	discord.EditMessage(disData, "submission")
}

//Rejudge function for rejudging verdict from OJ
func Rejudge(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	subIDText := strings.TrimPrefix(path, "/rejudge/subID=")
	subID, _ := strconv.Atoi(subIDText) //In DB, subID stored as int32

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	submissionCollection := DB.Collection("submission")
	userCollection := DB.Collection("user")
	rankOJCollection := DB.Collection("rankOJ")

	//getting data from DB
	var submissionData model.SubmissionData
	err := submissionCollection.FindOne(ctx, bson.M{"subID": subID}).Decode(&submissionData)
	errorhandling.Check(err)

	//initializing variables
	rVerdict := "Waiting"
	rRuntime := "N/A"
	rMemory := "N/A"
	rTerminalVerdict := false

	if submissionData.OJ == "DimikOJ" {
		//first login to DimikOJ
		if dimik.Login() != "success" { //if login unsuccessful
			errorPage(w, http.StatusInternalServerError)
			return
		}
		//DimikOJ login success

		apiURL := "https://dimikoj.com/submissions/" + submissionData.VID
		req, err := http.NewRequest("GET", apiURL, nil)
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		//getting verdict
		rVerdict = document.Find("div[class='col-xl-4 col-lg-5 col-md-8']").Find("span").Text()

		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], rVerdict) { //got terminal/final verdict
			rTerminalVerdict = true

			//now getting runtime & memory
			var tempValue string
			var tempValueArray []string
			document.Find("div[class='col-xl-4 col-lg-5 col-md-8']").Find("p").Each(func(index int, mixedStr *goquery.Selection) {
				tempValue = mixedStr.Text()
				tempValue = strings.TrimSpace(tempValue)
				tempValueArray = append(tempValueArray, tempValue)
			})

			if len(tempValueArray) >= 2 {
				rRuntime = tempValueArray[0]
				//removing extra text from runtime
				need := " "
				index := strings.Index(rRuntime, need)
				rRuntime = rRuntime[0:index]
				rRuntime += " s"
			}
		}
	} else if submissionData.OJ == "Toph" {
		//first login to Toph
		if toph.Login() != "success" { //if login unsuccessful
			errorPage(w, http.StatusInternalServerError)
			return
		}
		//Toph login success

		apiURL := "https://toph.co/s/" + submissionData.VID
		req, err := http.NewRequest("GET", apiURL, nil)
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		var tempArray []string
		var temp string
		document.Find("td").Each(func(index int, IDSeg *goquery.Selection) {
			temp = IDSeg.Text()
			temp = strings.TrimSpace(temp)
			tempArray = append(tempArray, temp)
		})

		rVerdict = tempArray[5]
		rRuntime = tempArray[6]
		rMemory = tempArray[7]

		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], rVerdict) {
			rTerminalVerdict = true
		}
	} else if submissionData.OJ == "URI" {
		//first login to URI
		if uri.Login() != "success" { //if login unsuccessful
			errorPage(w, http.StatusInternalServerError)
			return
		}
		//URI login success

		apiURL := "https://www.beecrowd.com.br/judge/en/runs/code/" + submissionData.VID
		response := uri.GETRequest(apiURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		rVerdict = document.Find("dd").Find("span").Text()
		rVerdict = strings.TrimSpace(rVerdict)

		//checking if this is final verdict or not. This checking is for wheather this function should be called again or not
		if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], rVerdict) {
			rTerminalVerdict = true

			//Finding runtime...
			var runtimeList []string
			var tempTime string
			document.Find("dd").Each(func(index int, mixedStr *goquery.Selection) {
				tempTime = mixedStr.Text()
				tempTime = strings.TrimSpace(tempTime)
				runtimeList = append(runtimeList, tempTime)
			})
			rRuntime = runtimeList[3]
			//got runtime
		}
	} else {
		//rejudging process on vjudge completed by 2 steps

		//preparing data for POST Request
		postData := url.Values{}

		//submitting to Vjudge for rejudging
		apiURL := "https://vjudge.net/problem/rejudge/" + submissionData.VID
		req, _ := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()
		//subbmission done
		body, _ := ioutil.ReadAll(response.Body)

		if string(body) == "OK" {
			//preparing data for POST Request
			postData := url.Values{
				"runIds[]": {submissionData.VID},
			}

			//submitting to Vjudge for rejudging
			apiURL := "https://vjudge.net/status/dataById"
			req, _ := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

			response, err := model.Client.Do(req)
			errorhandling.Check(err)
			defer response.Body.Close()
			//subbmission done
			body, _ := ioutil.ReadAll(response.Body)

			type Res struct {
				Status  string `json:"status"`
				Runtime int    `json:"runtime"`
				Memory  int    `json:"memory"`
			}
			res := map[string]Res{}
			json.Unmarshal(body, &res)

			rVerdict = res[submissionData.VID].Status
			rRuntime = strconv.Itoa(res[submissionData.VID].Runtime) + " ms"
			rMemory = strconv.Itoa(res[submissionData.VID].Memory) + " kB"

			if model.IsExistInTV(submissionData.OJ, model.TerminalVerdict[submissionData.OJ], res[submissionData.VID].Status) {
				rTerminalVerdict = true
			}
		} else {
			errorPage(w, http.StatusInternalServerError)
			return
		}
	}

	//Now Updating totalSolved for this user if solution is accepted - before that checking if this user already solved this problem or not
	var tempSubData model.SubmissionData
	res := submissionCollection.FindOne(ctx, bson.M{"username": submissionData.Username, "OJ": submissionData.OJ, "pNum": submissionData.PNum, "verdict": "Accepted"}).Decode(&tempSubData)

	if res == mongo.ErrNoDocuments && rVerdict == "Accepted" { //not solved before & soved right now
		//updating user totalSolved to DB
		updateField := bson.D{
			{Key: "$inc", Value: bson.D{ //incrementing totalSolved by 1
				{Key: "totalSolved", Value: 1},
			}},
		}
		_, err = userCollection.UpdateOne(ctx, bson.M{"username": submissionData.Username}, updateField)
		errorhandling.Check(err)
	}

	//Now Updating totalSolved for rankOJ if solution is accepted - before that checking if this problem already solved or not
	var tempSubData2 model.SubmissionData
	res = submissionCollection.FindOne(ctx, bson.M{"OJ": submissionData.OJ, "pNum": submissionData.PNum, "verdict": "Accepted"}).Decode(&tempSubData2)

	if res == mongo.ErrNoDocuments && rVerdict == "Accepted" { //not solved before & soved right now
		//updating rankOJ totalSolved to DB
		updateField := bson.D{
			{Key: "$inc", Value: bson.D{ //incrementing totalSolved by 1
				{Key: "totalSolved", Value: 1},
			}},
		}
		_, err = rankOJCollection.UpdateOne(ctx, bson.M{"OJ": submissionData.OJ}, updateField)
		errorhandling.Check(err)
	}

	//finally updating submission result fields in DB
	updateField := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "verdict", Value: rVerdict},
			{Key: "timeExec", Value: rRuntime},
			{Key: "memoryExec", Value: rMemory},
			{Key: "terminalVerdict", Value: rTerminalVerdict},
		}},
	}
	_, err = submissionCollection.UpdateOne(ctx, bson.M{"subID": subID}, updateField)
	errorhandling.Check(err)

	//preparing for returning data
	mapD := map[string]interface{}{
		"Status":          rVerdict,
		"Runtime":         rRuntime,
		"Memory":          rMemory,
		"TerminalVerdict": rTerminalVerdict,
	}
	mapB, _ := json.Marshal(mapD)
	returnData := []byte(mapB)

	w.Header().Set("Content-Type", "application/json")
	w.Write(returnData)

	// notofy to discord
	disData := model.SubmissionData{
		SubID:      subID,
		TimeExec:   rRuntime,
		MemoryExec: rMemory,
		Verdict:    rVerdict,
	}
	discord := discord.Init()
	discord.EditMessage(disData, "submission")
}

//SubHistory function for retrieving User's previous submission history of a specific problem
func SubHistory(w http.ResponseWriter, r *http.Request) {
	OJList := r.URL.Query()["OJ"]
	pNumList := r.URL.Query()["pNum"]
	userList := r.URL.Query()["user"]

	var OJ, pNum, user string
	if len(OJList) > 0 {
		OJ = OJList[0]
	}
	if len(pNumList) > 0 {
		pNum = pNumList[0]
	}
	if len(userList) > 0 {
		user = userList[0]
	}

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	submissionCollection := DB.Collection("submission")

	//getting data for this user from DB
	var subHistoryList []model.SubmissionData

	//setting uo options for retrieving data from DB
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "submittedAt", Value: -1}}) //sorting by submiitedAt in descending
	cursor, err := submissionCollection.Find(ctx, bson.M{"username": user, "OJ": OJ, "pNum": pNum}, opts)
	errorhandling.Check(err)

	//Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp model.SubmissionData
		err := cursor.Decode(&temp)
		errorhandling.Check(err)

		subHistoryList = append(subHistoryList, temp)
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(subHistoryList)
	w.Write(b)
}

//GetCaptcha function for confirming captcha correct or not
func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	captchaUser := strings.TrimPrefix(r.URL.Path, "/captcha/")
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	//preparing data for post
	postData := url.Values{
		"secret":   {vault.CaptchaKey},
		"response": {captchaUser},
	}

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := model.Client.Do(req)
	errorhandling.Check(err)
	defer response.Body.Close()

	res, _ := ioutil.ReadAll(response.Body)

	type captchaRes struct {
		Success bool `json:"success"`
	}
	var capRes captchaRes
	json.Unmarshal(res, &capRes)

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(capRes)
	w.Write(b)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := model.Store.Get(r, "mysession")
	var response string

	if session.Values["isLogin"] == true {
		if model.IsAccVerifed(r) {
			response = "true"
		} else {
			model.PopUpCause = "verifyRequired"
			response = "notVerified"
		}
	} else {
		model.PopUpCause = "loginRequired"
		response = "false"
	}
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(response)
	w.Write(b)
}

func GetCombinedStandings(w http.ResponseWriter, r *http.Request) {
	cID := r.FormValue("ids")
	cID = cID[1 : len(cID)-1]
	cIDs := strings.Split(cID, ",")

	type perStatus = struct {
		ConID     int
		PerSolved int
		PerTime   int64
	}
	type contestant = struct {
		Username         string
		TotalSolved      int
		TotalTime        int64
		PerContestStatus []perStatus
	}
	cData := make(map[string]contestant) //map for tracking same username from different contest
	var rData map[string]interface{}

	for _, val := range cIDs {
		apiURL := "https://ajudge.net/dataContest/" + val
		req, err := http.NewRequest("GET", apiURL, nil)
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

		response, err := model.Client.Do(req)
		errorhandling.Check(err)

		body, _ := ioutil.ReadAll(response.Body)

		json.Unmarshal(body, &rData)

		rd1 := rData["CContestantData"]
		if rd1 != nil { //if standings/submission is not empty
			rd2 := rd1.([]interface{})

			for j := 0; j < len(rd2); j++ {
				rd3 := rd2[j].(map[string]interface{})

				username := rd3["Username"].(string)

				ts := cData[username].TotalSolved
				ts += int(rd3["TotalSolved"].(float64))

				tt := cData[username].TotalTime
				tt += int64(rd3["TotalTime"].(float64))

				cid, _ := strconv.Atoi(val)
				tPer := perStatus{
					ConID:     cid,
					PerSolved: int(rd3["TotalSolved"].(float64)),
					PerTime:   int64(rd3["TotalTime"].(float64)),
				}

				tp := cData[username].PerContestStatus
				tp = append(tp, tPer)

				temp := contestant{
					TotalSolved:      ts,
					TotalTime:        tt,
					PerContestStatus: tp,
				}
				cData[username] = temp
			}
		}
	}

	//mow taking in []struct for sorting the list
	var contestantData []contestant
	for key, val := range cData {
		//first sorting the per Status according to Contest ID
		sort.SliceStable(val.PerContestStatus, func(i, j int) bool {
			a, b := val.PerContestStatus[i], val.PerContestStatus[j]
			if a.ConID != b.ConID {
				return a.ConID < b.ConID
			}
			return true
		})
		temp := contestant{
			Username:         key,
			TotalSolved:      val.TotalSolved,
			TotalTime:        val.TotalTime,
			PerContestStatus: val.PerContestStatus,
		}
		contestantData = append(contestantData, temp)
	}

	//now sorting the standings
	sort.SliceStable(contestantData, func(i, j int) bool {
		a, b := contestantData[i], contestantData[j]
		if a.TotalSolved != b.TotalSolved {
			return a.TotalSolved > b.TotalSolved
		}
		return a.TotalTime < b.TotalTime
	})
	//fmt.Fprintln(w, contestantData)

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(contestantData)
	w.Write(b)
}
