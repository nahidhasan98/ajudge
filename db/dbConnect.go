package db

import (
	"context"
	"time"

	// "vault"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Connect function for connenting to DB
func Connect() (*mongo.Database, context.Context, context.CancelFunc) {
	// dbUser := vault.DatabaseUsername
	// dbPass := vault.DatabasePassword
	// dbName := vault.DatabaseName

	//this is mongoDB atlas connection string
	//connectionString := "mongodb+srv://" + dbUser + ":" + dbPass + "@testcluster.kwwik.gcp.mongodb.net/" + dbName + "?retryWrites=true&w=majority"

	//this is mongoDB local connection string
	connectionString := "mongodb://localhost:27017"

	dbClient, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	errorhandling.Check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbClient.Connect(ctx)
	errorhandling.Check(err)

	err = dbClient.Ping(ctx, readpref.Primary())
	errorhandling.Check(err)

	//return db
	return dbClient.Database("ajudge"), ctx, cancel
}
