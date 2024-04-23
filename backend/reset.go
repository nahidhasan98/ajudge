package backend

import (
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func resetCommon(w http.ResponseWriter, r *http.Request, title string) {
	session, _ := model.Store.Get(r, "mysession")
	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "Reset"
	model.Info["PageTitle"] = title + " | AJudge"
	model.Info["LastPage"] = model.LastPage

	model.Tpl.ExecuteTemplate(w, "reset.gohtml", model.Info)
}

// Reset function for requesting password or verification token
func Reset(w http.ResponseWriter, r *http.Request) { //calling from submit of 1.Pass reset or 2.Token reset
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	path := r.URL.Path

	if r.Method != "POST" {
		var title string
		if path == "/resetPassword" {
			title = "Reset Password"
			resetCommon(w, r, title)
		} else if path == "/resetToken" { //account verify token
			title = "New Token Request"
			resetCommon(w, r, title)
		} else {
			errorPage(w, http.StatusNotFound)
			return
		}
	} else if r.Method == "POST" {
		email := html.EscapeString(strings.TrimSpace(r.FormValue("email")))

		//connecting to DB
		DB, ctx, cancel := db.Connect()
		defer cancel()
		defer DB.Client().Disconnect(ctx)

		//taking DB collection/table to a variable
		userCollection := DB.Collection("user")

		//getting username for this email from DB for later use
		var userData model.UserData
		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&userData)

		if err == mongo.ErrNoDocuments { //if error occur in finding email
			errorPage(w, http.StatusInternalServerError)
			return
		}
		newToken := model.GenerateToken()

		if path == "/resetPassword" { //Request for Password reset
			//updating (or inserting) new token in the DB
			updateField := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "passResetToken", Value: newToken},
					{Key: "passResetTokenSentAt", Value: time.Now().Unix()},
				}},
			}
			_, err := userCollection.UpdateOne(ctx, bson.M{"email": email}, updateField)
			errorhandling.Check(err)

			//sending mail to the user with a reset link
			link := "https://ajudge.net/passReset/token=" + newToken //sending link for mail
			model.SendMail(email, userData.Username, link, "passwordReset")

			model.PopUpCause = "passwordRequest"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if path == "/resetToken" { //Request for Token reset (account verify token)
			//updating new token in the DB
			updateField := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "accVerifyToken", Value: newToken},
					{Key: "accVerifyTokenSentAt", Value: time.Now().Unix()},
				}},
			}
			_, err := userCollection.UpdateOne(ctx, bson.M{"email": email}, updateField)
			errorhandling.Check(err)

			//sending mail to the user with a reset link
			link := "https://ajudge.net/verify-email/token=" + newToken
			model.SendMail(email, userData.Username, link, "accVerify")

			model.PopUpCause = "tokenRequest"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			//this segment is not necessary because it is in the POST method. Above two section will cacth the flow.
			errorPage(w, http.StatusBadRequest)
			return
		}
	}
}

// PassReset function for resetting a user's password
func PassReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := model.Store.Get(r, "mysession")

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	userCollection := DB.Collection("user")

	if r.Method != "POST" { //flow comes here from email link
		path := r.URL.Path
		token := strings.TrimPrefix(path, "/passReset/token=")

		//getting data for this user from DB
		var userData model.UserData
		res := userCollection.FindOne(ctx, bson.M{"passResetToken": token}).Decode(&userData)

		if res == mongo.ErrNoDocuments { //Row/document not found
			model.PopUpCause = "passTokenInvalid"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		//found a row/document

		//checking for token expired or not
		tokenReceived := time.Now().Unix() //current time
		timeDiff := model.Abs(tokenReceived - userData.PassResetTokenSentAt)

		if timeDiff > (30 * 60) { //30 minutes period
			model.PopUpCause = "passTokenExpired"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		model.Info["Username"] = session.Values["username"]
		model.Info["IsLogged"] = session.Values["isLogin"]
		model.Info["PageName"] = "Reset Password"
		model.Info["PageTitle"] = "Reset Password | AJudge"
		model.Info["LastPage"] = model.LastPage
		model.Info["Token"] = token

		model.Tpl.ExecuteTemplate(w, "passReset.gohtml", model.Info)
	} else if r.Method == "POST" {
		//gettting form data
		token := r.FormValue("token") //hidden - send by us through email earlier
		password := html.EscapeString(r.FormValue("password"))
		password = model.HashPassword(password) //hashing password

		// validating form data
		if len(password) < 8 {
			// msg: "password length should be at least 8 characters"

			model.PopUpCause = "passwordResetErr"
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		//updating new password in the DB
		updateField := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "password", Value: password},
			}},
		}
		_, err := userCollection.UpdateOne(ctx, bson.M{"passResetToken": token}, updateField)
		errorhandling.Check(err)

		model.PopUpCause = "passwordReset"
		http.Redirect(w, r, "/", http.StatusSeeOther)

		// notofy to discord
		disData := model.UserData{
			Username: session.Values["username"].(string),
		}
		discord := discord.Init()
		discord.SendMessage(disData, "resetPass")

		return
	}
}
