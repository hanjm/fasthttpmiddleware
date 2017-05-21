package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// NewRecoverMiddleware return a middleware which can let app recover from a panic in request handler.
// panic stack info will appear on the field named "stacktrace" in the log line
func NewRecoverMiddleware(logger *zap.Logger) Middleware {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.DPanic("recover")
					ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
			}()
			h(ctx)
		}
	}
}
