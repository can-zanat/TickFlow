package main

import (
	"log"

	"TickFlow/internal/binance"
	"TickFlow/internal/database"
	"TickFlow/internal/observer"
)

func main() {
	// connect mongoDB
	dbInstance := database.ConnectMongo()

	// create subject for observer pattern
	subject := observer.NewSubject()

	// create trade observer and attach to subject for listening
	tradeObs := observer.NewTradeObserver(dbInstance)
	subject.Attach(tradeObs)

	// connect binance websocket and start listening for trades and notify observers
	if err := binance.Connect(subject); err != nil {
		log.Fatalf("Binance connection failure: %v", err)
	}
}
