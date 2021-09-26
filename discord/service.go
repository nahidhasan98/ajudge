package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
)

type discordInterfacer interface {
	SendMessage(data interface{}, notifier string) (*api.WebhookMessage, error)
	EditMessage(data model.SubmissionData) (*api.WebhookMessage, error)
	DeleteMessage(msgID api.Snowflake) error
}

type discordStruct struct {
	repoService repoInterfacer
}

func (ds discordStruct) SendMessage(data interface{}, notifier string) (*api.WebhookMessage, error) {
	var webhookID, webhookToken string
	var disMsg string
	var subID int

	// preparing message to send
	switch notifier {

	case "submission":
		temp := data.(model.SubmissionData)
		subID = temp.SubID

		timeDotTime := time.Unix(temp.SubmittedAt, 0)
		formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")

		disMsg = prepareSubmissionMessage(temp, formattedTime)

		webhookID = vault.WebhookIDSub
		webhookToken = vault.WebhookTokenSub

	case "login":
		temp := data.(model.UserData)

		timeDotTime := time.Unix(temp.CreatedAt, 0)
		formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")

		disMsg = prepareLoginRegMessage(temp, formattedTime)

		webhookID = vault.WebhookIDLogin
		webhookToken = vault.WebhookTokenLogin

	case "registration":
		temp := data.(model.UserData)

		timeDotTime := time.Unix(temp.CreatedAt, 0)
		formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")

		disMsg = prepareLoginRegMessage(temp, formattedTime)

		webhookID = vault.WebhookIDReg
		webhookToken = vault.WebhookTokenReg

	case "contest":
		temp := data.(model.ContestData)

		timeDotTime := time.Unix(temp.StartAt, 0)
		formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")

		disMsg = prepareContestMessage(temp, formattedTime)

		webhookID = vault.WebhookIDContest
		webhookToken = vault.WebhookTokenContest

	case "resetPass":
		temp := data.(model.UserData)

		disMsg = prepareResetMessage(temp)

		webhookID = vault.WebhookIDReset
		webhookToken = vault.WebhookTokenReset

	case "feedback":
		temp := data.(model.FeedbackData)

		disMsg = prepareFeedbackMessage(temp)

		webhookID = vault.WebhookIDFeedback
		webhookToken = vault.WebhookTokenFeedback
	}

	// innitializing webhook
	webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(webhookID), webhookToken)
	errorhandling.Check(err)

	// sending msg to discord
	res, err := webhook.SendContent(disMsg)
	errorhandling.Check(err)

	// store to DB
	err = ds.repoService.storeMsgID(subID, fmt.Sprintf("%v", res.ID), disMsg, notifier)
	errorhandling.Check(err)

	return res, err
}

func (ds discordStruct) EditMessage(data model.SubmissionData) (*api.WebhookMessage, error) {
	// getting sent msg details by subID
	sentData, err := ds.repoService.getDetails(data.SubID)
	errorhandling.Check(err)

	// preparing message to send
	disMsg := prepareSubmissionEditedMessage(sentData.Message, data)

	// innitializing webhook
	webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(vault.WebhookIDSub), vault.WebhookTokenSub)
	errorhandling.Check(err)

	// editing sent msg to discord
	res, err := webhook.EditContent(api.Snowflake(sentData.MessageID), disMsg)
	errorhandling.Check(err)

	// update sent msg info to DB
	err = ds.repoService.updateMsg(data.SubID, fmt.Sprintf("%v", res.ID), disMsg)
	errorhandling.Check(err)

	return res, err
}

func (ds discordStruct) DeleteMessage(msgID api.Snowflake) error {
	// innitializing webhook
	webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(vault.WebhookIDSub), vault.WebhookTokenSub)
	errorhandling.Check(err)

	// deleting sent msg to discord
	err = webhook.DeleteMessage(msgID)
	errorhandling.Check(err)

	return err
}

func prepareSubmissionMessage(data model.SubmissionData, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + data.OJ + " - " + data.Username + "\n"
	disMsg += "Submission ID" + getSpace(3) + ": " + fmt.Sprintf("%v", data.SubID) + "\n"
	disMsg += "Remote ID" + getSpace(7) + ": " + data.VID + "\n"
	disMsg += "Username" + getSpace(8) + ": " + data.Username + "\n"
	disMsg += "OJ" + getSpace(14) + ": " + data.OJ + "\n"
	disMsg += "Problem Number" + getSpace(2) + ": " + data.PNum + "\n"
	disMsg += "Problem Name" + getSpace(4) + ": " + data.PName + "\n"
	disMsg += "Language" + getSpace(8) + ": " + data.Language + "\n"
	disMsg += "Time" + getSpace(12) + ": " + data.TimeExec + "\n"
	disMsg += "Memory" + getSpace(10) + ": " + data.MemoryExec + "\n"
	disMsg += "Submitted At" + getSpace(4) + ": " + formattedTime + "\n"
	disMsg += "Verdict" + getSpace(9) + ": " + data.Verdict + "\n"
	disMsg += "Contest ID" + getSpace(6) + ": " + fmt.Sprintf("%v", data.ContestID) + "\n"
	disMsg += "Serial Index" + getSpace(4) + ": " + data.SerialIndex + "\n"
	disMsg += "```"

	return disMsg
}

func prepareLoginRegMessage(data model.UserData, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + data.Username + "\n"
	disMsg += "Username" + getSpace(4) + ": " + data.Username + "\n"
	disMsg += "Email" + getSpace(7) + ": " + data.Email + "\n"
	disMsg += "Full Name" + getSpace(3) + ": " + data.FullName + "\n"
	disMsg += "Verified" + getSpace(4) + ": " + fmt.Sprintf("%v", data.IsVerified) + "\n"
	disMsg += "Member Since" + getSpace(0) + ": " + formattedTime + "\n"
	disMsg += "Total Solved" + getSpace(0) + ": " + fmt.Sprintf("%v", data.TotalSolved) + "\n"
	disMsg += "```"

	return disMsg
}

func prepareContestMessage(data model.ContestData, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + fmt.Sprintf("%v", data.ContestID) + " - " + data.Title + "\n"
	disMsg += "Contest ID" + getSpace(2) + ": " + fmt.Sprintf("%v", data.ContestID) + "\n"
	disMsg += "Title" + getSpace(7) + ": " + data.Title + "\n"
	disMsg += "Start At" + getSpace(4) + ": " + formattedTime + "\n"
	disMsg += "Duration" + getSpace(4) + ": " + fmt.Sprintf("%v", data.Duration) + "\n"
	disMsg += "Author" + getSpace(6) + ": " + data.Author + "\n"
	disMsg += "```"

	return disMsg
}

func prepareResetMessage(data model.UserData) string {
	disMsg := "```md\n"
	disMsg += "Username" + getSpace(0) + ": " + data.Username + "\n"
	disMsg += "```"

	return disMsg
}

func prepareFeedbackMessage(data model.FeedbackData) string {
	disMsg := "```md\n"
	disMsg += "# " + data.Email + "\n"
	disMsg += "Name" + getSpace(4) + ": " + data.Name + "\n"
	disMsg += "Email" + getSpace(3) + ": " + data.Email + "\n"
	disMsg += "Message" + getSpace(1) + ": " + data.Message + "\n"
	disMsg += "```"

	return disMsg
}

func prepareSubmissionEditedMessage(old string, data model.SubmissionData) string {
	// tt := "Time1234567891234: ---\nMemory"
	idx1 := strings.Index(old, "Time")
	idx2 := strings.Index(old, "Memory")
	disMsgV1 := old[:idx1+18] + data.TimeExec + "\n" + old[idx2:]

	// tt := "Memory12345678912: ---\nSubmitted"
	idx1 = strings.Index(disMsgV1, "Memory")
	idx2 = strings.Index(disMsgV1, "Submitted")
	disMsgV2 := disMsgV1[:idx1+18] + data.MemoryExec + "\n" + disMsgV1[idx2:]

	// tt := "Verdict1234567891: ---\nContest"
	idx1 = strings.Index(disMsgV2, "Verdict")
	idx2 = strings.Index(disMsgV2, "Contest")
	disMsg := disMsgV2[:idx1+18] + data.Verdict + "\n" + disMsgV2[idx2:]

	return disMsg
}

// for indentation purpose
func getSpace(num int) string {
	str := ""

	for i := 0; i < num; i++ {
		str += " "
	}

	return str
}

func newDiscordService(repo repoInterfacer) discordInterfacer {
	return &discordStruct{
		repoService: repo,
	}
}
