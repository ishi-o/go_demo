package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	requestsInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Number of HTTP requests in progress",
		},
		[]string{"method", "path"},
	)
)

func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		method := string(c.Method())
		path := c.URI().String()

		requestsInProgress.WithLabelValues(method, path).Inc()
		c.Next(ctx)
		requestsInProgress.WithLabelValues(method, path).Dec()

		duration := time.Since(start)
		status := c.Response.StatusCode()
		statusStr := strconv.Itoa(status)
		requestsTotal.WithLabelValues(method, path, statusStr).Inc()
		requestDuration.WithLabelValues(method, path, statusStr).Observe(duration.Seconds())
	}
}
