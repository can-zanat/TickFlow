package observer

import (
	"TickFlow/internal/database"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	tradeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "trade_processing_duration_seconds",
		Help:    "WebSocket message to MongoDB write duration",
		Buckets: prometheus.DefBuckets,
	})
	tradeSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "trade_success_total",
		Help: "Başarılı olarak MongoDB'ye eklenen trade sayısı",
	})
	tradeErrorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "trade_error_total",
		Help: "Hata ile MongoDB'ye eklenemeyen trade sayısı",
	})
)

func init() {
	prometheus.MustRegister(tradeDuration)
	prometheus.MustRegister(tradeSuccessCounter, tradeErrorCounter)
}

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
	startTime := time.Now()
	err := t.db.SaveTrade(data)
	duration := time.Since(startTime).Seconds()
	tradeDuration.Observe(duration)

	if err != nil {
		tradeErrorCounter.Inc()
		log.Fatalf("Update function error: %v", err)

		return
	}

	tradeSuccessCounter.Inc()
}
