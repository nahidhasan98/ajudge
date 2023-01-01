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

// Login function for login to DimikOJ
func Login() string {
	defer errorhandling.Recovery() //for panic() error errorhandling.Recovery

	//checking if already logged in or not
	apiURL := "https://dimikoj.com/login"
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	//finding loggin info
	resp, _ := document.Find("#logout-form").Html()
	resp = strings.TrimSpace(resp)

	// if already logged in then this will be used in submission
	// but if not logged in then this will be used in new login
	hiddenToken, _ = document.Find("input[name='_token']").Attr("value")

	//if already logged in - then return
	if len(resp) > 0 {
		return "success"
	}

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

	resp, _ = document.Find("#logout-form").Html()
	resp = strings.TrimSpace(resp)

	if len(resp) > 0 {
		// after login, new token is generated
		// and this will be used in submission
		hiddenToken, _ = document.Find("input[name='_token']").Attr("value")
		return "success"
	}

	return "failed"
}
