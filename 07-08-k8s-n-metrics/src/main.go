package main

import (
        "github.com/gorilla/mux"
        "math/rand"
        "net/http"
        "strconv"
        "time"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        counter = promauto.NewCounterVec(prometheus.CounterOpts{
                Name: "myapp_counts_total",
                Help: "The total number of somethins",
        }, []string{"code", "method", "user_id"})
        gauge = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "myapp_gauge_total",
                Help: "Current level of some metric",
                ConstLabels: prometheus.Labels{"session": "design_system"},
        })
        histogram = promauto.NewHistogram(prometheus.HistogramOpts{
                Name:    "myapp_random_numbers_hist",
                Help:    "A histogram of normally distributed random numbers.",
                Buckets: prometheus.LinearBuckets(-3, .1, 61), // from -3 with step 0.1
        })
        summary = promauto.NewSummary(prometheus.SummaryOpts{
                Name:    "myapp_random_numbers_sum",
                Help:    "A summary of normally distributed random numbers.",
                Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, //quantiles with absolute errors
        })
)

func recordMetrics() {
        go func() {
                for {
                        counter.WithLabelValues("404", "POST").Add(float64(rand.Intn(10)))
                        time.Sleep(100 * time.Millisecond)
                }
        }()
        go func() {
                for {
                        gauge.Set(float64(rand.Intn(10) + 100))
                        time.Sleep(100 * time.Millisecond)
                }
        }()
        go func() {
                for {
                        histogram.Observe(rand.NormFloat64())
                        time.Sleep(100 * time.Millisecond)
                }
        }()
        go func() {
                for {
                        summary.Observe(rand.NormFloat64())
                        time.Sleep(100 * time.Millisecond)
                }
        }()
}

var (
        httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
                Name: "myapp_http_duration_seconds",
        }, []string{"path"})

        httpCodes = promauto.NewCounterVec(prometheus.CounterOpts{
                Name: "myapp_http_status_codes_total",
        }, []string{"path", "code"})
)

func handleFast(w http.ResponseWriter, r *http.Request) {
        _, err := w.Write([]byte("hello"))
        if err != nil {
                http.Error(w, err.Error(), 400)
        }
}

func handleRandom(w http.ResponseWriter, r *http.Request) {
        time.Sleep(100 * time.Duration(rand.Intn(100))  * time.Millisecond)
        _, err := w.Write([]byte("biba"))
        if err != nil {
                http.Error(w, err.Error(), 400)
        }
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
        statuses := []int{
                http.StatusOK,
                http.StatusAccepted,
                http.StatusCreated,
                http.StatusBadRequest,
                http.StatusForbidden,
                http.StatusBadGateway,
                http.StatusInternalServerError,
        }
        w.WriteHeader(statuses[rand.Intn(len(statuses))])
}

type loggingResponseWriter struct {
        http.ResponseWriter
        statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
        return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
        lrw.statusCode = code
        lrw.ResponseWriter.WriteHeader(code)
}

func prometheusMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                route := mux.CurrentRoute(r)
                path, _ := route.GetPathTemplate()
                timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

                lrw := NewLoggingResponseWriter(w)
                next.ServeHTTP(lrw, r)

                timer.ObserveDuration()
                statusCode := lrw.statusCode
                httpCodes.WithLabelValues(path, strconv.Itoa(statusCode)).Inc()
        })
}

func main() {
        r := mux.NewRouter()
        r.Use(prometheusMiddleware)
        recordMetrics()

        r.Path("/metrics").Handler(promhttp.Handler())
        r.Path("/api/fast").HandlerFunc(handleFast)
        r.Path("/api/random").HandlerFunc(handleRandom)
        r.Path("/api/status").HandlerFunc(handleStatus)

        srv := &http.Server{Addr: "0.0.0.0:2112", Handler: r}
        srv.ListenAndServe()
}