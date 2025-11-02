package load

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type RequestResult struct {
	Started   time.Time
	Completed time.Time
	Latency   time.Duration
	Status    int
	Err       string
	Bytes     int64
}

type Metrics struct {
	TotalRequests uint64
	Success       uint64
	Failure       uint64
	Bytes         uint64
	LatenciesMu   sync.Mutex
	Latencies     []float64
	Results       []RequestResult
}

func NewMetrics() *Metrics {
	return &Metrics{
		Latencies: make([]float64, 0, 1000),
	}
}

func (m *Metrics) Record(rr RequestResult) {
	atomic.AddUint64(&m.TotalRequests, 1)
	if rr.Err == "" && rr.Status >= 200 && rr.Status < 400 {
		atomic.AddUint64(&m.Success, 1)
	} else {
		atomic.AddUint64(&m.Failure, 1)
	}
	atomic.AddUint64(&m.Bytes, uint64(rr.Bytes))

	m.LatenciesMu.Lock()
	m.Latencies = append(m.Latencies, rr.Latency.Seconds()*1000.0)
	m.Results = append(m.Results, rr)
	m.LatenciesMu.Unlock()
}

func (m *Metrics) Percentiles(ps []float64) map[float64]float64 {
	m.LatenciesMu.Lock()
	defer m.LatenciesMu.Unlock()

	out := map[float64]float64{}
	if len(m.Latencies) == 0 {
		for _, p := range ps {
			out[p] = 0.0
		}
		return out
	}

	arr := append([]float64(nil), m.Latencies...)
	sort.Float64s(arr)

	for _, p := range ps {
		if p <= 0 {
			out[p] = arr[0]
			continue
		}
		if p >= 100 {
			out[p] = arr[len(arr)-1]
			continue
		}
		k := int(math.Ceil((p / 100.0) * float64(len(arr))))
		if k <= 0 {
			k = 1
		}
		out[p] = arr[k-1]
	}
	return out
}

func (m *Metrics) Summary() map[string]uint64 {
	return map[string]uint64{
		"total":   m.TotalRequests,
		"success": m.Success,
		"failure": m.Failure,
	}
}
