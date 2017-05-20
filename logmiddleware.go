package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"time"
)

func NewLogMiddleware(logger *zap.Logger) Middleware {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			startTime := time.Now()
			h(ctx)
			if ctx.Response.StatusCode()/100 == 2 {
				logger.Info("access", zap.Int("code", ctx.Response.StatusCode()), zap.Duration("time", time.Since(startTime)), zap.ByteString("method", ctx.Method()), zap.ByteString("path", ctx.Path()), zap.String("addr", ctx.RemoteAddr().String()))
			} else {
				logger.Warn("access", zap.Int("code", ctx.Response.StatusCode()), zap.Duration("time", time.Since(startTime)), zap.ByteString("method", ctx.Method()), zap.ByteString("path", ctx.Path()), zap.String("addr", ctx.RemoteAddr().String()))
			}
		}
	}
}
