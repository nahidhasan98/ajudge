package backend

import (
	"fmt"
	"net/http"

	"github.com/nahidhasan98/ajudge/mail"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	m := mail.Init()
	m.SendMailForPasswordReset("mnh.nahid35@gmail.com", "nahidhasan98", "https://ajudge.net")

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
