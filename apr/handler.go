package apr

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

type handler struct {
	Tpl        *template.Template
	AprService aprInterfacer
}

// makeHTTPHandlers function defines endpoints
func makeHTTPHandlers(router *mux.Router, aprService aprInterfacer) {
	h := &handler{
		Tpl:        template.Must(template.ParseGlob("frontend/html/*")),
		AprService: aprService,
	}

	router.HandleFunc("/apr", h.previewApr).Methods("GET")
	router.HandleFunc("/apr/pull", h.pull).Methods("POST")
	router.HandleFunc("/apr/build", h.build).Methods("POST")
	router.HandleFunc("/apr/restart", h.restart).Methods("POST")
}

// logout function for logging out from our own site
func (h *handler) previewApr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	if session.Values["isLogin"] == false || model.Info["Username"] != "admin" {
		//errorPage(w, http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound) //status code such as: 400, 404 etc.
		model.Info["StatusCode"] = http.StatusNotFound
		h.Tpl.ExecuteTemplate(w, "pageNotFound.gohtml", model.Info)
		return
	}

	model.Info["PageName"] = "apr"
	model.Info["PageTitle"] = "Auto Pull Restart | AJudge"

	h.Tpl.ExecuteTemplate(w, "apr.html", model.Info)
}

type aprResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// logout function for logging out from our own site
func (h *handler) pull(w http.ResponseWriter, r *http.Request) {
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	if session.Values["isLogin"] == false || model.Info["Username"] != "admin" {
		res := aprResponse{
			Status:  "error",
			Message: "access denied",
		}

		b, err := json.Marshal(res)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	out, err := h.AprService.pull()
	if err != nil {
		res := aprResponse{
			Status:  "error",
			Message: err.Error(),
		}

		b, err := json.Marshal(res)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	res := aprResponse{
		Status:  "success",
		Message: string(out),
	}
	b, err := json.Marshal(res)
	errorhandling.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// logout function for logging out from our own site
func (h *handler) build(w http.ResponseWriter, r *http.Request) {
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	if session.Values["isLogin"] == false || model.Info["Username"] != "admin" {
		res := aprResponse{
			Status:  "error",
			Message: "access denied",
		}

		b, err := json.Marshal(res)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	out, err := h.AprService.build()
	if err != nil {
		res := aprResponse{
			Status:  "error",
			Message: err.Error(),
		}

		b, err := json.Marshal(res)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	res := aprResponse{
		Status:  "success",
		Message: string(out),
	}
	b, err := json.Marshal(res)
	errorhandling.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// logout function for logging out from our own site
func (h *handler) restart(w http.ResponseWriter, r *http.Request) {
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	if session.Values["isLogin"] == false || model.Info["Username"] != "admin" {
		res := aprResponse{
			Status:  "error",
			Message: "access denied",
		}

		b, err := json.Marshal(res)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	out, err := h.AprService.restart()
	if err != nil {
		res := aprResponse{
			Status:  "error",
			Message: err.Error(),
		}

		b, err := json.Marshal(res)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	res := aprResponse{
		Status:  "success",
		Message: string(out),
	}
	b, err := json.Marshal(res)
	errorhandling.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
