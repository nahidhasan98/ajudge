package backend

import (
	"fmt"
	"net/http"

	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
	discordtexthook "github.com/nahidhasan98/discord-text-hook"
)

func init() {

}

//Test function for testing a piece of code
func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// hand("Login Done")
	disData := model.UserData{
		Username: "Nahid98",
	}
	discord := discord.Init()
	discord.SendMessage(disData, "login")

	msg := "```md"
	msg += `content:nahid98`
	msg += "```"
	fmt.Println(msg)
	hand(msg)

	fmt.Println("ENDDDDDD")
	fmt.Println("Happy coding.")
	//model.Tpl.ExecuteTemplate(w, "test.html", nil)
}

func hand(msg string) {
	fmt.Println(27)
	jobs := make(chan int, 5)
	go send(jobs, msg)
	jobs <- 1
	close(jobs)
}

func send(jobs <-chan int, msg string) {
	fmt.Println(35)
	for range jobs {
		ds := discordtexthook.NewDiscordTextHookService(vault.WebhookIDLogin, vault.WebhookTokenLogin)
		h, e := ds.SendMessage(msg)
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(h)
	}
}
