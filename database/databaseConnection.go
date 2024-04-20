package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {

	var mongo_uri string = "mongodb://localhost:27017/"

    // Create a new MongoDB client with the specified URI
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongo_uri))
    if err != nil {
        log.Fatal(err)
    }

    // Context with timeout for the connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connect to MongoDB
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")

    // Return the MongoDB client
    return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("whispr").Collection(collectionName)
	return collection
}