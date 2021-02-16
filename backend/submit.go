package backend

import (
	"net/http"
	"net/url"
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

//Submit function for submitting a problem solution
func Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := model.Store.Get(r, "mysession")
	model.LastPage = r.URL.Path

	contestID := 0 //for differentiate from contest submission or normal submission
	serialIndex := ""

	if r.Method != "POST" {
		if session.Values["isLogin"] == true {
			if model.IsAccVerifed(r) {
				path := r.URL.Path
				OJpNum := strings.TrimPrefix(path, "/submit/")

				need := "-"
				index := strings.Index(OJpNum, need)
				OJ, pNum := "", ""

				if index == -1 { //url is not like this "/submit/OJ-pNum"
					model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""
					errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
					return
				}
				//url is something like this "/submit/OJ-pNum"
				OJ = OJpNum[0:index]
				pNum = OJpNum[index+1:]

				if model.OJSet[OJ] == false || pNum == "" { //bad url, not OJ & pNum specified
					model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""
					errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
					return
				}
				// got something in OJ and pNum
				model.PDesSrcVJ = "" //resetting for now
				allowSubmit := false

				if OJ == "DimikOJ" {
					dimik.ProbDes(pNum)

					if model.PTitle == "" { //didn't get any problem
						model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""
						errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
						return
					}
					//else got a problem with this OJ & pNum
					allowSubmit = true
				} else if OJ == "Toph" {
					toph.ProbDes(pNum)

					if model.PTitle == "" { //didn't get any problem
						model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""
						errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
						return
					}
					//else got a problem with this OJ & pNum
					allowSubmit = true
				} else if OJ == "URI" {
					uri.ProbDes(pNum)

					if model.PTitle == "" { //didn't get any problem
						model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""
						errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
						return
					}
					//else got a problem with this OJ & pNum
					allowSubmit = true
				} else {
					//(Finding problem) Verifying that problem exist with this OJ & pNum
					var status int
					_, allowSubmit, status = vjudge.ProbDes(OJ, pNum)

					if model.PDesSrcVJ == "" { //didn't get any problem (model.PTitle/model.PDesSrcVJ both will be empty)
						model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""
						errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
						return
					}
					//got a problem with this OJ & pNum
					//checking whether problem submission allowed or not
					if allowSubmit == true && status == 0 {
						allowSubmit = true
					}
				}

				if allowSubmit == true {
					model.Info["Username"] = session.Values["username"]
					model.Info["Password"] = session.Values["password"]
					model.Info["IsLogged"] = session.Values["isLogin"]
					model.Info["PageName"] = "Submission"
					model.Info["PageTitle"] = "Submission | AJudge"
					model.Info["Lastpage"] = model.LastPage
					model.Info["PopUpCause"] = model.PopUpCause
					model.Info["ErrorType"] = model.ErrorType
					model.Info["OJ"] = OJ
					model.Info["PNum"] = pNum
					model.Info["PName"] = model.PTitle
					model.Info["ContestID"] = contestID

					model.Tpl.ExecuteTemplate(w, "submit.gohtml", model.Info)

					//clearing up values (because it may be used in wrong place unintentionally)
					model.PopUpCause = ""
					model.Info["PopUpCause"] = model.PopUpCause
				} else if allowSubmit == false {
					link := "/problemView/" + OJ + "-" + pNum
					http.Redirect(w, r, link, http.StatusSeeOther)
					return
				}
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

//Result function for verdict page
func Result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	//here only 2 info added. some info are come from submit page
	model.Info["PageName"] = "Result"
	model.Info["PageTitle"] = "Result | AJudge"

	model.Tpl.ExecuteTemplate(w, "result.gohtml", model.Info)

	//clearing things up
	model.Info["SubID"] = ""
	model.Info["OJ"] = ""
	model.Info["PNum"] = ""
	model.Info["Language"] = ""
	model.Info["SourceCode"] = ""
	model.Info["SubmittedAt"] = ""
}
