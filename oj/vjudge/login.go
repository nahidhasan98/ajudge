package vjudge

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
)

//Login function for login to VJ
func Login() string {
	defer errorhandling.Recovery() //for panic() error Recovery

	//checking if already logged in or not
	apiURL := "https://vjudge.net"
	response := GETRequest(apiURL)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	errorhandling.Check(err)

	resp := document.Find("a[id='userNameDropdown']").Text()
	resp = strings.TrimSpace(resp)

	if resp != "ajudgebd" { //if not already logged in - then try to login
		apiURL := "https://vjudge.net/user/login"
		postData := url.Values{
			"username": {vault.VJudgeUsername},
			"password": {vault.VJudgePassword},
		}

		req, err := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
		errorhandling.Check(err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		response, err := model.Client.Do(req)
		errorhandling.Check(err)
		defer response.Body.Close()

		resp, err := ioutil.ReadAll(response.Body)
		errorhandling.Check(err)

		if string(resp) != "success" { //successful login returns this text
			return "failed"
		}
	}
	//already logged in
	return "success"
}
