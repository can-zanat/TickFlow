package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const trashHold = 10 * time.Second

var (
	memAllocGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_mem_alloc_bytes",
		Help: "Current memory allocation",
	})
)

func init() {
	prometheus.MustRegister(memAllocGauge)

	monitorMemStats := func() {
		var memStats runtime.MemStats

		for {
			runtime.ReadMemStats(&memStats)
			memAllocGauge.Set(float64(memStats.Alloc))
			time.Sleep(trashHold)
		}
	}

	go monitorMemStats()
}
