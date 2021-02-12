package toph

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
)

//Search function for searching problem in Toph
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

	if sQuery == "Toph Only" {
		apiURL := "https://toph.co/problems"
		response := GETRequest(apiURL)
		defer response.Body.Close()
		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		var categoryList []string
		var categoryName string

		//finding all category list
		document.Find("td").Each(func(index int, mixedStr *goquery.Selection) {
			categoryName, _ = mixedStr.Find("a").Attr("href") //Toph have no pNum. So we taking problemLink (like:/p/copycat)
			categoryList = append(categoryList, categoryName)
		})

		for k := 1; k <= 6; k++ { //taking 6 category problem
			//selecting a random category
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := 92 //there are about 31 categories (each categories have three <td>)
			category := rand.Intn(max-min+1) + min

			sQuery := categoryList[category]

			apiURL := "https://toph.co" + sQuery
			response := GETRequest(apiURL)
			defer response.Body.Close()
			document, err = goquery.NewDocumentFromReader(response.Body)
			errorhandling.Check(err)

			document.Find("div[class='col-md-9']").Find("td").Each(func(index int, mixedStr *goquery.Selection) {
				tempNum, _ = mixedStr.Find("a").Attr("href") //Toph have no pNum. So we taking pLink (like:/p/copycat)

				tempNum = tempNum[3:] //avoiding /p/ from the link
				probNumList = append(probNumList, tempNum)

				tempName = mixedStr.Find("h4").Text()
				probNameList = append(probNameList, tempName)
			})
		}
	} else { //searching with a problem name
		sQuery = strings.ReplaceAll(sQuery, ` `, `+`) //making a valid query like: leap year => leap+year
		apiURL := "https://toph.co/search?type=problems&q=" + sQuery
		response := GETRequest(apiURL)
		defer response.Body.Close()
		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		document.Find("td").Each(func(index int, mixedStr *goquery.Selection) {
			tempNum, _ = mixedStr.Find("a").Attr("href") //Toph have no pNum. So we taking pLink (like:/p/copycat)

			tempNum = tempNum[3:] //avoiding /p/ from the link
			probNumList = append(probNumList, tempNum)

			tempName = mixedStr.Find("h4").Text()
			probNameList = append(probNameList, tempName)
		})
	}
	//preparing final list from two separate lists
	for i := 0; i < len(probNumList); i++ {
		temp := list{probNumList[i], probNameList[i]}
		problemList = append(problemList, temp)
	}
	res, _ = json.Marshal(problemList) //making json file

	return res
}
