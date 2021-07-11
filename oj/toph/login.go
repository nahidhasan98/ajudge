package toph

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
)

//Login function for login to Toph
func Login() string {
	defer errorhandling.Recovery() //for panic() error Recovery

	apiURL := "https://toph.co"
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	//finding loggin info
	resp := document.Find("a[href='/u/ajudge.bd']").Text()
	resp = strings.TrimSpace(resp)

	//getting for later usesage (for submit)
	getToken(document)

	if resp != "A Judgeajudge.bdProfile" { //if not already logged in - then try to login
		apiURL := "https://toph.co/login"
		postData := url.Values{
			"handle":   {vault.TophUsername},
			"password": {vault.TophPassword},
		}

		req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		resp := document.Find("a[href='/u/ajudge.bd']").Text()
		resp = strings.TrimSpace(resp)
		//fmt.Println("Login=", resp)

		//getting for later usesage (for submit)
		getToken(document)

		if resp != "A Judgeajudge.bdProfile" { //successful login page has somthing like this
			return "failed"
		}
	}
	//already logged in
	return "success"
}

func getToken(document *goquery.Document) {
	result := document.Find("script").Text()
	result = strings.TrimSpace(result)

	need := "tokenId"
	index := strings.Index(result, need)

	tk := result[index+9 : index+41] //token length is 32
	tokenID = "Token " + tk

	//fmt.Println("Token=", tokenID)
}
