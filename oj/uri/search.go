package uri

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//Search function for searching problem in URI
func Search(sQuery string) []model.ProblemList {
	defer errorhandling.Recovery() //for panic() error Recovery

	//data that will be collected
	var problemList []model.ProblemList

	sQuery = strings.ReplaceAll(sQuery, ` `, `+`) //making a valid query like: leap year => leap+year
	apiURL := "https://www.beecrowd.com.br/judge/en/search?q=" + sQuery + "&for=problems"
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	document.Find("td[class='large']").Each(func(index int, mixedStr *goquery.Selection) {
		tempNum, _ := mixedStr.Find("a").Attr("href")

		tempNum = strings.TrimPrefix(tempNum, "/judge/en/problems/view/")
		tempName := mixedStr.Find("a").Text()

		var temp model.ProblemList
		temp.OJ = "URI"
		temp.PNum = tempNum
		temp.PName = tempName

		problemList = append(problemList, temp)
	})

	return problemList
}
