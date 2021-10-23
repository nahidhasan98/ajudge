package backend

import (
	"net/http"
	"strings"
	"time"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Register function for registering to our own site
// func Register(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

// 	if r.Method != "POST" {
// 		session, _ := model.Store.Get(r, "mysession")
// 		if session.Values["isLogin"] == true {
// 			http.Redirect(w, r, "/", http.StatusSeeOther)
// 		} else {
// 			model.Info["PageName"] = "Registration"
// 			model.Info["PageTitle"] = "Registration | AJudge"

// 			model.Tpl.ExecuteTemplate(w, "register.gohtml", model.Info)
// 		}
// 	} else if r.Method == "POST" {
// 		//gettting form data
// 		fullName := html.EscapeString(strings.TrimSpace(r.FormValue("fullName")))
// 		email := html.EscapeString(strings.TrimSpace(r.FormValue("email")))
// 		username := html.EscapeString(strings.TrimSpace(r.FormValue("username")))
// 		password := html.EscapeString(r.FormValue("password"))
// 		password = model.HashPassword(password) //hashing password
// 		newToken := model.GenerateToken()       //generating token for account verification

// 		//do regitration
// 		//connecting to DB
// 		DB, ctx, cancel := db.Connect()
// 		defer cancel()
// 		defer DB.Client().Disconnect(ctx)

// 		//taking DB collection/table to a variable
// 		userCollection := DB.Collection("user")
// 		counterCollection := DB.Collection("counter")

// 		//getting LastUserID from DB to assign a user ID(LastUserID+1) for this user
// 		var lastUserID model.LastUsedID
// 		err := counterCollection.FindOne(ctx, bson.M{}).Decode(&lastUserID)
// 		errorhandling.Check(err)

// 		//preparing data for inserting to DB
// 		userData := model.UserData{
// 			UserID:               lastUserID.LastUserID + 1,
// 			FullName:             fullName,
// 			Email:                email,
// 			Username:             username,
// 			Password:             password,
// 			CreatedAt:            time.Now().Unix(),
// 			IsVerified:           false,
// 			AccVerifyToken:       newToken,
// 			AccVerifyTokenSentAt: time.Now().Unix(),
// 			PassResetToken:       "",
// 			PassResetTokenSentAt: 0,
// 		}

// 		_, err = userCollection.InsertOne(ctx, userData)
// 		errorhandling.Check(err)

// 		//updating LastUserID to DB for later use/next user
// 		updateField := bson.D{
// 			{Key: "$inc", Value: bson.D{ //incrementing LastUserID by 1
// 				{Key: "lastUserID", Value: 1},
// 			}},
// 		}
// 		_, err = counterCollection.UpdateOne(ctx, bson.M{}, updateField)
// 		errorhandling.Check(err)

// 		//sending mail to user email with a verification link
// 		linkforMail := "https://ajudge.net/verify-email/token=" + newToken
// 		model.SendMail(email, username, linkforMail, "accVerify")

// 		model.PopUpCause = "registrationDone" //login page will give a popup
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)

// 		// notofy to discord
// 		disData := userData
// 		discord := discord.Init()
// 		discord.SendMessage(disData, "registration")

// 		return
// 	}
// }

//EmailVerifiation function for verify registerd email
func EmailVerifiation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	token := strings.TrimPrefix(path, "/verify-email/token=")

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	userCollection := DB.Collection("user")

	//retrieving data from DB
	var userData model.UserData
	err := userCollection.FindOne(ctx, bson.M{"accVerifyToken": token}).Decode(&userData)

	//checking weather the account is already verified or not
	if userData.IsVerified { //already verified
		model.PopUpCause = "tokenAlreadyVerified"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//if account already not verified then go for next move
	if err == mongo.ErrNoDocuments { //no row found (token not found) (returned no document/row)
		model.PopUpCause = "tokenInvalid"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//at this point account not already verified & token is valid
	//checking for token expired or not
	tokenReceivedAt := time.Now().Unix()
	timeDiff := tokenReceivedAt - userData.AccVerifyTokenSentAt

	if timeDiff > (30 * 60) { //30 minutes period (converting to seconds)
		model.PopUpCause = "tokenExpired"
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//token not expired. SO this is a valid request
	//updating account verify status for this user
	updateField := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "isVerified", Value: true},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"accVerifyToken": token}, updateField)
	errorhandling.Check(err)

	model.PopUpCause = "tokenVerifiedNow"
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
