package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"note_app_server/global"
)

func InitMongoDB() {
	ctx := context.Background()
	username := AC.Mongo.Username
	password := AC.Mongo.Password
	host := AC.Mongo.Host
	port := AC.Mongo.Port
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+username+":"+password+"@"+host+":"+port))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	global.MongoClient = client
}
