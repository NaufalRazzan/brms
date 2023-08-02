package db

import (
	"brms/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB(collectionName string) (*mongo.Client, *mongo.Collection, error){
	client, err := mongo.NewClient(options.Client().ApplyURI(config.GetConfig().MongoDBUrl))
	if err != nil{
		log.Panicln("Error while creating new client: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil{
		log.Panicln("Error while establishing a connection to DB: ", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil{
		log.Panicln("Error while pinging: ", err)
	}

	log.Println("Connected to MongoDB")
	
	collection := client.Database(config.GetConfig().MongoDBname).Collection(collectionName)

	return client, collection, nil
}