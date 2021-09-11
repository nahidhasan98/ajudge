package model

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/vault"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//IsAccVerifed func
func IsAccVerifed(r *http.Request) bool {
	session, _ := Store.Get(r, "mysession")
	username := session.Values["username"]

	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking  DB collection/table to a variable
	userCollection := DB.Collection("user")

	//retrieving data from DB
	var userData UserData
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userData)
	errorhandling.Check(err)

	return userData.IsVerified
}

//Min function
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Abs function
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
func removeStyle(styleBody string) string {
	need1 := "<style"
	index1 := strings.Index(styleBody, need1)
	need2 := "</style>"
	index2 := strings.Index(styleBody, need2)

	var part1, part2 string
	if index1 != -1 {
		part1 = styleBody[0:index1]
		part2 = styleBody[index2+8:]

		styleBody = part1 + part2
	}
	return styleBody
}

//HashPassword function
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

//SendMail function
func SendMail(email, username, link, resetType string) {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", vault.GmailUsername, vault.GmailPassword, "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}

	var subject, description1, buttonText, description2 string

	//var msg []byte
	var body bytes.Buffer
	if resetType == "passwordReset" {
		subject = "Ajudge Password Reset Link"
		description1 = `Someone has requested a new password for your Ajudge account. Please click the link below to reset your password.`
		buttonText = "Reset Password"
		description2 = `If you did not request a password reset, please ignore this email or reply to let us know. This password reset link is valid for next 30 minutes.`
	} else if resetType == "accVerify" {
		subject = "Ajudge Account Verification"
		description1 = `Thank you for signed up. To start journey with us, please click the link below to verify your email address.`
		buttonText = "Verify Account"
		description2 = `This account verification link is valid for next 30 minutes. After link expiration, the above link will help you to get a new link.`
	}

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("From: Ajudge Team \nSubject: %s \nTo:%s \n%s\n\n", subject, email, mimeHeaders)))

	Tpl.ExecuteTemplate(&body, "sendMail.gohtml", struct {
		Username, Link, Description1, Description2, ButtonText string
	}{
		Username:     username,
		Link:         link,
		Description1: description1,
		Description2: description2,
		ButtonText:   buttonText,
	})

	err := smtp.SendMail("smtp.gmail.com:587", auth, "", to, body.Bytes())
	errorhandling.Check(err)
}

//GenerateToken function
func GenerateToken() string {
	b := make([]byte, 16)
	rand.Read(b)

	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

//IsExistInTV function
func IsExistInTV(OJ string, arr []string, verdict string) bool {
	var extra1, extra2, extra3, extra4, extra5, extra6, extra7 string

	if OJ == "CodeForces" || OJ == "Gym" || OJ == "SGU" {
		if len(verdict) > 12 { //for Wrong answer CodeForces gives like: "Wrong answer on test 5"
			extra1 = verdict[:12] //so we are taking only "Wrong answer" for checking existance in terminal verdict
		}
		if len(verdict) > 19 { //for Time limit exceeded CodeForces gives like: "Time limit exceeded on test 27"
			extra2 = verdict[:19] //so we are taking only "Time limit exceeded" for checking existance in terminal verdict
		}
		if len(verdict) > 21 { //for Memory limit exceeded CodeForces gives like: "Memory limit exceeded on test 2"
			extra3 = verdict[:21] //so we are taking only "Memory limit exceeded" for checking existance in terminal verdict
		}
		if len(verdict) > 13 { //for Runtime error CodeForces gives like: "Runtime error on test 1"
			extra4 = verdict[:13] //so we are taking only "Runtime error" for checking existance in terminal verdict
		}
	} else if OJ == "HDU" {
		if len(verdict) > 13 { //"Runtime Error(INTEGER_DIVIDE_BY_ZERO)" --> "Runtime Error"
			extra1 = verdict[:13]
		}
	} else if OJ == "SPOJ" {
		if len(verdict) > 12 { //"Wrong answer #0" --> "Wrong answer"
			extra1 = verdict[:12]
		}
		if len(verdict) > 13 { //"Runtime error (SIGSEGV)" --> "Runtime error"
			extra2 = verdict[:13]
		}
	} else if OJ == "UESTC" {
		if len(verdict) > 18 { //"Presentation Error on test 1" --> "Presentation error"
			extra1 = verdict[:18]
		}
		if len(verdict) > 12 { //"Wrong Answer on test 1" --> "Wrong answer"
			extra2 = verdict[:12]
		}
		if len(verdict) > 12 { //"System Error on test 1" --> "System Error"
			extra3 = verdict[:12]
		}
		if len(verdict) > 19 { //"Time Limit Exceeded on test 27" --> "Time Limit Exceeded"
			extra4 = verdict[:19]
		}
		if len(verdict) > 21 { //"Memory Limit Exceeded on test 24" --> "Memory Limit Exceeded"
			extra5 = verdict[:21]
		}
		if len(verdict) > 21 { //"Output Limit Exceeded on test 1" --> "Output Limit Exceeded"
			extra6 = verdict[:21]
		}
		if len(verdict) > 19 { //"Restricted Function on test 2" --> "Restricted Function"
			extra7 = verdict[:19]
		}
	} else if OJ == "UniversalOJ" {
		if len(verdict) > 17 { //"Extra Test Failed : Wrong Answer on 7" --> "Extra Test Failed"
			extra1 = verdict[:17]
		}
	} else if OJ == "URAL" {
		if len(verdict) > 13 { //"Runtime error (non-zero exit code)" --> "Runtime error"
			extra1 = verdict[:13]
		}
	} else if OJ == "URI" {
		if len(verdict) > 18 { //"Presentation error (10%)" --> "Presentation error"
			extra1 = verdict[:18]
		}
		if len(verdict) > 12 { //"Wrong answer (10%)" --> "Wrong answer"
			extra2 = verdict[:12]
		}
	}

	//checking in array
	for k := 0; k < len(arr); k++ {
		if arr[k] == verdict || arr[k] == extra1 || arr[k] == extra2 || arr[k] == extra3 || arr[k] == extra4 || arr[k] == extra5 || arr[k] == extra6 || arr[k] == extra7 {
			return true
		}
	}
	return false
}
