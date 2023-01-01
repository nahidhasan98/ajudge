package dimik

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

// ProbDes function for grabbing problem description
func ProbDes(pNum string) string {
	defer errorhandling.Recovery() //for panic() error errorhandling.Recovery

	//resetting previous value
	model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""

	//defining a variable for returning data
	var DimikOJDes string

	//checking if the pNum is a int or not
	tempPNum, _ := strconv.Atoi(pNum) //if pNum contains only digit, it will remain same("12"->12), otherwise become 0("12abc"->0)
	pNum = strconv.Itoa(tempPNum)

	apiURL := "https://dimikoj.com/problems/" + pNum
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	model.PTitle = strings.TrimSpace(document.Find("h2").Text())
	//DimikOJ have no pTimeLimit an pMemoryLimit

	if model.PTitle != "" { //if desired problem exist
		Des1, _ := document.Find("div[class='card-body']").Html()

		idx1 := strings.Index(Des1, "<hr/>")
		if idx1 != -1 {
			Des1 = strings.TrimSpace(Des1[idx1+5:])
		}

		var Des2 string
		document.Find("div[class='card mb-3']").Each(func(i int, s *goquery.Selection) {
			if i == 1 {
				Des2, _ = s.Html()
			}
		})
		Des2 = strings.TrimSpace(Des2)

		DimikOJDes = Des1 + Des2
		DimikOJDes += `<script data-n-head="ssr" src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.2/MathJax.js?config=TeX-AMS_HTML"></script>`
	}

	//got Title, Description. time limit & memory limit not specified
	return DimikOJDes
}
