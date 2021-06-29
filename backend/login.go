package backend

import (
	"html"
	"net/http"
	"strings"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//Login function for login to our own site
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := model.Store.Get(r, "mysession")

	if r.Method != "POST" {
		if session.Values["isLogin"] == true {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			//preparing info for sending to frontend
			model.Info["PageName"] = "Login"
			model.Info["PageTitle"] = "Login | AJudge"
			model.Info["LastPage"] = model.LastPage
			model.Info["PopUpCause"] = model.PopUpCause

			model.Tpl.ExecuteTemplate(w, "login.gohtml", model.Info)

			//clearing up variables
			model.Info["Username"] = ""
			model.Info["ErrUsername"] = ""
			model.Info["ErrPassword"] = ""
			model.PopUpCause = ""
			model.Info["PopUpCause"] = model.PopUpCause
		}
	} else if r.Method == "POST" {
		//getting form data
		username := html.EscapeString(strings.TrimSpace(r.FormValue("username")))
		password := html.EscapeString(r.FormValue("password"))

		//clearing up html display error (wheather error exist or not will be checked up below/later)
		model.Info["ErrUsername"] = ""
		model.Info["ErrPassword"] = ""

		//connecting to DB
		DB, ctx, cancel := db.Connect()
		defer cancel()
		defer DB.Client().Disconnect(ctx)

		//taking DB collection/table to a variable
		userCollection := DB.Collection("user")

		//getting original password for this user from DB
		var userData model.UserData
		err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userData)

		if err == mongo.ErrNoDocuments { //if username not exist/found in DB (returned no document/row)
			model.Info["Username"] = username
			model.Info["ErrUsername"] = "Username not found."

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if checkPasswordHash(password, userData.Password) { //if password matched
			session.Values["username"] = username
			session.Values["password"] = password
			session.Values["isLogin"] = true
			session.Save(r, w)

			//preparing info for sending to frontend
			model.Info["Username"] = session.Values["username"]
			model.Info["Password"] = session.Values["password"]
			model.Info["IsLogged"] = session.Values["isLogin"]

			http.Redirect(w, r, model.LastPage, http.StatusSeeOther)
			return
		}
		//if password not matched
		model.Info["Username"] = username
		model.Info["ErrPassword"] = "Invalid password"

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

//Logout function for logging out from our own site
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := model.Store.Get(r, "mysession")
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Values["isLogin"] = false

	session.Options.MaxAge = -1 //cookies will be deleted immediately.
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//function used above by this particular file
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
