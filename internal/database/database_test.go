package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStore struct {
	Client *mongo.Client
}

const (
	mongoImage = "mongo:7.0.4"
)

func NewStoreWithURI(uri string) *MongoDBStore {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("Connection failure:", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Unable to access MongoDB:", err)
	}

	return &MongoDBStore{
		Client: client,
	}
}

func prepareTestStore(t *testing.T) (store *MongoDBStore, clean func()) {
	t.Helper()

	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage(mongoImage))
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}

	clean = func() {
		if terminateErr := mongodbContainer.Terminate(ctx); terminateErr != nil {
			t.Fatalf("Failed to terminate MongoDB container: %v", terminateErr)
		}
	}

	containerURI, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get container connection string: %v", err)
	}

	s := NewStoreWithURI(containerURI)

	return s, clean
}

func TestMongoDBStore_SaveTrade(t *testing.T) {
	store, clean := prepareTestStore(t)
	defer clean()

	mongoInstance := &MongoDB{db: store.Client.Database("rates")}

	oid, err := primitive.ObjectIDFromHex("67af25ccf98bf2e0fbfb2095")
	if err != nil {
		t.Fatalf("Invalid ObjectID: %v", err)
	}

	trade := map[string]interface{}{
		"_id": oid,
		"Q":   "0.01489000",
		"b":   "97076.93000000",
		"h":   "97248.60000000",
		"p":   "830.69000000",
		"w":   "96198.67603745",
		"B":   "9.81902000",
		"o":   "96246.24000000",
		"l":   "95217.36000000",
		"q":   "1643937768.37928790",
		"n":   3232443,
		"s":   "BTCUSDT",
		"P":   "0.863",
		"x":   "96246.24000000",
		"c":   "97076.93000000",
		"O":   int64(1739445324865),
		"F":   int64(4554306623),
		"E":   int64(1739531724866),
		"a":   "97076.94000000",
		"A":   "4.68856000",
		"v":   "17088.98538000",
		"C":   int64(1739531724865),
		"L":   int64(4557539065),
		"e":   "24hrTicker",
	}

	if err = mongoInstance.SaveTrade(trade); err != nil {
		t.Fatalf("SaveTrade failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result map[string]interface{}

	filter := bson.M{"_id": oid}

	if err = mongoInstance.db.Collection("rateBTC/USDT").FindOne(ctx, filter).Decode(&result); err != nil {
		t.Fatalf("Cannot find data: %v", err)
	}

	if result["s"] != "BTCUSDT" {
		t.Errorf("Expected s: BTCUSDT, got: %v", result["s"])
	}
}
