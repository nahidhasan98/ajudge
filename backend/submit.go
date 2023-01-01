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

// Submit function for submitting a problem solution
func Submit(w http.ResponseWriter, r *http.Request) {
	contestID := 0 //for differentiate from contest submission or normal submission
	serialIndex := ""

	if r.Method == "POST" {
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

			title := strings.TrimSpace(document.Find("h2").Text())

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
			apiURL := "https://www.beecrowd.com.br/judge/en/problems/view/" + pNum
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
