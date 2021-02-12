package dimik

import (
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//ProbDes function for grabbing problem description
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

	model.PTitle = document.Find("h2[class='card-title']").Text()
	//DimikOJ have no pTimeLimit an pMemoryLimit

	if model.PTitle != "" { //if desired problem exist
		Des1, _ := document.Find("div[class='card-body pstatement']").Html()
		Des2, _ := document.Find("div[class='card-body bg-light']").Html()

		DimikOJDes = Des1 + Des2
	}

	//got Title, Description. time limit & memory limit not specified
	return DimikOJDes
}
