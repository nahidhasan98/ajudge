package discord

import (
	"fmt"
	"time"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
	discordtexthook "github.com/nahidhasan98/discord-text-hook"
)

func sendWorker(jobs <-chan int, data interface{}, notifier string, ds discordStruct) {
	// fmt.Println("Worker pool is working..")

	for range jobs {
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
			temp := data.(UserModel)

			timeDotTime := time.Unix(temp.CreatedAt, 0)
			formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")

			disMsg = prepareLoginRegMessage(temp, formattedTime)

			webhookID = vault.WebhookIDLogin
			webhookToken = vault.WebhookTokenLogin

		case "registration":
			temp := data.(UserModel)

			timeDotTime := time.Unix(temp.CreatedAt, 0)
			formattedTime := timeDotTime.Format("02-Jan-2006 (15:04:05)")

			disMsg = prepareLoginRegMessage(temp, formattedTime)

			webhookID = vault.WebhookIDReg
			webhookToken = vault.WebhookTokenReg

		case "contest":
			temp := data.(model.ContestData)
			subID = temp.ContestID // we need contest ID for retrieve messageID (for editing message). using subID as contestID without adding a new field to the model

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
		webhook := discordtexthook.NewDiscordTextHookService(webhookID, webhookToken)

		// sending msg to discord
		res, err := webhook.SendMessage(disMsg)
		errorhandling.Check(err)

		// store to DB
		err = ds.repoService.storeMsgID(subID, fmt.Sprintf("%v", res.ID), disMsg, notifier)
		errorhandling.Check(err)
	}
}

func editWorker(jobs <-chan int, data interface{}, notifier string, ds discordStruct) {
	// fmt.Println("Worker pool is working..")

	for range jobs {
		var webhookID, webhookToken string
		var disMsg string
		var subID int
		var sentData discordModel
		var err error

		// preparing message to send
		switch notifier {

		case "submission":
			temp := data.(model.SubmissionData)
			subID = temp.SubID

			// getting sent msg details by subID
			sentData, err = ds.repoService.getDetails(subID, notifier)
			errorhandling.Check(err)
			if err != nil {
				return
			}

			// preparing message to send
			disMsg = prepareSubmissionEditedMessage(sentData.Message, temp)

			webhookID = vault.WebhookIDSub
			webhookToken = vault.WebhookTokenSub

		case "contest":
			temp := data.(model.ContestData)
			subID = temp.ContestID // we need contest ID for retrieve messageID (for editing message). using subID as contestID without adding a new field to the model

			// getting sent msg details by subID
			sentData, err = ds.repoService.getDetails(subID, notifier)
			errorhandling.Check(err)
			if err != nil {
				return
			}

			// preparing message to send
			disMsg = prepareContestEditedMessage(sentData.Message, temp)

			webhookID = vault.WebhookIDContest
			webhookToken = vault.WebhookTokenContest
		}

		// innitializing webhook
		webhook := discordtexthook.NewDiscordTextHookService(webhookID, webhookToken)
		errorhandling.Check(err)

		// editing sent msg to discord
		res, err := webhook.EditMessage(disMsg, sentData.MessageID)
		errorhandling.Check(err)

		// update sent msg info to DB
		err = ds.repoService.updateMsg(subID, fmt.Sprintf("%v", res.ID), disMsg)
		errorhandling.Check(err)
	}
}
