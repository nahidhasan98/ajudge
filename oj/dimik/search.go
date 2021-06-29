package dimik

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//Search function for searching problem in DimikOJ
func Search(sQuery string) []model.ProblemList {
	defer errorhandling.Recovery() //for panic() error errorhandling.Recovery

	//data that will be collected
	var problemList []model.ProblemList

	sQuery = strings.ReplaceAll(sQuery, ` `, `+`) //making a valid query like: leap year => leap+year
	apiURL := "https://dimikoj.com/problems/" + sQuery
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	var tempNum, tempName string
	tempNum = sQuery
	tempName = document.Find("h2[class='card-title']").Text()

	if tempName != "" {
		var temp model.ProblemList
		temp.OJ = "DimikOJ"
		temp.PNum = tempNum
		temp.PName = tempName

		problemList = append(problemList, temp)
	}

	return problemList
}
