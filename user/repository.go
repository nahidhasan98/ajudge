package user

import (
	"github.com/nahidhasan98/ajudge/db"
	"github.com/nahidhasan98/ajudge/errorhandling"
	"go.mongodb.org/mongo-driver/bson"
)

type repoInterfacer interface {
	getUserByEmail(email string) (*UserModel, error)
	getUserByUsername(username string) (*UserModel, error)
	getLastUserID() int
	updateLastUserID() error
	createUser(userData *UserModel) error
}

type repo struct {
	DBTable string
}

func (r *repo) getUserByEmail(email string) (*UserModel, error) {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	userCollection := DB.Collection(r.DBTable)

	// getting original password for this user from DB
	var userData UserModel
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&userData)
	if err != nil { // if username not exist/found in DB (returned no document/row)
		return nil, err
	}

	return &userData, nil
}

func (r *repo) getUserByUsername(username string) (*UserModel, error) {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	userCollection := DB.Collection(r.DBTable)

	// getting original password for this user from DB
	var userData UserModel
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&userData)
	if err != nil { // if username not exist/found in DB (returned no document/row)
		return nil, err
	}

	return &userData, nil
}

func (r *repo) getLastUserID() int {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	counterCollection := DB.Collection("counter")

	// getting LastUserID from DB to assign a user ID(LastUserID+1) for this user
	var lastUsedID lastUsedIDModel
	err := counterCollection.FindOne(ctx, bson.M{}).Decode(&lastUsedID)
	errorhandling.Check(err)

	return lastUsedID.LastUserID
}

func (r *repo) updateLastUserID() error {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	counterCollection := DB.Collection("counter")

	// updating LastUserID to DB for later use/next user
	updateField := bson.D{
		{Key: "$inc", Value: bson.D{ // incrementing LastUserID by 1
			{Key: "lastUserID", Value: 1},
		}},
	}
	_, err := counterCollection.UpdateOne(ctx, bson.M{}, updateField)
	return err
}

func (r *repo) createUser(userData *UserModel) error {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	userCollection := DB.Collection(r.DBTable)

	// getting original password for this user from DB
	_, err := userCollection.InsertOne(ctx, userData)
	errorhandling.Check(err)

	return err
}

func newRepository() repoInterfacer {
	return &repo{
		DBTable: "user",
	}
}
