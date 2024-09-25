package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() {
	err := godotenv.Load(".env")
	if err !=nil {
		log.Fatalf("Error loading environment")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err!= nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to db")
	}

}