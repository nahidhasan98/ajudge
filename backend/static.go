package backend

import (
	"html"
	"net/http"
	"os"
	"strings"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/discord"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Index function for Homepage
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	//removing temporary images (downloaded earlier for problemView pages)
	os.RemoveAll("assets/temp/")
	//fmt.Println("File deleted")

	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "Homepage"
	model.Info["PageTitle"] = "AJudge | All your favourite OJ's are in one platform"
	model.Info["LastPage"] = model.LastPage
	model.Info["PopUpCause"] = model.PopUpCause

	model.Tpl.ExecuteTemplate(w, "index.gohtml", model.Info)

	//clearing up values (because it may be used in wrong place unintentionally)
	model.PopUpCause = ""
	model.Info["PopUpCause"] = model.PopUpCause
}

// About funtcion for about page
func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	model.LastPage = r.URL.Path
	session, _ := model.Store.Get(r, "mysession")

	model.Info["Username"] = session.Values["username"]
	model.Info["IsLogged"] = session.Values["isLogin"]
	model.Info["PageName"] = "About"
	model.Info["PageTitle"] = "About | AJudge"

	model.Tpl.ExecuteTemplate(w, "about.gohtml", model.Info)
}

// Contact funtion for contact page
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		model.LastPage = r.URL.Path
		session, _ := model.Store.Get(r, "mysession")

		model.Info["Username"] = session.Values["username"]
		model.Info["IsLogged"] = session.Values["isLogin"]
		model.Info["PageName"] = "Contact"
		model.Info["PageTitle"] = "Contact | AJudge"
		model.Info["LastPage"] = model.LastPage
		model.Info["PopUpCause"] = model.PopUpCause

		model.Tpl.ExecuteTemplate(w, "contact.gohtml", model.Info)

		//clearing up values (because it may be used in wrong place unintentionally)
		model.PopUpCause = ""
		model.Info["PopUpCause"] = model.PopUpCause
	} else if r.Method == "POST" {
		//getting form data
		name := html.EscapeString(strings.TrimSpace(r.FormValue("mailName")))
		email := html.EscapeString(strings.TrimSpace(r.FormValue("mailEmail")))
		message := html.EscapeString(strings.TrimSpace(r.FormValue("mailMessage")))

		// connecting to DB
		DB, ctx, cancel := db.Connect()
		defer cancel()
		defer DB.Client().Disconnect(ctx)

		// taking DB collection/table to a variable
		userCollection := DB.Collection("user")

		// retrieving data from DB
		var userData model.UserData
		err := userCollection.FindOne(ctx, bson.M{"email": email, "isVerified": true}).Decode(&userData)

		// if email not verified then reject
		if err == mongo.ErrNoDocuments {
			model.PopUpCause = "userFeedbackReject"
			http.Redirect(w, r, "/contact", http.StatusSeeOther)
			return
		}

		// pause mail temporarily for stop spam

		// //sending mail to our email
		// auth := smtp.PlainAuth("", vault.GmailUsername, vault.GmailPassword, vault.SmtpServiceHost)
		// to := []string{"mnh.nahid35@gmail.com"}

		// var msg = []byte("From: " + name + "\r\n" +
		// 	"To: mnh.nahid35@gmail.com \r\n" +
		// 	"Subject: Ajudge User Feedback\r\n" +
		// 	"Cc: mugdo1997@gmail.com \r\n" +
		// 	"\r\n" +
		// 	"Sender's Name: " + name + "\r\n" +
		// 	"Sender's Email: " + email + "\r\n" +
		// 	"Message: " + message)

		// err := smtp.SendMail(fmt.Sprintf("%s:%s", vault.SmtpServiceHost, vault.SmtpServicePort), auth, "", to, msg)
		// errorhandling.Check(err)

		model.PopUpCause = "userFeedback"
		http.Redirect(w, r, "/contact", http.StatusSeeOther)

		// notofy to discord
		disData := model.FeedbackData{
			Name:    name,
			Email:   email,
			Message: message,
		}
		discord := discord.Init()
		discord.SendMessage(disData, "feedback")

		return
	}
}

// PageNotFound function for bad URL handling
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	errorPage(w, http.StatusNotFound)
}
func errorPage(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode) //status code such as: 400, 404 etc.

	model.Info["StatusCode"] = statusCode

	model.Tpl.ExecuteTemplate(w, "pageNotFound.gohtml", model.Info)
}
