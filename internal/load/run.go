package load

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func Run() error {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("[INFO] Starting load test...")

	config, err := ParseFlags()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Target: %s | Method: %s | Concurrency: %d | Duration: %s",
		config.URL, config.Method, config.Concurrency, config.Duration)

	tmpl, err := LoadBodyTemplate(config.BodyFile)
	if err != nil {
		log.Printf("[ERROR] Failed to load body template: %v", err)
		return err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	metrics := NewMetrics()
	var wg sync.WaitGroup

	log.Printf("[INFO] Launching %d workers...", config.Concurrency)
	stopAt := time.Now().Add(config.Duration)
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorker(id, client, config.URL, config.Method, tmpl, metrics, stopAt)
		}(i + 1)
	}

	wg.Wait()
	log.Println("[INFO] All workers finished. Generating summary...")

	summary := metrics.Summary()
	pct := metrics.Percentiles([]float64{50, 90, 95, 99})
	log.Printf("[INFO] Results - Total: %d | Success: %d | Failure: %d",
		summary["total"], summary["success"], summary["failure"])
	for p, v := range pct {
		log.Printf("[INFO] Latency p%.0f: %.2f ms", p, v)
	}

	return GenerateReport(config.ReportFile, summary, config.Duration, metrics.Latencies)
}
