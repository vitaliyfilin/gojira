package load

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"
)

func runWorker(id int, client *http.Client, url, method string, tmpl *template.Template, metrics *Metrics, stopAt time.Time) {
	for time.Now().Before(stopAt) {
		bodyBytes, err := RenderBody(tmpl)
		if err != nil {
			log.Printf("[WARN] Worker %d: body template error: %v", id, err)
			continue
		}
		var body io.Reader
		if bodyBytes != nil {
			body = bytes.NewReader(bodyBytes)
		}

		start := time.Now()
		log.Printf("[INFO] Worker %d: sending %s request to %s", id, method, url)
		req, _ := http.NewRequest(method, url, body)
		resp, err := client.Do(req)
		lat := time.Since(start)

		rr := RequestResult{Started: start, Completed: start.Add(lat), Latency: lat}

		if err != nil {
			rr.Err = err.Error()
			log.Printf("[WARN] Worker %d: request error: %v", id, err)
			metrics.Record(rr)
			continue
		}

		n, _ := io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		rr.Status = resp.StatusCode
		rr.Bytes = n
		metrics.Record(rr)

		log.Printf("[INFO] Worker %d: received response with status %d, body size %d bytes", id, resp.StatusCode, n)

		if resp.StatusCode >= 400 {
			log.Printf("[WARN] Worker %d: received status %d", id, resp.StatusCode)
		}
	}
}

func StartWorkers(concurrency int, client *http.Client, url, method string, tmpl *template.Template, metrics *Metrics, duration time.Duration) {
	stopAt := time.Now().Add(duration)
	for i := 0; i < concurrency; i++ {
		go runWorker(i+1, client, url, method, tmpl, metrics, stopAt)
	}
}
