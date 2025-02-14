package observer

import (
	"TickFlow/internal/database"
	"log"
)

// TradeObserver get the trade data and write it to MongoDB.
type TradeObserver struct {
	db database.Database
}

// NewTradeObserver create a new TradeObserver with dependency injection.
func NewTradeObserver(db database.Database) *TradeObserver {
	return &TradeObserver{db: db}
}

// Update get data from subject and write it to MongoDB using the injected database.
func (t *TradeObserver) Update(data map[string]interface{}) {
	err := t.db.SaveTrade(data)
	if err != nil {
		log.Fatalf("Update function error: %v", err)
		return
	}
}
