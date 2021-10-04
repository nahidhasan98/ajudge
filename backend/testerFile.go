package backend

import (
	"fmt"
	"net/http"

	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/model"
)

func init() {

}

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ds := discord.Init()
	l := model.UserData{
		Username: "nahidhasan98",
	}
	ds.SendMessage(l, "login")

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
