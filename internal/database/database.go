package database

import (
	"TickFlow/configs"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TenSecondsTimeout = 10 * time.Second
const FiveSecondsTimeout = 5 * time.Second

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
	ctx, cancel := context.WithTimeout(context.Background(), TenSecondsTimeout)

	config, err := configs.LoadConfig()
	if err != nil {
		cancel()
		log.Fatalf("Error loading config: %v", err)
	}

	clientOptions := options.Client().ApplyURI(config.MongoDB.URI)
	client, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		cancel()
		log.Fatalf("MongoDB connection failure: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		cancel()
		log.Fatalf("Unable to ping MongoDB: %v", err)
	}

	DB = client.Database("rates")

	log.Println("Connected to MongoDB properly, using 'rates' database.")
	cancel()

	return &MongoDB{db: DB}
}

// SaveTrade adds trade data to the 'rateBTC/USDT' collection.
// This method now implements the Database interface.
func (m *MongoDB) SaveTrade(trade map[string]interface{}) error {
	collection := m.db.Collection("rateBTC/USDT")
	ctx, cancel := context.WithTimeout(context.Background(), FiveSecondsTimeout)

	defer cancel()

	_, err := collection.InsertOne(ctx, trade)
	if err != nil {
		log.Printf("Data cannot be added to MongoDB: %v", err)

		return err
	}

	log.Println("Data added to MongoDB.")

	return nil
}
