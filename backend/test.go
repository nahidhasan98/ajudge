package backend

import (
	"fmt"
	"net/http"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}

/*
	//searching problem from DimikOJ first
	var pRes []byte
	if OJ == "" || OJ == "All" || OJ == "DimikOJ" {
		if pNum != "" {
			pRes = dimik.Search(pNum)
		} else if pName != "" {
			pRes = dimik.Search(pName)
		} else {
			if OJ == "DimikOJ" { //again checking cause it may contains only DimikOJ or All OJ
				pRes = dimik.Search("DimikOJ Only")
			} else {
				rand.Seed(time.Now().UnixNano())
				min := 1   //the first problem DimikOJ is 1
				max := 100 //the last problem DimikOJ is approx. 100
				sQuery := rand.Intn(max-min+1) + min
				pRes = dimik.Search(strconv.Itoa(sQuery))
			}
		}
		type PResList struct {
			Num  string `json:"Num"`
			Name string `json:"Name"`
		}
		var pResList []PResList
		json.Unmarshal(pRes, &pResList)

		for i := 0; i < len(pResList); i++ {
			var temp model.List

			temp.OJ = "DimikOJ"
			temp.Num = pResList[i].Num
			temp.Name = pResList[i].Name

			pListFinal = append(pListFinal, temp)
		}
	}
	if OJ == "" || OJ == "All" || OJ == "Toph" {
		if pNum != "" {
			pRes = toph.Search(pNum)
		} else if pName != "" {
			pRes = toph.Search(pName)
		} else {
			if OJ == "Toph" { //again checking cause it may contains only DimikOJ or All OJ
				pRes = toph.Search("Toph Only")
			} else { //actually this do nothing because Toph has no problem number
				rand.Seed(time.Now().UnixNano())
				min := 1   //the first problem DimikOJ is 1
				max := 100 //the last problem DimikOJ is approx. 100
				sQuery := rand.Intn(max-min+1) + min
				pRes = toph.Search(strconv.Itoa(sQuery))
			}
		}
		type PResList struct {
			Num  string `json:"Num"`
			Name string `json:"Name"`
		}
		var pResList []PResList
		json.Unmarshal(pRes, &pResList)

		for i := 0; i < len(pResList); i++ {
			var temp model.List

			temp.OJ = "Toph"
			temp.Num = pResList[i].Num
			temp.Name = pResList[i].Name

			pListFinal = append(pListFinal, temp)
		}
	}
	if OJ == "" || OJ == "All" || OJ == "URI" {
		if pNum != "" {
			pRes = uri.Search(pNum)
		} else if pName != "" {
			pRes = uri.Search(pName)
		} else {
			if OJ == "URI" {
				pRes = uri.Search("URI Only")
			} else {
				rand.Seed(time.Now().UnixNano())
				min := 1001 //the first problem in URI is 1001
				max := 3100 //the last problem in URI is approx. 3100
				sQuery := rand.Intn(max-min+1) + min
				pRes = uri.Search(strconv.Itoa(sQuery))
			}
		}
		type PResList struct {
			Num  string `json:"Num"`
			Name string `json:"Name"`
		}
		var pResList []PResList
		json.Unmarshal(pRes, &pResList)

		for i := 0; i < len(pResList); i++ {
			var temp model.List

			temp.OJ = "URI"
			temp.Num = pResList[i].Num
			temp.Name = pResList[i].Name

			pListFinal = append(pListFinal, temp)
		}
	}
	if OJ != "DimikOJ" && OJ != "Toph" && OJ != "URI" {
		probData := vjudge.Search(OJ, pNum, pName)

		var temp model.List
		for i := 0; i < len(probData); i++ {
			temp.OJ = probData[i].OJ
			temp.Num = probData[i].Num
			temp.Name = probData[i].Name

			pListFinal = append(pListFinal, temp)
		}
	}

*/
