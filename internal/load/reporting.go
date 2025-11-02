package load

import (
	"log"
	"time"

	"gojira/internal/report"
)

// GenerateReport generates an HTML report with enhanced metrics.
func GenerateReport(file string, summary map[string]uint64, duration time.Duration, latencies []float64) error {
	var avgLatency float64
	for _, l := range latencies {
		avgLatency += l
	}
	if len(latencies) > 0 {
		avgLatency /= float64(len(latencies))
	}

	rps := float64(summary["total"]) / duration.Seconds()
	percentiles := map[float64]float64{}
	if len(latencies) > 0 {
		sorted := make([]float64, len(latencies))
		copy(sorted, latencies)
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i] > sorted[j] {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		for _, p := range []float64{50, 90, 95, 99} {
			k := int(float64(len(sorted)) * p / 100.0)
			if k >= len(sorted) {
				k = len(sorted) - 1
			}
			percentiles[p] = sorted[k]
		}
	}

	data := map[string]any{
		"Summary": map[string]any{
			"total":       summary["total"],
			"success":     summary["success"],
			"failure":     summary["failure"],
			"bytes":       summary["bytes"],
			"duration":    duration.String(),
			"avg_latency": avgLatency,
			"rps":         rps,
		},
		"Percentiles": percentiles,
		"Latencies":   latencies,
	}

	if err := report.GenerateHTML(file, data); err != nil {
		log.Printf("[ERROR] Failed to generate report: %v", err)
		return err
	}

	log.Printf("[INFO] Report generated: %s", file)
	return nil
}
