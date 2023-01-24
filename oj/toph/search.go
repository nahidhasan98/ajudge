package toph

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

// Search function for searching problem in Toph
func Search(sQuery string) []model.ProblemList {
	defer errorhandling.Recovery() //for panic() error Recovery

	//data that will be collected
	var problemList []model.ProblemList

	sQuery = strings.ReplaceAll(sQuery, ` `, `+`) //making a valid query like: leap year => leap+year
	apiURL := "https://toph.co/search?type=problems&q=" + sQuery
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	document.Find("td[class='flow']").Each(func(index int, mixedStr *goquery.Selection) {
		tempNum, _ := mixedStr.Find("a").Attr("href") //Toph have no pNum. So we taking pLink (like:/p/copycat)

		tempNum = tempNum[3:] //avoiding /p/ from the link
		tempName := mixedStr.Find("a[href^='/p/']").Text()

		var temp model.ProblemList
		temp.OJ = "Toph"
		temp.PNum = tempNum
		temp.PName = tempName

		problemList = append(problemList, temp)
	})

	return problemList
}
