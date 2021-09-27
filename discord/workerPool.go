package discord

import (
	"fmt"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/model"
	"github.com/nahidhasan98/ajudge/vault"
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
	}
}

func editWorker(jobs <-chan int, data model.SubmissionData, ds discordStruct) {
	// fmt.Println("Worker pool is working..")

	for range jobs {
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
	}
}

func deleteWorker(jobs <-chan int, msgID api.Snowflake) {
	// fmt.Println("Worker pool is working..")

	for range jobs {
		// innitializing webhook
		webhook, err := disgohook.NewWebhookClientByIDToken(nil, nil, api.Snowflake(vault.WebhookIDSub), vault.WebhookTokenSub)
		errorhandling.Check(err)

		// deleting sent msg to discord
		err = webhook.DeleteMessage(msgID)
		errorhandling.Check(err)
	}
}
