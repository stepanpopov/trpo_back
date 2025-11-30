package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ResponseWriterStatusCodeSaver struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriterStatusCodeSaver) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriterStatusCodeSaver) StatusCode() int {
	if w.statusCode == 0 {
		return 200
	}
	return w.statusCode
}

var responseTimeMetrics = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "fluire",
		Subsystem:  "api",
		Name:       "response_time",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01},
	},
	[]string{"method", "route", "code"},
)

func observeResponseTime(duration time.Duration, method, route, code string) {
	responseTimeMetrics.WithLabelValues(method, route, code).
		Observe(float64(duration.Microseconds()))
}

// Metrics is a middleware that observes and records metrics for HTTP requests.
// It tracks the response time, HTTP method, route pattern, and status code of each request.
// The metrics are recorded using the `observeResponseTime` function.
//
// Behavior:
//   - If the route pattern is unavailable, it defaults to "NIL".
//   - The middleware wraps the ResponseWriter to capture the status code.
//
// Returns:
//   - A middleware function that wraps an http.Handler and records metrics for each request.
func Metrics() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// API route (URL Pattern)
			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern == "" {
				routePattern = "NIL"
			}

			writerSaver := &ResponseWriterStatusCodeSaver{
				ResponseWriter: w,
			}

			start := time.Now()
			defer func() {
				// Status code
				code := writerSaver.StatusCode()
				codeStr := ""
				if code != 0 {
					codeStr = strconv.Itoa(code)
				}

				observeResponseTime(time.Since(start), r.Method, routePattern, codeStr)
			}()

			next.ServeHTTP(writerSaver, r)
		})
	}
}
