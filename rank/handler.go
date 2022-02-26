package rank

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
)

type handler struct {
	rankService rankInterfacer
}

// makeHTTPHandlers function defines endpoints
func 	makeHTTPHandlers(router *mux.Router, rankService rankInterfacer) {
	h := &handler{
		rankService: rankService,
	}

	router.HandleFunc("/rank/{type}", h.previewRankingPage).Methods("GET")
	router.HandleFunc("/rank/list/{type}", h.rankList).Methods("GET")
}

func (h *handler) previewRankingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, err := model.Store.Get(r, "mysession")
	errorhandling.Check(err)

	model.LastPage = r.URL.Path

	id := mux.Vars(r)["type"]
	if id != "oj" && id != "user" {
		w.WriteHeader(http.StatusNotFound) // status code such as: 400, 404 etc.
		model.Info["StatusCode"] = http.StatusNotFound
		model.Tpl.ExecuteTemplate(w, "pageNotFound.gohtml", model.Info)
		return
	}

	var rankType string
	if id == "oj" {
		rankType = "OJ"
	} else if id == "user" {
		rankType = "User"
	}

	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "Rank"
	model.Info["PageTitle"] = rankType + " Rank | AJudge"
	model.Info["RankType"] = rankType

	model.Tpl.ExecuteTemplate(w, "rank.gohtml", model.Info)
}

type rankResponse struct {
	Status   string      `json:"status,omitempty"`
	Message  string      `json:"message,omitempty"`
	RankList []rankModel `json:"rankList,omitempty"`
}

func (h *handler) rankList(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["type"]
	if id != "oj" && id != "user" {
		response := rankResponse{
			Status:  "error",
			Message: "invalid url",
		}
		byteData, err := json.Marshal(response)
		errorhandling.Check(err)

		w.WriteHeader(http.StatusNotFound) // status code such as: 400, 404 etc.
		w.Header().Set("Content-Type", "application/json")
		w.Write(byteData)
		return
	}

	rankList, err := h.rankService.getRankList(id)
	if err != nil {
		response := rankResponse{
			Status:  "error",
			Message: err.Error(),
		}
		byteData, err := json.Marshal(response)
		errorhandling.Check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(byteData)
		return
	}

	// preparing data to response
	response := rankResponse{
		Status:   "success",
		RankList: rankList,
	}

	b, err := json.Marshal(response)
	errorhandling.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
