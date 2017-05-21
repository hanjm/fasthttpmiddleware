package fasthttpmiddleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

var (
	promLabelNames = []string{"code", "method", "path"}
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts(prometheus.Opts{
			Name: "request_count",
			Help: "request count",
		}), promLabelNames)
	responseTimeSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "response_time",
		Help: "response time (second)",
	}, promLabelNames)
)

// NewPrometheusMiddleware return a middleware which can be used by [prometheus](https://github.com/prometheus/prometheus)
// The prometheus is a monitoring system and time series database.
// Note: the returned middleware is contain the function of logmiddleware.
func NewPrometheusMiddleware(bindAddr string, xRealIp bool, logger *zap.Logger) Middleware {
	go func() {
		prometheus.MustRegister(requestCounter)
		prometheus.MustRegister(responseTimeSummary)
		http.Handle("/metrics", promhttp.Handler())
		logger.Debug("prometheus metrics server start at " + bindAddr)
		http.ListenAndServe(bindAddr, nil)
	}()
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			startTime := time.Now()
			h(ctx)
			var addrField zapcore.Field
			if xRealIp {
				addrField = zap.ByteString("addr", ctx.Request.Header.Peek("X-Real-IP"))
			} else {
				addrField = zap.String("addr", ctx.RemoteAddr().String())
			}
			promLabels := prometheus.Labels{"code": strconv.Itoa(ctx.Response.StatusCode()), "method": string(ctx.Method()), "path": string(ctx.Path())}
			responseTime := time.Since(startTime).Seconds()
			responseTimeSummary.With(promLabels).Observe(responseTime)
			requestCounter.With(promLabels).Inc()
			if ctx.Response.StatusCode()/100 == 2 {
				logger.Info("access", zap.Int("code", ctx.Response.StatusCode()), zap.Float64("time", responseTime), zap.String("method", promLabels["method"]), zap.String("path", promLabels["path"]), addrField)
			} else {
				logger.Warn("access", zap.Int("code", ctx.Response.StatusCode()), zap.Float64("time", responseTime), zap.String("method", promLabels["method"]), zap.String("path", promLabels["path"]), addrField)
			}
		}
	}
}
