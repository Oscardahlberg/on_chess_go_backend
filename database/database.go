package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"time"
)

var Client *mongo.Client
var Database *mongo.Database
var OpenLobbys *mongo.Collection
var OngoingLobbys *mongo.Collection
var ArchivedLobbys *mongo.Collection

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dbConnectionString := fmt.Sprintf("mongodb://127.0.0.1:27017/on_chess")

	var err error
	Client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(dbConnectionString))
	if err != nil {
		return err
	}
	Database := Client.Database("on_chess")
	OpenLobbys = Database.Collection("open_lobbys")
	OngoingLobbys = Database.Collection("ongoing_lobbys")
	ArchivedLobbys = Database.Collection("archived_lobbys")

	return nil
}
