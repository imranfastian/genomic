package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/imranfastian/genomic/config"
	"github.com/imranfastian/genomic/handlers"
	"github.com/imranfastian/genomic/middleware"
)

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix               // check :  if this is the desired
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}) // check :  for human-readable during dev
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// ObservabilityMiddleware instruments and logs each request
func ObservabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusResponseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)
		dur := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		status := rw.status

		httpRequestsTotal.WithLabelValues(path, method, http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(path, method).Observe(dur)

		log.Info().
			Str("path", path).
			Str("method", method).
			Int("status", status).
			Float64("duration_s", dur).
			Str("remote", r.RemoteAddr).
			Msg("http request")
	})
}

// helper to capture status code
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func main() {
	// Initialize PostgreSQL connection
	if err := config.InitDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}
	defer config.CloseDB()

	mux := http.NewServeMux()

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Public route
	mux.HandleFunc("/login", handlers.LoginHandler)

	// Protected route (uses JWT middleware)
	mux.HandleFunc("/genomes", middleware.JWTMiddleware(handlers.GenomesHandler))

	// Server config
	server := &http.Server{
		Addr:         ":8080",
		Handler:      ObservabilityMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Server is running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
