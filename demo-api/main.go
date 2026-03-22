package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
    "sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_requests_total",
	}, []string{"endpoint", "status"})

	duration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "demo_request_duration_seconds",
		Buckets: []float64{0.1, 0.5, 1, 2, 5},
	}, []string{"endpoint"})
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func track(endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		rw := &rw{ResponseWriter: w, code: 200}
		next(rw, r)
		elapsed := time.Since(t).Seconds()
		requests.WithLabelValues(endpoint, fmt.Sprint(rw.code)).Inc()
		duration.WithLabelValues(endpoint).Observe(elapsed)
		log.Info("req", "endpoint", endpoint, "status", rw.code, "ms", elapsed*1000)
	}
}

type rw struct {
	http.ResponseWriter
	code int
}

func (r *rw) WriteHeader(code int) { r.code = code; r.ResponseWriter.WriteHeader(code) }

func headers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/normal", track("/normal", func(w http.ResponseWriter, r *http.Request) {
		headers(w)
		fmt.Fprint(w, `{"status":"ok"}`)
	}))

	http.HandleFunc("/slow", track("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(2000+rand.Intn(3000)) * time.Millisecond)
		headers(w)
		fmt.Fprint(w, `{"status":"ok"}`)
	}))

	http.HandleFunc("/error", track("/error", func(w http.ResponseWriter, r *http.Request) {
		headers(w)
		w.WriteHeader(500)
		fmt.Fprint(w, `{"status":"error"}`)
	}))

	http.HandleFunc("/burst", track("/burst", func(w http.ResponseWriter, r *http.Request) {
        var wg sync.WaitGroup
        for i := 0; i < 10; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
                status := "200"
                if rand.Intn(5) == 0 {
                    status = "500"
                }
                requests.WithLabelValues("/burst-sub", status).Inc()
            }()
        }
        wg.Wait()
        headers(w)
        fmt.Fprint(w, `{"status":"ok","fired":10}`)
    }))

	log.Info("demo-api ready", "port", 8080)
	http.ListenAndServe(":8080", nil)
}