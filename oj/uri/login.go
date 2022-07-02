package uri

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
)

//Login function for login to URI
func Login() string {
	defer errorhandling.Recovery() //for panic() error Recovery

	//checking if already logged in or not
	apiURL := "https://www.beecrowd.com.br/judge/en"
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	//finding loggin info
	resp := document.Find("div[class='h-user']").Find("i").Text()

	if resp != "ajudge.bd@gmail.com" { //if not already logged in - then try to login
		apiURL := "https://www.beecrowd.com.br/judge/en/login"
		response := GETRequest(apiURL)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		//finding data for POST request
		method, _ = document.Find("input[name='_method']").Attr("value")
		csrfToken, _ = document.Find("input[name='_csrfToken']").Attr("value")
		tokenFields, _ = document.Find("input[name='_Token[fields]']").Attr("value")
		tokenUnlocked, _ = document.Find("input[name='_Token[unlocked]']").Attr("value")
		//fmt.Println(method, csrfToken, tokenFields, tokenUnlocked)

		postData := url.Values{
			"_method":          {method},
			"_csrfToken":       {csrfToken},
			"email":            {vault.URIUsername},
			"password":         {vault.URIPassword},
			"remember_me":      {"0"},
			"_Token[fields]":   {tokenFields},
			"_Token[unlocked]": {tokenUnlocked},
		}

		req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		response, err = model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		document, err = goquery.NewDocumentFromReader(response.Body)
		errorhandling.Check(err)

		resp := document.Find("div[class='h-user']").Find("i").Text()

		if resp != "ajudge.bd@gmail.com" { //successful login page has somthing like this
			return "failed"
		}
	}
	//already logged in
	return "success"
}
