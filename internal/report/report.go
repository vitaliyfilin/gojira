package report

import (
	"encoding/json"
	"html/template"
	"os"
	"strings"
	"time"
)

const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8"/>
<title>Gojira Report</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<style>
body {font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 30px; max-width: 1000px;}
h1 {color: #333;}
h2 {margin-top: 40px; color: #444;}
.summary-table, .percentile-table {
  border-collapse: collapse;
  margin-top: 15px;
  width: 100%;
}
.summary-table th, .summary-table td, .percentile-table th, .percentile-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}
.summary-table tr:nth-child(even), .percentile-table tr:nth-child(even) {background-color: #f9f9f9;}
.summary-table th, .percentile-table th {background-color: #007bff; color: white;}
.chart-container {width: 100%; height: 400px; margin-top: 30px;}
footer {margin-top: 50px; font-size: 0.9em; color: #777;}
</style>
</head>
<body>
  <h1>Gojira Load Test Report</h1>

  <h2>Summary</h2>
  <table class="summary-table">
    <tr><th>Metric</th><th>Value</th></tr>
    <tr><td>Total Requests</td><td>{{index .Summary "total"}}</td></tr>
    <tr><td>Successful Requests</td><td>{{index .Summary "success"}}</td></tr>
    <tr><td>Failed Requests</td><td>{{index .Summary "failure"}}</td></tr>
    <tr><td>Total Data Transferred</td><td>{{printf "%.2f MB" (divMB (index .Summary "bytes"))}}</td></tr>
    <tr><td>Test Duration</td><td>{{index .Summary "duration"}}</td></tr>
    <tr><td>Average Latency (ms)</td><td>{{printf "%.2f" (index .Summary "avg_latency")}}</td></tr>
    <tr><td>Requests Per Second</td><td>{{printf "%.2f" (index .Summary "rps")}}</td></tr>
  </table>

  <h2>Latency Percentiles (ms)</h2>
  <table class="percentile-table">
    <tr><th>Percentile</th><th>Latency</th></tr>
    {{range $p, $v := .Percentiles}}
      <tr><td>p{{$p}}</td><td>{{printf "%.2f" $v}}</td></tr>
    {{end}}
  </table>

  <h2>Latency Distribution</h2>
  <div class="chart-container">
    <canvas id="histogram"></canvas>
  </div>

  <h2>Latency over Time</h2>
  <div class="chart-container">
    <canvas id="chart"></canvas>
  </div>

  <footer>Generated at {{.GeneratedAt}} by Gojira Load Tester</footer>

  <script>
  const latencies = {{ toJSON .Latencies }};
  const ctx1 = document.getElementById('chart').getContext('2d');
  const ctx2 = document.getElementById('histogram').getContext('2d');

  // Line chart (latency over time)
  new Chart(ctx1, {
    type: 'line',
    data: {
      labels: latencies.map((_, i) => i),
      datasets: [{
        label: 'Latency (ms)',
        data: latencies,
        borderColor: '#007bff',
        backgroundColor: 'rgba(0,123,255,0.1)',
        pointRadius: 0,
      }]
    },
    options: {
      scales: {
        x: {display: true, title: {display: true, text: 'Request #'}},
        y: {display: true, title: {display: true, text: 'Latency (ms)'}}
      }
    }
  });

  // Histogram (distribution)
  const bins = Array(20).fill(0);
  const max = Math.max(...latencies);
  for (const l of latencies) {
    const idx = Math.min(Math.floor(l / max * bins.length), bins.length - 1);
    bins[idx]++;
  }
  new Chart(ctx2, {
    type: 'bar',
    data: {
      labels: bins.map((_, i) => Math.round(i / bins.length * max)),
      datasets: [{
        label: 'Count',
        data: bins,
        backgroundColor: 'rgba(0, 123, 255, 0.6)'
      }]
    },
    options: {
      scales: {
        x: {title: {display: true, text: 'Latency (ms)'}},
        y: {title: {display: true, text: 'Count'}}
      }
    }
  });
  </script>
</body>
</html>`

func GenerateHTML(file string, data map[string]interface{}) error {
	t := template.Must(template.New("r").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			s := strings.ReplaceAll(string(b), "</", "<\\/")
			return template.JS(s)
		},
		"divMB": func(v interface{}) float64 {
			switch n := v.(type) {
			case uint64:
				return float64(n) / (1024 * 1024)
			case int64:
				return float64(n) / (1024 * 1024)
			default:
				return 0
			}
		},
	}).Parse(tpl))

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	data["GeneratedAt"] = time.Now().Format(time.RFC1123)
	return t.Execute(f, data)
}
