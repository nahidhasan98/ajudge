package recaptcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/nahidhasan98/ajudge/errorhandling"
)

type recaptchaInterfacer interface {
	ValidateCaptcha(captcha string) error
}

type recaptcha struct{}

type recaptchaResponse struct {
	Success bool          `json:"success"`
	Codes   []interface{} `json:"error-codes"`
}

//GetCaptcha function for confirming captcha correct or not
func (rc *recaptcha) ValidateCaptcha(captcha string) error {
	response, err := sendAPIRequest(captcha)
	errorhandling.Check(err)
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	errorhandling.Check(err)

	var capRes recaptchaResponse
	json.Unmarshal(resBody, &capRes)

	if !capRes.Success {
		e := fmt.Sprintf("%v", capRes.Codes[0])
		errorhandling.Check(errors.New("Error from google recaptcha verification: " + e))
		return errors.New(e)
	}

	return nil
}

func NewRecaptchaService() recaptchaInterfacer {
	return &recaptcha{}
}
