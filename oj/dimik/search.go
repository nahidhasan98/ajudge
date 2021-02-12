package dimik

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
)

//Search function for searching problem in DimikOJ
func Search(sQuery string) []byte {
	defer errorhandling.Recovery() //for panic() error errorhandling.Recovery

	//defining a variable for returning data
	var res []byte

	//data that will be collected
	type list struct {
		Num, Name string
	}
	var problemList []list

	//two separate lists(number & name) of problem. will be entered in final problem list later
	var probNumList []string
	var tempNum, tempName string

	if sQuery == "DimikOJ Only" {
		for k := 1; k <= 4; k++ { //there are about 4 pages problem in DimikOJ
			//Getting Problem List
			apiURL := "https://dimikoj.com/problems?page=" + strconv.Itoa(k)
			response := GETRequest(apiURL)
			defer response.Body.Close()

			document, err := goquery.NewDocumentFromReader(response.Body)
			errorhandling.Check(err)

			document.Find("tr").Find("td").Each(func(index int, mixedStr *goquery.Selection) {
				tempNum = mixedStr.Text()
				probNumList = append(probNumList, tempNum)
			})
		}
		for i := 0; i < len(probNumList); i++ {
			if i%5 == 0 {
				tempNum = probNumList[i]    //prob number index like:0,5,10,15...
				tempName = probNumList[i+1] //prob Name index like:1,6,11,16...

				temp := list{tempNum, tempName}
				problemList = append(problemList, temp)
			}
		}
	} else { //Getting Only One Problem
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
			temp := list{tempNum, tempName}
			problemList = append(problemList, temp)
		}
	}
	res, _ = json.Marshal(problemList) //making json file

	return res
}
