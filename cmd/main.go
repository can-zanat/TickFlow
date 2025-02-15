package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"TickFlow/internal/binance"
	"TickFlow/internal/database"
	"TickFlow/internal/observer"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "TickFlow/internal/metrics"
)

const trashHoldTen = 10 * time.Second
const trashHoldFifteen = 15 * time.Second

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

	startMetricsServer()
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	prometheusPort := "8080"
	log.Printf("Starting metrics server on port %s ...", prometheusPort)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", prometheusPort),
		ReadTimeout:  trashHoldTen,
		WriteTimeout: trashHoldTen,
		IdleTimeout:  trashHoldFifteen,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
