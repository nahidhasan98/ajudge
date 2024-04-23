package user

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
	"github.com/nahidhasan98/ajudge/recaptcha"
)

type handler struct {
	userService userInterfacer
}

// makeHTTPHandlers function defines endpoints
func makeHTTPHandlers(router *mux.Router, userService userInterfacer) {
	h := &handler{
		userService: userService,
	}

	router.HandleFunc("/register", h.previewRegistrationPage).Methods("GET")
	router.HandleFunc("/register", h.register).Methods("POST")
}

// getLoginPage function for login to our own site
func (h *handler) previewRegistrationPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		model.Info["PageName"] = "Registration"
		model.Info["PageTitle"] = "Registration | AJudge"

		model.Tpl.ExecuteTemplate(w, "register.gohtml", model.Info)
	}
}

type registrationResponse struct {
	Status      string   `json:"status"`
	RedirectURL string   `json:"redirectURL"`
	Errors      []Errors `json:"errors"`
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	// gettting form data
	fullName := html.EscapeString(strings.TrimSpace(r.FormValue("fullName")))
	email := html.EscapeString(strings.TrimSpace(r.FormValue("email")))
	username := html.EscapeString(strings.TrimSpace(r.FormValue("username")))
	password := html.EscapeString(r.FormValue("password"))
	confirmPassword := html.EscapeString(r.FormValue("confirmPassword"))
	captcha := html.EscapeString(r.FormValue("g-recaptcha-response"))

	// validating form data
	errs := validate(fullName, email, username, password, confirmPassword, captcha, h)
	if len(errs) > 0 {
		response := registrationResponse{
			Status: "error",
			Errors: errs,
		}
		byteData, err := json.Marshal(response)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(byteData)
		return
	}

	userData, err := h.userService.AddUser(fullName, email, username, password)
	if err != nil { // if registration failed
		response := registrationResponse{
			Status: "error",
			Errors: []Errors{
				{Type: "other", Message: fmt.Sprintf("%v", err)},
			},
		}
		byteData, err := json.Marshal(response)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(byteData)
		return
	}

	model.PopUpCause = "registrationDone" // login page will give a popup

	// preparing data to response
	response := registrationResponse{
		Status:      "success",
		RedirectURL: "/login",
	}
	b, err := json.Marshal(response)
	errorhandling.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	// notifying to discord
	disData := discord.UserModel(*userData)
	dis := discord.Init()
	dis.SendMessage(disData, "registration")
}

type Errors struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func validate(fullName, email, username, password, confirmPassword, captcha string, h *handler) []Errors {
	var errs []Errors

	// taking care of fullName
	if len(fullName) == 0 {
		temp := Errors{
			Type:    "fullName",
			Message: "name should no be empty",
		}
		errs = append(errs, temp)
	}

	// taking care of email
	res := h.userService.isAvailableEmail(email)
	if !res {
		temp := Errors{
			Type:    "email",
			Message: "email already registered, choose another one",
		}
		errs = append(errs, temp)
	}

	// taking care of username
	// PART 1: length
	if len(username) == 0 {
		temp := Errors{
			Type:    "username",
			Message: "username should no be empty",
		}
		errs = append(errs, temp)
	} else {
		// PART 2: space
		index := strings.Index(username, " ")
		if index > -1 { // contains space
			temp := Errors{
				Type:    "username",
				Message: "username can't contains space",
			}
			errs = append(errs, temp)
		} else {
			// PART 3: availability
			res = h.userService.isAvailableUsername(username)
			if !res { // no error means a user found by this username, so username not available
				temp := Errors{
					Type:    "username",
					Message: "username already taken, choose another one",
				}
				errs = append(errs, temp)
			}
		}
	}

	// taking care of password
	// PART 1: length
	if len(password) < 8 {
		temp := Errors{
			Type:    "password",
			Message: "password length should be at least 8 characters",
		}
		errs = append(errs, temp)
	} else {
		// PART 2: matching
		if password != confirmPassword {
			temp := Errors{
				Type:    "password",
				Message: "password mismatched, put cautiously",
			}
			errs = append(errs, temp)
		}
	}

	// taking care of recaptcha
	//
	// google recaptcha can only be verified once to prevent replay attacks
	// so captcha verification process will be done if all other fields are valid/okay.
	if len(errs) == 0 {
		rs := recaptcha.NewRecaptchaService()
		err := rs.ValidateCaptcha(captcha)
		if err != nil {
			temp := Errors{
				Type:    "captcha",
				Message: "captcha error, please verify you are not a robot",
			}
			errs = append(errs, temp)
		}
	}

	return errs
}
