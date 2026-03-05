package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

// NewMetrics returns a middleware that records per-route Prometheus metrics and
// registers the /metrics endpoint on the provided app.
func NewMetrics(app *fiber.App) fiber.Handler {
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())

		// Use the route pattern (e.g. /example/:id) instead of the actual path
		// to avoid label cardinality explosion in Prometheus.
		routePath := c.Route().Path

		httpRequestsTotal.WithLabelValues(c.Method(), routePath, status).Inc()
		httpRequestDuration.WithLabelValues(c.Method(), routePath).Observe(duration)

		return err
	}
}
