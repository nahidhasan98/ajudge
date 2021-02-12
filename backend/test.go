package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nahidhasan98/ajudge/errorhandling"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	contestDateTime := "324242:00"
	_, err := time.Parse(time.RFC3339, contestDateTime) //format RFC3339 = "2006-01-02T15:04:05Z07:00"
	errorhandling.Check(err)

	//not found: !=nil - mongo: no document in result
	fmt.Println("ENDDDDDD")
	//tpl.ExecuteTemplate(w, "test.html", Info)
}

//URISearch function for searching problem in URI
// func collectProbList(w http.ResponseWriter) []byte {
// 	defer recovery() //for panic() error recovery

// 	//defining a variable for returning data
// 	var res []byte

// 	//data that will be collected
// 	type list struct {
// 		Num, Name string
// 	}
// 	var problemList []list

// 	//two separate lists(number & name) of problem. will be entered in final problem list later
// 	var probNumList, probNameList []string
// 	var tempNum, tempName string
// 	var uriPageNumList []int

// 	fmt.Println(time.Now())
// 	for k := 1; k <= 85; k++ { //taking 10 pages problem
// 		var apiURL string

// 		apiURL = "https://www.urionlinejudge.com.br/judge/en/problems/all?page=" + strconv.Itoa(k)

// 		response := URIGETRequest(apiURL)
// 		defer response.Body.Close()
// 		document, err := goquery.NewDocumentFromReader(response.Body)
// 		CheckErr(err)

// 		//collecting a list of problem numbers
// 		document.Find("td[class='large']").Each(func(index int, mixedStr *goquery.Selection) {
// 			tempNum, _ = mixedStr.Find("a").Attr("href")

// 			need := "/judge/en/problems/view/"
// 			match := strings.Index(tempNum, need)

// 			if match != -1 {
// 				tempNum = strings.TrimPrefix(tempNum, "/judge/en/problems/view/")
// 				probNumList = append(probNumList, tempNum)

// 				tempName = mixedStr.Find("a").Text()
// 				probNameList = append(probNameList, tempName)

// 				uriPageNumList = append(uriPageNumList, k)
// 			}
// 		})
// 	}
// 	fmt.Println(time.Now())
// 	DB := dbConn()
// 	defer DB.Close()

// 	//preparing final list from two separate lists
// 	for i := 0; i < len(probNumList); i++ {
// 		temp := list{probNumList[i], probNameList[i]}
// 		problemList = append(problemList, temp)

// 		fmt.Fprintln(w, probNumList[i], probNameList[i])

// 		var id int
// 		exist := DB.QueryRow("SELECT id FROM uri WHERE pNum=?", probNumList[i]).Scan(&id)

// 		if exist == sql.ErrNoRows {
// 			insertQuery, err := DB.Prepare("INSERT INTO uri (pNum,pName,uriPageNum) VALUES (?,?,?)")
// 			CheckErr(err)
// 			insertQuery.Exec(probNumList[i], probNameList[i], uriPageNumList[i])
// 		}
// 	}
// 	fmt.Println(time.Now())
// 	res, _ = json.Marshal(problemList) ///making json file

// 	return res
// }

// subIDs := [2]int{4, 6}
// //setting uo options for retrieving data from DB
// opts := options.Find()
// opts.SetSort(bson.D{{Key: "subID", Value: -1}}) //descending order
// cursor, err := submissionCollection.Find(ctx, bson.M{"subID": bson.M{"$in": subIDs}}, opts)
// CheckErr(err)
