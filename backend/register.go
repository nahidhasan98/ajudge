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
