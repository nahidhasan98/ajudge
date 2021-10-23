package recaptcha

import (
	"net/http"
	"strings"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/vault"
)

func sendAPIRequest(captcha string) (*http.Response, error) {
	method := "POST"
	apiURL := "https://www.google.com/recaptcha/api/siteverify"

	// preparing payload
	payload := strings.NewReader("secret=" + vault.CaptchaKey + "&response=" + captcha)

	// setting up request
	req, err := http.NewRequest(method, apiURL, payload)
	errorhandling.Check(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	return client.Do(req)
}
