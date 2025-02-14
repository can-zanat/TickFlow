package database

import (
	"context"
	"log"
	"time"

	"TickFlow/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database interface exists for testing and dependency injection.
type Database interface {
	SaveTrade(data map[string]interface{}) error
}

var (
	DB     *mongo.Database
	client *mongo.Client
)

type MongoDB struct {
	db *mongo.Database
}

// ConnectMongo starts the MongoDB connection.
func ConnectMongo() *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	clientOptions := options.Client().ApplyURI(config.MongoDB.URI)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection failure: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Unable to ping MongoDB: %v", err)
	}

	DB = client.Database("rates")
	log.Println("Connected to MongoDB properly, using 'rates' database.")

	return &MongoDB{db: DB}
}

// SaveTrade adds trade data to the 'rateBTC/USDT' collection.
// This method now implements the Database interface.
func (m *MongoDB) SaveTrade(trade map[string]interface{}) error {
	collection := m.db.Collection("rateBTC/USDT")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, trade)
	if err != nil {
		log.Printf("Data cannot be added to MongoDB: %v", err)
		return err
	}
	log.Println("Data added to MongoDB.")
	return nil
}
