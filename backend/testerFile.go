package backend

import (
	"fmt"
	"net/http"
)

// Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// out, err := exec.Command("go", "env").CombinedOutput()
	// if err != nil {
	// 	fmt.Println("16", err)
	// }
	// fmt.Println("18#", string(out))

	// out, err = exec.Command("echo", os.Getenv("PATH")).CombinedOutput()
	// if err != nil {
	// 	fmt.Println("22", err)
	// }
	// fmt.Println("24#", string(out))

	// out, err = exec.Command("git", "version").CombinedOutput()
	// if err != nil {
	// 	fmt.Println("28", err)
	// }
	// fmt.Println("30#", string(out))

	fmt.Fprintln(w, "Hello from test")

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}

// package dimik

// import (
// 	"errors"
// 	"net/http"
// 	"net/url"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// 	"github.com/nahidhasan98/ajudge/errorhandling"
// 	"github.com/nahidhasan98/ajudge/model"
// 	"github.com/nahidhasan98/ajudge/vault"
// )

// type dimik struct {
// 	// Username, Password string
// }

// type OJInterfacer interface {
// 	Login() error
// 	Submission(pNum, lang, src string) error
// }

// func NewOJ() OJInterfacer {
// 	return &dimik{}
// }

// func (oj *dimik) Login() error {
// 	defer errorhandling.Recovery() //for panic() error errorhandling.Recovery

// 	//checking if already logged in or not
// 	apiURL := "https://dimikoj.com/login"
// 	response := GETRequest(apiURL)
// 	defer response.Body.Close()

// 	document, err := goquery.NewDocumentFromReader(response.Body)
// 	errorhandling.Check(err)

// 	//finding loggin info
// 	resp, _ := document.Find("#logout-form").Html()
// 	resp = strings.TrimSpace(resp)

// 	//if already logged in - then return
// 	if len(resp) > 0 {
// 		return nil
// 	}

// 	// trying to login now
// 	hiddenToken, _ := document.Find("input[name='_token']").Attr("value")

// 	postData := url.Values{
// 		"_token":   {hiddenToken},
// 		"email":    {vault.DimikOJUsername},
// 		"password": {vault.DimikOJPassword},
// 	}

// 	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(postData.Encode()))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
// 	response, err = model.Client.Do(req)
// 	errorhandling.Check(err)
// 	defer response.Body.Close()

// 	document, err = goquery.NewDocumentFromReader(response.Body)
// 	errorhandling.Check(err)

// 	resp, _ = document.Find("#logout-form").Html()
// 	resp = strings.TrimSpace(resp)

// 	if len(resp) > 0 {
// 		return nil
// 	}

// 	return errors.New("login error")
// }

// func (oj *dimik) Submission(pNum, lang, src string) error {
// 	return nil
// }
