package binance

import (
	"TickFlow/configs"
	"TickFlow/internal/observer"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Connect to Binance WebSocket and start listening for trades.
func Connect(subject *observer.Subject) error {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, resp, err := websocket.DefaultDialer.Dial(config.BinanceURL, http.Header{})
	if err != nil {
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
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
