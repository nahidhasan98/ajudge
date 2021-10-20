package backend

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/oj/dimik"
	"github.com/nahidhasan98/ajudge/oj/toph"
	"github.com/nahidhasan98/ajudge/oj/uri"
	"github.com/nahidhasan98/ajudge/oj/vjudge"
)

//Problem function for searching problem from different OJ
func Problem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	//problem list will be gathered by frontend ajax call

	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "Problem"
	model.Info["PageTitle"] = "Problem | AJudge"

	model.Tpl.ExecuteTemplate(w, "problem.gohtml", model.Info)
}

//ProblemView function for grabbing a problem description
func ProblemView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	OJpNum := strings.TrimPrefix(path, "/problemView/")

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

	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "ProblemView"
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

	model.Tpl.ExecuteTemplate(w, "problemView.gohtml", model.Info)
}

//Origin function for finding a problem origin
func Origin(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	OJpNum := strings.TrimPrefix(path, "/origin/")

	need := "-"
	index := strings.Index(OJpNum, need)
	OJ, pNum := "", ""

	if index == -1 { //url is not like this "/origin/OJ-pNum"
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
		return
	}
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
	if OJ == "DimikOJ" {
		model.POrigin = "https://dimikoj.com/problems/" + pNum

		//checking wheather this problem exist or not
		response := dimik.GETRequest(model.POrigin)
		defer response.Body.Close()

		res, _ := ioutil.ReadAll(response.Body)
		sRes := string(res)

		need := "card-body pstatement"
		index := strings.Index(sRes, need)

		if index == -1 { //didn't get any problem
			errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
			return
		}
		//redirecting to the origin site
		http.Redirect(w, r, model.POrigin, http.StatusSeeOther)
		return
	} else if OJ == "Toph" {
		model.POrigin = "https://toph.co/p/" + pNum

		//checking wheather this problem exist or not
		response := toph.GETRequest(model.POrigin)
		defer response.Body.Close()

		res, _ := ioutil.ReadAll(response.Body)
		sRes := string(res)

		need := "artifact"
		index := strings.Index(sRes, need)

		if index == -1 { //didn't get any problem
			errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
			return
		}
		//redirecting to the origin site
		http.Redirect(w, r, model.POrigin, http.StatusSeeOther)
		return
	} else if OJ == "URI" {
		model.POrigin = "https://www.urionlinejudge.com.br/judge/en/problems/view/" + pNum

		//checking wheather this problem exist or not
		response := uri.GETRequest(model.POrigin)
		defer response.Body.Close()

		res, _ := ioutil.ReadAll(response.Body)
		sRes := string(res)

		need := "iframe"
		index := strings.Index(sRes, need)

		if index == -1 { //didn't get any problem
			errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
			return
		}
		//redirecting to the origin site
		http.Redirect(w, r, model.POrigin, http.StatusSeeOther)
		return
	} else {
		//Finding origin
		pURL := "https://vjudge.net/problem/" + OJ + "-" + pNum

		//checking wheather this problem exist or not
		response := vjudge.GETRequest(pURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		pOriginText, _ := document.Find("span[class='origin']").Find("a").Attr("href")

		if pOriginText == "" { //didn't get any problem
			errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
			return
		}
		//got a problem and origin
		//getting origin link
		model.POrigin = getOriginLink("https://vjudge.net" + pOriginText)

		//redirecting to the origin site
		http.Redirect(w, r, model.POrigin, http.StatusSeeOther)
	}
}

//function used above by this particular file
func getOriginLink(apiURL string) string {
	req, _ := http.NewRequest("GET", apiURL, nil)

	//setting up request to prevent auto redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	errorhandling.Check(err)
	defer resp.Body.Close()
	header, _ := resp.Location()

	return header.String()
}
