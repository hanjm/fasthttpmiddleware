package fasthttpmiddleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

var (
	isRegistered = false // fix when to invoke NewPrometheusMiddleware twice lead a panic: duplicate metrics collector registration attempted
	promLabelNames = []string{"code", "method", "path"}
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts(prometheus.Opts{
			Name: "http_requests_total",
			Help: "http requests total",
		}), promLabelNames)
	responseTimeSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_request_duration_microseconds",
		Help: "http request duration microseconds",
	}, promLabelNames)
)

// NewPrometheusMiddleware return a middleware which can be used by [prometheus](https://github.com/prometheus/prometheus) collecting metrics.
// The prometheus is a monitoring system and time series database.
func NewPrometheusMiddleware(bindAddr string, logger *zap.Logger) Middleware {
	go func() {
		if !isRegistered {
			prometheus.MustRegister(requestCounter)
			prometheus.MustRegister(responseTimeSummary)
			isRegistered = true
			http.Handle("/metrics", promhttp.Handler())
			logger.Debug("prometheus metrics server start at " + bindAddr)
			http.ListenAndServe(bindAddr, nil)
		}
	}()
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			startTime := time.Now()
			h(ctx)
			promLabels := prometheus.Labels{"code": strconv.Itoa(ctx.Response.StatusCode()), "method": string(ctx.Method()), "path": string(ctx.Path())}
			responseTime := time.Since(startTime).Seconds() * 1000
			responseTimeSummary.With(promLabels).Observe(responseTime)
			requestCounter.With(promLabels).Inc()
		}
	}
}
