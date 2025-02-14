package main

import (
	"fmt"
	"log"
	"net/http"

	"TickFlow/internal/binance"
	"TickFlow/internal/database"
	"TickFlow/internal/observer"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "TickFlow/internal/metrics"
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
	go func() {
		if err := binance.Connect(subject); err != nil {
			log.Fatalf("Binance connection failure: %v", err)
		}
	}()

	// start prometheus metrics server
	http.Handle("/metrics", promhttp.Handler())
	prometheusPort := "8080"
	log.Printf("Starting metrics server on port %s ...", prometheusPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", prometheusPort), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
