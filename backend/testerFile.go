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

func worker(jobs <-chan int) {
	fmt.Println("Register the worker")
	for i := range jobs {
		fmt.Println("worker processing job", i)
		d1 := discord.Init()

		t := model.SubmissionData{
			SubID:    123,
			Username: "Nahid",
		}

		d1.SendMessage(t, "submission")
		fmt.Println("Msg sent")

		fmt.Println("worker exiting job", i)
		//time.Sleep(time.Second * 5)
	}
}
