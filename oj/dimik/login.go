package dimik

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
)

//Login function for login to DimikOJ
func Login() string {
	defer errorhandling.Recovery() //for panic() error errorhandling.Recovery

	//checking if already logged in or not
	apiURL := "https://dimikoj.com"
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	//finding loggin info
	resp := document.Find("span[class='badge badge-secondary']").Text()
	hiddenToken, _ = document.Find("input[name='_token']").Attr("value") //will be used for submission

	if resp != "♠" { //if not already logged in - then try to login
		apiURL := "https://dimikoj.com/login"
		response := GETRequest(apiURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		hiddenToken, _ = document.Find("input[name='_token']").Attr("value")

		postData := url.Values{
			"_token":   {hiddenToken},
			"email":    {vault.DimikOJUsername},
			"password": {vault.DimikOJPassword},
		}

		req, _ := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		response, err = model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err = goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		resp := document.Find("span[class='badge badge-secondary']").Text()
		hiddenToken, _ = document.Find("input[name='_token']").Attr("value") //will be used for submission

		if resp != "♠" { //successful login page has somthing like this
			return "failed"
		}
	}
	//already logged in
	return "success"
}
