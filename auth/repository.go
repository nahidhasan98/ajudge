package auth

import (
	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/model"
	"go.mongodb.org/mongo-driver/bson"
)

type repoInterfacer interface {
	getUserByUsername(username string) (*model.UserData, error)
}

type repo struct {
	DBTable string
}

func (r *repo) getUserByUsername(username string) (*model.UserData, error) {
	//connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	//taking DB collection/table to a variable
	userCollection := DB.Collection(r.DBTable)

	//getting original password for this user from DB
	var userData model.UserData
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userData)

	if err != nil { //if username not exist/found in DB (returned no document/row)
		return nil, err
	}

	return &userData, nil
}

func newRepository() repoInterfacer {
	return &repo{
		DBTable: "user",
	}
}
