package binance

import (
	"TickFlow/configs"
	"log"
	"net/http"

	"TickFlow/internal/observer"
	"github.com/gorilla/websocket"
)

// Connect to Binance WebSocket and start listening for trades.
func Connect(subject *observer.Subject) error {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(config.BinanceURL, http.Header{})
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Binance WebSocket connection create properly!")

	for {
		var message map[string]interface{}
		err = conn.ReadJSON(&message)
		if err != nil {
			log.Println("JSON parse error:", err)
			break
		}
		subject.Notify(message)
	}

	return nil
}
