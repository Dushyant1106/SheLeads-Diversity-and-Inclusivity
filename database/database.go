package database

import (
	"context"
	"log"
	"sheleads-backend/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURL)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	DB = Client.Database(config.AppConfig.DBName)
	log.Println("MongoDB connected successfully to database:", config.AppConfig.DBName)
}

func Migrate() {
	// MongoDB doesn't require migrations like SQL databases
	// Collections are created automatically when first document is inserted
	// We can create indexes here if needed

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create indexes for users collection
	usersCollection := DB.Collection("users")
	_, err := usersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("Warning: Failed to create email index:", err)
	}

	// Create indexes for worklogs collection
	worklogsCollection := DB.Collection("worklogs")
	_, err = worklogsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"user_id": 1, "created_at": -1},
	})
	if err != nil {
		log.Println("Warning: Failed to create worklogs index:", err)
	}

	log.Println("MongoDB indexes created successfully")
}

func Disconnect() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			log.Println("MongoDB disconnected successfully")
		}
	}
}

