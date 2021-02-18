package backend

import (
	"net/http"
	"strings"
	"time"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Profile function for a user's statistics
func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	model.LastPage = "profile/"
	session, _ := model.Store.Get(r, "mysession")

	path := r.URL.Path
	user := strings.TrimPrefix(path, "/profile/")

	if (user == "" || user == "/profile") && session.Values["isLogin"] == true {
		http.Redirect(w, r, "/profile/"+session.Values["username"].(string), http.StatusSeeOther)
		return
	}

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	userCollection := DB.Collection("user")

	//getting data for this user from DB
	var userData model.UserData
	err := userCollection.FindOne(ctx, bson.M{"username": user}).Decode(&userData)

	if err == mongo.ErrNoDocuments { //wrong username (no records in DB)
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
		return
	}

	self := false //for displaying source code (if visit self profile)
	if user == session.Values["username"] {
		self = true
	}

	//formating createdAt time to display
	timeDotTime := time.Unix(userData.CreatedAt, 0)
	createdAt := timeDotTime.Format("02-Jan-2006")

	//user submissioon list will be gathered by frontend ajax call

	//preparing info for sending to frontend
	model.Info["Username"] = session.Values["username"]
	model.Info["Password"] = session.Values["password"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "Profile"
	model.Info["PageTitle"] = user + `'s Profile | AJudge`
	model.Info["Self"] = self
	model.Info["FullName"] = userData.FullName
	model.Info["Email"] = userData.Email
	model.Info["CreatedAt"] = createdAt
	model.Info["PopUpCause"] = model.PopUpCause

	model.Tpl.ExecuteTemplate(w, "profile.gohtml", model.Info)

	//clearing up values (because it may be used in wrong place unintentionally)
	model.Info["FullName"] = ""
	model.Info["Email"] = ""
	model.PopUpCause = ""
	model.Info["PopUpCause"] = model.PopUpCause
}
