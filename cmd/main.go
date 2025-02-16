package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Start metrics and health server with graceful shutdown
	startServerWithGracefulShutdown()
}

func startServerWithGracefulShutdown() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  trashHoldTen,
		WriteTimeout: trashHoldTen,
		IdleTimeout:  trashHoldFifteen,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting metrics and health server on port 8080...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}
