package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

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
