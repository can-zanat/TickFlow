package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	memAllocGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_alloc_bytes",
		Help: "Current memory allocation",
	})
)

func init() {
	prometheus.MustRegister(memAllocGauge)
	go monitorMemStats()
}

func monitorMemStats() {
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)
		memAllocGauge.Set(float64(memStats.Alloc))
		time.Sleep(10 * time.Second)
	}
}
