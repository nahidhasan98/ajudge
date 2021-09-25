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
	SendMessage(data model.SubmissionData) (*api.WebhookMessage, error)
	EditMessage(data model.SubmissionData) (*api.WebhookMessage, error)
	DeleteMessage(msgID api.Snowflake) error
}

type discordStruct struct {
	repoService repoInterfacer
}

func (ds discordStruct) SendMessage(data model.SubmissionData) (*api.WebhookMessage, error) {
	// preparing message to send
	timeDotTime := time.Unix(data.SubmittedAt, 0)
	formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")
	disMsg := prepareMessage(data, formattedTime)

	// innitializing webhook
	webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(vault.WebhookID), vault.WebhookToken)
	errorhandling.Check(err)

	// sending msg to discord
	res, err := webhook.SendContent(disMsg)
	errorhandling.Check(err)

	// store to DB
	err = ds.repoService.storeMsgID(data.SubID, fmt.Sprintf("%v", res.ID), disMsg)
	errorhandling.Check(err)

	return res, err
}

func (ds discordStruct) EditMessage(data model.SubmissionData) (*api.WebhookMessage, error) {
	// getting sent msg details by subID
	sentData, err := ds.repoService.getDetails(data.SubID)
	errorhandling.Check(err)

	// preparing message to send
	disMsg := prepareEditedMessage(sentData.Message, data)

	// innitializing webhook
	webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(vault.WebhookID), vault.WebhookToken)
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
	webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(vault.WebhookID), vault.WebhookToken)
	errorhandling.Check(err)

	// deleting sent msg to discord
	err = webhook.DeleteMessage(msgID)
	errorhandling.Check(err)

	return err
}

func prepareMessage(data model.SubmissionData, formattedTime string) string {
	disMsg := "```md\n"
	disMsg += "# " + data.OJ + " - " + data.Username + "\n"
	disMsg += "Submission ID" + getSpace(4) + ": " + fmt.Sprintf("%v", data.SubID) + "\n"
	disMsg += "Remote ID" + getSpace(8) + ": " + data.VID + "\n"
	disMsg += "Username" + getSpace(9) + ": " + data.Username + "\n"
	disMsg += "OJ" + getSpace(15) + ": " + data.OJ + "\n"
	disMsg += "Problem Number" + getSpace(3) + ": " + data.PNum + "\n"
	disMsg += "Problem Name" + getSpace(5) + ": " + data.PName + "\n"
	disMsg += "Language" + getSpace(9) + ": " + data.Language + "\n"
	disMsg += "Time" + getSpace(13) + ": " + data.TimeExec + "\n"
	disMsg += "Memory" + getSpace(11) + ": " + data.MemoryExec + "\n"
	disMsg += "Submitted At" + getSpace(5) + ": " + formattedTime + "\n"
	disMsg += "Verdict" + getSpace(10) + ": " + data.Verdict + "\n"
	disMsg += "Contest ID" + getSpace(7) + ": " + fmt.Sprintf("%v", data.ContestID) + "\n"
	disMsg += "SerialIndex" + getSpace(6) + ": " + data.SerialIndex + "\n"
	disMsg += "```"

	return disMsg
}

func prepareEditedMessage(old string, data model.SubmissionData) string {
	// tt := "Time1234567891234: ---\nMemory"
	idx1 := strings.Index(old, "Time")
	idx2 := strings.Index(old, "Memory")
	disMsgV1 := old[:idx1+19] + data.TimeExec + "\n" + old[idx2:]

	// tt := "Memory12345678912: ---\nSubmitted"
	idx1 = strings.Index(disMsgV1, "Memory")
	idx2 = strings.Index(disMsgV1, "Submitted")
	disMsgV2 := disMsgV1[:idx1+19] + data.MemoryExec + "\n" + disMsgV1[idx2:]

	// tt := "Verdict1234567891: ---\nContest"
	idx1 = strings.Index(disMsgV2, "Verdict")
	idx2 = strings.Index(disMsgV2, "Contest")
	disMsg := disMsgV2[:idx1+19] + data.Verdict + "\n" + disMsgV2[idx2:]

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
