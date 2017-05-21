package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	_ "net/http/pprof"
	"net/http"
)

// NewPrometheusMiddleware return a middleware which can be used by [prometheus](https://github.com/prometheus/prometheus)
// The prometheus is a monitoring system and time series database.
func NewPrometheusMiddleware(bindAddr string) Middleware {
	go func() {
		http.ListenAndServe(bindAddr, nil)
	}()
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			h(ctx)
		}
	}
}
