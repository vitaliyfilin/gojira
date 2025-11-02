package load

import (
	"flag"
	"fmt"
	"time"
)

type LoadTestConfig struct {
	URL         string
	Method      string
	Concurrency int
	Duration    time.Duration
	BodyFile    string
	ReportFile  string
}

func ParseFlags() (*LoadTestConfig, error) {
	url := flag.String("url", "", "Target URL (required)")
	method := flag.String("method", "GET", "HTTP method")
	concurrency := flag.Int("c", 10, "Concurrency level")
	durationSec := flag.Int("t", 10, "Duration in seconds")
	bodyFile := flag.String("body-file", "", "Template file for body")
	reportFile := flag.String("report", "gojira_report.html", "HTML report output")
	flag.Parse()

	if *url == "" {
		return nil, fmt.Errorf("-url is required")
	}

	return &LoadTestConfig{
		URL:         *url,
		Method:      *method,
		Concurrency: *concurrency,
		Duration:    time.Duration(*durationSec) * time.Second,
		BodyFile:    *bodyFile,
		ReportFile:  *reportFile,
	}, nil
}
