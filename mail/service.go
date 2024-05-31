package mail

import (
	"bytes"
	"fmt"
	"text/template"
)

type mailInterfacer interface {
	SendMailForRegistration(email, username, link string)
	SendMailForPasswordReset(email, username, link string)
}

type mailStruct struct {
	Tpl *template.Template
}

func (m mailStruct) SendMailForRegistration(email, username, link string) {
	mailBody := prepareMailBody("verification", m.Tpl, email, username, link)
	go sendWorker(email, mailBody)
}

func (m mailStruct) SendMailForPasswordReset(email, username, link string) {
	mailBody := prepareMailBody("reset", m.Tpl, email, username, link)
	go sendWorker(email, mailBody)
}

func prepareMailBody(what string, tpl *template.Template, email, username, link string) []byte {
	var subject, description1, buttonText, description2 string

	switch what {
	case "verification":
		subject = "Ajudge Account Verification"
		description1 = `Thank you for signed up. To start journey with us, please click the link below to verify your email address.`
		buttonText = "Verify Account"
		description2 = `This account verification link is valid for next 30 minutes. After link expiration, the above link will help you to get a new link.`
	case "reset":
		subject = "Ajudge Password Reset Link"
		description1 = `Someone has requested a new password for your Ajudge account. Please click the link below to reset your password.`
		buttonText = "Reset Password"
		description2 = `If you did not request a password reset, please ignore this email or reply to let us know. This password reset link is valid for next 30 minutes.`
	}

	// preparing body
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("From: Ajudge Team \nSubject: %s \nTo:%s \n%s\n\n", subject, email, mimeHeaders)))

	tData := struct {
		Username, Link, Description1, Description2, ButtonText string
	}{
		Username:     username,
		Link:         link,
		Description1: description1,
		Description2: description2,
		ButtonText:   buttonText,
	}

	tpl.ExecuteTemplate(&body, "mail.gohtml", tData)

	return body.Bytes()
}

func NewMailService() mailInterfacer {
	return &mailStruct{
		Tpl: template.Must(template.ParseGlob("frontend/html/mail.gohtml")),
	}
}
