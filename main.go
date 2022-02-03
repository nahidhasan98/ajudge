package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nahidhasan98/ajudge/apr"
	"github.com/nahidhasan98/ajudge/auth"
	"github.com/nahidhasan98/ajudge/backend"
	"github.com/nahidhasan98/ajudge/user"
)

func main() {
	//(instead of default 'http' router) using Gorilla mux router
	r := mux.NewRouter()
	auth.Init(r)
	user.Init(r)
	apr.Init(r)

	//just a message for ensuring that local server is running
	fmt.Println("Local Server is running...")

	// for serving perspective pages
	r.HandleFunc("/", backend.Index)
	r.PathPrefix("/rank").HandlerFunc(backend.Rank)
	r.HandleFunc("/about", backend.About)
	r.HandleFunc("/contact", backend.Contact)

	// r.HandleFunc("/register", backend.Register)
	r.PathPrefix("/verify-email/token=").HandlerFunc(backend.EmailVerifiation)

	// r.HandleFunc("/login", backend.Login)
	// r.HandleFunc("/logout", backend.Logout)
	r.PathPrefix("/profile").HandlerFunc(backend.Profile)

	r.PathPrefix("/reset").HandlerFunc(backend.Reset)
	r.PathPrefix("/passReset/token=").HandlerFunc(backend.PassReset)

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
	r.PathPrefix("/listRank").HandlerFunc(backend.GetRankList)
	r.HandleFunc("/listContest", backend.GetContestList)
	r.PathPrefix("/problemSet/").HandlerFunc(backend.GetProblemSet)
	r.PathPrefix("/dataContest/").HandlerFunc(backend.GetContestData)
	// r.PathPrefix("/captcha/").HandlerFunc(backend.GetCaptcha)
	r.HandleFunc("/getCombinedStandings", backend.GetCombinedStandings)

	// for testing any piece of code (Not essential for this site)
	r.HandleFunc("/test", backend.Test)

	// for serving images, javascripts & css files
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// A Custom Page Not Found route
	r.NotFoundHandler = http.HandlerFunc(backend.PageNotFound)

	// for localhost server
	http.ListenAndServe(":8080", r)
}
