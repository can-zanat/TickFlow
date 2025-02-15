package binance

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"TickFlow/configs"
	"TickFlow/internal/observer"

	"github.com/gorilla/websocket"
)

const maxRetries = 5 // Maksimum yeniden bağlantı denemesi
const trashHold = 5 * time.Second

// Connect establishes a connection to the Binance WebSocket and returns an error
// if the maximum retry count is exceeded.
func Connect(subject *observer.Subject) error {
	config, err := configs.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	return dialAndListen(subject, config.BinanceURL, 0)
}

// dialAndListen attempts to connect to the given URL and listen for messages.
// If an error occurs, it retries with a delay until the maximum number of retries is reached.
func dialAndListen(subject *observer.Subject, url string, attempt int) error {
	if attempt >= maxRetries {
		return fmt.Errorf("maximum reconnection attempts (%d) exceeded", maxRetries)
	}

	conn, resp, err := websocket.DefaultDialer.Dial(url, http.Header{})
	if err != nil {
		log.Printf("Attempt %d: Failed to establish WebSocket connection: %v", attempt+1, err)
		time.Sleep(trashHold)

		return dialAndListen(subject, url, attempt+1)
	}

	if resp != nil {
		resp.Body.Close()
	}

	log.Println("Binance WebSocket connection established successfully!")

	// If the connection drops, handleConnection returns an error, which triggers a reconnection attempt.
	err = handleConnection(conn, subject)
	if err != nil {
		log.Printf("Connection error: %v. Retrying reconnection (Attempt %d)...", err, attempt+1)
		time.Sleep(trashHold)

		return dialAndListen(subject, url, attempt+1)
	}

	return nil
}

// handleConnection continuously reads messages from the active WebSocket connection.
// If the connection is lost, it returns an error to trigger a reconnection attempt.
func handleConnection(conn *websocket.Conn, subject *observer.Subject) error {
	defer conn.Close()

	for {
		var message map[string]interface{}
		if err := conn.ReadJSON(&message); err != nil {
			return err
		}

		subject.Notify(message)
	}
}
