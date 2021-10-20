package auth

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

// type handlerInterface interface {
// 	loginHandler(w http.ResponseWriter, r *http.Request)
// }

type handler struct {
	authService authInterfacer
}

// getLoginPage function for login to our own site
func (h *handler) getLoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// preparing info for sending to frontend
		model.Info["PageName"] = "Login"
		model.Info["PageTitle"] = "Login | AJudge"
		model.Info["LastPage"] = model.LastPage
		model.Info["PopUpCause"] = model.PopUpCause

		model.Tpl.ExecuteTemplate(w, "login.gohtml", model.Info)

		// clearing up variables
		model.Info["Username"] = ""
		model.Info["ErrUsername"] = ""
		model.Info["ErrPassword"] = ""
		model.PopUpCause = ""
		model.Info["PopUpCause"] = model.PopUpCause
	}
}

type loginResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	RedirectURL string `json:"redirectURL"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	// getting form data
	username := html.EscapeString(strings.TrimSpace(r.FormValue("username")))
	password := html.EscapeString(r.FormValue("password"))

	userData, err := h.authService.Authenticate(username, password)
	if err != nil { // if login credentials failed
		response := loginResponse{
			Status:  "error",
			Message: fmt.Sprintf("%v", err),
		}
		byteData, err := json.Marshal(response)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(byteData)
		return
	}

	// saving session data
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)
	session.Values["username"] = username
	session.Values["isLogin"] = true
	session.Save(r, w)

	// preparing info for frontend
	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]

	// preparing data to response
	response := loginResponse{
		Status:      "success",
		RedirectURL: model.LastPage,
	}
	byteDate, err := json.Marshal(response)
	errorhandling.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteDate)

	// notifying to discord
	disData := *userData
	discord := discord.Init()
	discord.SendMessage(disData, "login")
}

// logout function for logging out from our own site
func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)
	session.Values["username"] = ""
	session.Values["isLogin"] = false

	session.Options.MaxAge = -1 // cookies will be deleted immediately
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// makeHTTPHandlers function defines endpoints
func makeHTTPHandlers(router *mux.Router, authService authInterfacer) {
	h := &handler{
		authService: authService,
	}

	router.HandleFunc("/login", h.getLoginPage).Methods("GET")
	router.HandleFunc("/login", h.login).Methods("POST")
	router.HandleFunc("/logout", h.logout).Methods("GET")
}
