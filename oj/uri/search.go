package uri

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
)

//Search function for searching problem in URI
func Search(sQuery string) []byte {
	defer errorhandling.Recovery() //for panic() error Recovery

	//defining a variable for returning data
	var res []byte

	//data that will be collected
	type list struct {
		Num, Name string
	}
	var problemList []list

	//two separate lists(number & name) of problem. will be entered in final problem list later
	var probNumList, probNameList []string
	var tempNum, tempName string

	for k := 1; k <= 10; k++ { //taking 10 pages problem
		var apiURL string

		if sQuery == "URI Only" {
			apiURL = "https://www.urionlinejudge.com.br/judge/en/search?q=&page=" + strconv.Itoa(k)
		} else { //searching with a problem name/number
			sQuery = strings.ReplaceAll(sQuery, ` `, `+`) //making a valid query like: leap year => leap+year
			apiURL = "https://www.urionlinejudge.com.br/judge/en/search?q=" + sQuery + "&for=problems"
		}

		response := GETRequest(apiURL)
		defer response.Body.Close()
		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		//collecting a list of problem-numbers
		document.Find("td[class='large']").Each(func(index int, mixedStr *goquery.Selection) {
			tempNum, _ = mixedStr.Find("a").Attr("href")
			tempNum = strings.TrimPrefix(tempNum, "/judge/en/problems/view/")
			probNumList = append(probNumList, tempNum)
		})

		//collecting a list of problem-names
		document.Find("td[class='large']").Each(func(index int, mixedStr *goquery.Selection) {
			tempName = mixedStr.Find("a").Text()
			probNameList = append(probNameList, tempName)
		})

		if sQuery != "URI Only" { //if specific problem is searched then only 1 time search
			break
		}
	}
	//preparing final list from two separate lists
	for i := 0; i < len(probNumList); i++ {
		temp := list{probNumList[i], probNameList[i]}
		problemList = append(problemList, temp)
	}
	res, _ = json.Marshal(problemList) ///making json file

	return res
}
