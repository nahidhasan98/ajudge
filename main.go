package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nahidhasan98/ajudge/apr"
	"github.com/nahidhasan98/ajudge/auth"
	"github.com/nahidhasan98/ajudge/backend"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/rank"
	"github.com/nahidhasan98/ajudge/user"
	"github.com/nahidhasan98/ajudge/vault"
	discordtexthook "github.com/nahidhasan98/discord-text-hook"
)

func logMe(r *http.Request) {
	session, _ := model.Store.Get(r, "mysession")
	username := session.Values["username"]

	msg := fmt.Sprintf("Time: %v\nIP     : %v\nIPF   : %v\nUser : %v\nURL : %v", time.Now(), r.RemoteAddr, r.Header.Get("X-FORWARDED-FOR"), username, r.RequestURI)

	// innitializing webhook
	webhook := discordtexthook.NewDiscordTextHookService(vault.WebhookIDLogger, vault.WebhookTokenLogger)

	// sending msg to discord
	_, err := webhook.SendMessage(msg)
	errorhandling.Check(err)
}

func loggingMiddlewareTemp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		go logMe(r)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	// using Gorilla mux router
	r := mux.NewRouter()

	// temporary logging for suspicious activity of the site
	r.Use(loggingMiddlewareTemp)

	auth.Init(r)
	user.Init(r)
	apr.Init(r)
	rank.Init(r)

	// for serving perspective pages
	r.HandleFunc("/", backend.Index)
	// r.PathPrefix("/rank").HandlerFunc(backend.Rank)
	r.HandleFunc("/about", backend.About)
	r.HandleFunc("/contact", backend.Contact)

	// r.HandleFunc("/register", backend.Register)
	r.PathPrefix("/verify-email/token=").HandlerFunc(backend.EmailVerifiation)

	// r.HandleFunc("/login", backend.Login)
	// r.HandleFunc("/logout", backend.Logout)
	r.PathPrefix("/profile").HandlerFunc(backend.Profile)

	r.PathPrefix("/reset").HandlerFunc(backend.Reset)
	r.PathPrefix("/passReset").HandlerFunc(backend.PassReset)

	r.HandleFunc("/problem", backend.Problem)
	r.PathPrefix("/problemView/").HandlerFunc(backend.ProblemView)
	r.PathPrefix("/origin/").HandlerFunc(backend.Origin)

	r.HandleFunc("/contest", backend.Contest)
	r.HandleFunc("/contest/create", backend.CreateContest)
	r.PathPrefix("/contestUpdate/").HandlerFunc(backend.ContestUpadte)
	r.PathPrefix("/contest/").HandlerFunc(backend.ContestGround)

	// XHR request
	r.HandleFunc("/checkLogin", backend.CheckLogin)
	r.PathPrefix("/check").HandlerFunc(backend.CheckDB)
	r.PathPrefix("/problemList").HandlerFunc(backend.ProblemList)
	r.PathPrefix("/userSubmission/").HandlerFunc(backend.GetUserSubmission)
	r.PathPrefix("/subHistory").HandlerFunc(backend.SubHistory)
	r.PathPrefix("/lang").HandlerFunc(backend.GetLanguage)
	r.HandleFunc("/submit", backend.Submit)
	r.HandleFunc("/submitC", backend.SubmitC)
	r.PathPrefix("/verdict/subID=").HandlerFunc(backend.Verdict)
	r.PathPrefix("/rejudge/subID=").HandlerFunc(backend.Rejudge)
	// r.PathPrefix("/listRank").HandlerFunc(backend.GetRankList)
	r.HandleFunc("/listContest", backend.GetContestList)
	r.PathPrefix("/problemSet/").HandlerFunc(backend.GetProblemSet)
	r.PathPrefix("/dataContest/").HandlerFunc(backend.GetContestData)
	r.PathPrefix("/captcha/").HandlerFunc(backend.GetCaptcha)
	r.HandleFunc("/getCombinedStandings", backend.GetCombinedStandings)

	// for testing any piece of code (Not essential for this site)
	r.HandleFunc("/test", backend.Test)

	// for serving images, javascripts & css files
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// A Custom Page Not Found route
	r.NotFoundHandler = http.HandlerFunc(backend.PageNotFound)

	// printing a message for displaying that local server is running
	fmt.Println("Local Server is running on port 6001...")

	// for localhost server
	http.ListenAndServe(":6001", r)
}
