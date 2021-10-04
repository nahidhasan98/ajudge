package discord

import (
	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"go.mongodb.org/mongo-driver/bson"
)

type repoInterfacer interface {
	storeMsgID(subID int, msgID, msg, notifier string) error
	getDetails(subID int, notifier string) (discordModel, error)
	updateMsg(subID int, msgID, msg string) error
}

type repoStruct struct {
	DBTable string
}

func (r *repoStruct) storeMsgID(subID int, msgID, msg, notifier string) error {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// selecting colllection
	coll := DB.Collection(r.DBTable)

	disData := discordModel{
		SubID:            subID,
		MessageID:        msgID,
		Message:          msg,
		NotificationType: notifier,
	}

	// inserting info to DB
	_, err := coll.InsertOne(ctx, disData)
	errorhandling.Check(err)

	return err
}

func (r *repoStruct) getDetails(subID int, notifier string) (discordModel, error) {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// selecting colllection
	coll := DB.Collection(r.DBTable)

	//getting msgID from DB
	var temp discordModel
	err := coll.FindOne(ctx, bson.M{"subID": subID, "notificationType": notifier}).Decode(&temp)
	errorhandling.Check(err)

	return temp, err
}

func (r *repoStruct) updateMsg(subID int, msgID, msg string) error {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// selecting colllection
	coll := DB.Collection(r.DBTable)

	// updating info to DB
	updateField := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "message", Value: msg},
		}},
	}

	_, err := coll.UpdateOne(ctx, bson.M{"subID": subID, "messageID": msgID}, updateField)
	errorhandling.Check(err)

	return err
}

func newRepository() repoInterfacer {
	return &repoStruct{
		DBTable: "discord",
	}
}
