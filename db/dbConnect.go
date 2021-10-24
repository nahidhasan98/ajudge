package db

import (
	"context"
	"time"

	// "vault"

	"github.com/nahidhasan98/ajudge/errorhandling"
	"github.com/nahidhasan98/ajudge/vault"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Connect function for connenting to DB
func Connect() (*mongo.Database, context.Context, context.CancelFunc) {
	dbName := vault.DatabaseName
	dbConnectionString := vault.DatabaseConnectionString

	dbClient, err := mongo.NewClient(options.Client().ApplyURI(dbConnectionString))
	errorhandling.Check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbClient.Connect(ctx)
	errorhandling.Check(err)

	err = dbClient.Ping(ctx, readpref.Primary())
	errorhandling.Check(err)

	//return db
	return dbClient.Database(dbName), ctx, cancel
}
