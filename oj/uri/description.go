package uri

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//ProbDes function for grabbing problem description
func ProbDes(pNum string) string {
	defer errorhandling.Recovery() //for panic() error Recovery

	//resetting previous value
	model.PTitle, model.PTimeLimit, model.PMemoryLimit, model.PSourceLimit, model.PDesSrcVJ, model.POrigin = "", "", "", "", "", ""

	//defining a variable for returning data
	var URIDes = ""

	//checking if the pNum is a int or not
	tempPNum, _ := strconv.Atoi(pNum) //if pNum contains only digit, it will remain same("12"->12), otherwise become 0("12abc"->0)
	pNum = strconv.Itoa(tempPNum)

	//getting problem description
	apiURL := "https://www.beecrowd.com.br/judge/en/problems/view/" + pNum
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)
	pDesSrcURI, _ := document.Find("iframe").Attr("src") //this is prob description source

	apiURL = "https://www.beecrowd.com.br/" + pDesSrcURI
	response = GETRequest(apiURL)
	defer response.Body.Close()

	document, err = goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	model.PTitle = document.Find("div[class='header']").Find("h1").Text()

	if model.PTitle != "" { //if desired problem exist
		model.PTimeLimit = document.Find("div[class='header']").Find("strong").Text()
		model.PTimeLimit = strings.TrimPrefix(model.PTimeLimit, "Timelimit: ")

		URIDes, _ = document.Find("div[class='problem']").Html()
		if URIDes == "" { //then it is SQL category problem
			URIDes, _ = document.Find("div[class='problem-sql']").Html()
			model.PTitle += " (SQL Problem)"
			model.PTimeLimit = strings.TrimPrefix(model.PTimeLimit, "SQLTimelimit: ")
		}
		model.PTimeLimit += "s"
		//got Title,TimeLimit,Description

		//Getting Memory Limit
		apiURL := "https://www.beecrowd.com.br/judge/en/problems/view/" + pNum
		response := GETRequest(apiURL)
		defer response.Body.Close()

		document, err = goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		var tempMemoryArray []string
		var tempMemory string
		document.Find("div[id='page-name-c']").Find("ul").Find("li").Each(func(index int, mixedStr *goquery.Selection) {
			tempMemory = mixedStr.Text()
			tempMemory = strings.TrimSpace(tempMemory)
			tempMemoryArray = append(tempMemoryArray, tempMemory)
		})
		model.PMemoryLimit = tempMemoryArray[8]                                       //memoryLimit on the 8 indexed list
		model.PMemoryLimit = strings.TrimPrefix(model.PMemoryLimit, "Memory Limit: ") //Taking only Value, leaving the extra text
		//got memory limit
	}

	//adding mathjax script for math equations/Katex
	URIDes = URIDes + `<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
	<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3.0.1/es5/tex-mml-chtml.js"></script>`

	return URIDes
}
