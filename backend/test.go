package backend

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nahidhasan98/ajudge/oj/toph"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	apiURL := "https://toph.co/problems/all?start=1"
	response := toph.GETRequest(apiURL)
	fmt.Println(response)
	defer response.Body.Close()
	// _, err := goquery.NewDocumentFromReader(response.Body)
	// errorhandling.Check(err)

	resp, _ := ioutil.ReadAll(response.Body)
	body := string(resp)
	fmt.Println("Hellooooo", body)

	fmt.Println("ENDDDDDD")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
