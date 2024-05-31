package mail

import (
	"fmt"
	"net/smtp"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/vault"
)

func sendWorker(email string, mailBody []byte) {
	// one or more emails address to which mail will be sent
	to := []string{email}

	// Choose auth method and set it up
	auth := smtp.PlainAuth("", vault.GmailUsername, vault.GmailPassword, vault.SmtpServiceHost)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", vault.SmtpServiceHost, vault.SmtpServicePort), auth, "", to, mailBody)
	errorhandling.Check(err)
}
