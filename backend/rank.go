package backend

import (
	"net/http"

	"github.com/nahidhasan98/ajudge/model"
)

//Rank function for ranking OJ & user based on number of problem solved
func Rank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	path := r.URL.Path

	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	if path != "/rankOJ" && path != "/rankUser" {
		errorPage(w, http.StatusBadRequest)
		return
	}
	rankType := path[5:]

	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "Rank"
	model.Info["PageTitle"] = rankType + " Rank | AJudge"
	model.Info["RankType"] = rankType

	model.Tpl.ExecuteTemplate(w, "rank.gohtml", model.Info)
}
