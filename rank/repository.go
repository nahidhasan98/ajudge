package rank

import (
	"github.com/nahidhasan98/ajudge/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repoInterfacer interface {
	getUserRank() ([]rankModel, error)
	getOJRank() ([]rankModel, error)
}

type repo struct {
	DBTableForUser string
	DBTableForOJ   string
}

func (r *repo) getUserRank() ([]rankModel, error) {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	userCollection := DB.Collection(r.DBTableForUser)

	// setting up options for retrieving data from DB
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "totalSolved", Value: -1}, {Key: "username", Value: 1}}) //sorting by totalSolved & then OJ/user name
	cursor, err := userCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var rankList []rankModel
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp rankModel
		err := cursor.Decode(&temp)
		if err != nil {
			return nil, err
		}

		rankList = append(rankList, temp)
	}

	return rankList, nil
}

func (r *repo) getOJRank() ([]rankModel, error) {
	// connecting to DB
	DB, ctx, cancel := db.Connect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	// taking DB collection/table to a variable
	rankOJCollection := DB.Collection(r.DBTableForOJ)

	// setting up options for retrieving data from DB
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "totalSolved", Value: -1}, {Key: "OJ", Value: 1}}) // sorting by totalSolved & then OJ/user name
	cursor, err := rankOJCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var rankList []rankModel
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var temp rankModel
		err := cursor.Decode(&temp)
		if err != nil {
			return nil, err
		}

		rankList = append(rankList, temp)
	}

	return rankList, nil
}

func newRepository() repoInterfacer {
	return &repo{
		DBTableForUser: "user",
		DBTableForOJ:   "rankOJ",
	}
}
