package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDb() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	Mongodb := os.Getenv("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(Mongodb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

var Client *mongo.Client = InitDb()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database(os.Getenv("DATABASE_NAME")).Collection(collectionName)
	return collection
}
