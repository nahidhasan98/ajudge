package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/model"
)

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	d1 := discord.Init()

	t := model.SubmissionData{
		SubID:    123,
		Username: "Nahid",
	}

	_, _ = d1.SendMessage(t)
	fmt.Println("Msg sent")

	time.Sleep(3 * time.Second)

	d2 := discord.Init()
	t = model.SubmissionData{
		SubID:    123,
		Username: "Hasan",
	}

	d2.EditMessage(t)

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}
