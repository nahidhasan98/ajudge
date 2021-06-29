package vjudge

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

//Search function for searching problem in VJ
func Search(OJ, pNum, pName string) []model.ProblemList {
	defer errorhandling.Recovery() //for panic() error Recovery

	//defining a variable for returning data
	var problemList []model.ProblemList

	if OJ == "计蒜客" || OJ == "黑暗爆炸" {
		OJ = url.QueryEscape(OJ)
	}
	var length = "1000" //we will take about 1000 problem from VJudge

	pName = strings.ReplaceAll(pName, ` `, `%20`) //making a valid query like: leap year => leap%20year
	apiURL := "https://vjudge.net/problem/data?draw=1&start=0&length=" + length + "&sortDir=desc&sortCol=5&OJId=" + OJ + "&probNum=" + pNum + "&title=" + pName + "&source=&category=all"
	response := GETRequest(apiURL)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	//VJudge returns data in json format.
	//json file contains a structure
	//structure contains an array & some other integer
	//the elements of this array contains another structure
	type list struct {
		OJ    string `json:"originOJ"`
		PNum  string `json:"originProb"`
		PName string `json:"title"`
	}
	type searchResult struct {
		Data []list `json:"data"`
	}
	var respBody searchResult
	json.Unmarshal(body, &respBody) //extracting the json file

	for i := 0; i < len(respBody.Data); i++ { //getting problem one by one
		if model.OJSet[respBody.Data[i].OJ] { //if problem come from desired OJ
			var temp model.ProblemList
			temp.OJ = respBody.Data[i].OJ
			temp.PNum = respBody.Data[i].PNum
			temp.PName = respBody.Data[i].PName

			problemList = append(problemList, temp)
		}
	}

	return problemList
}
