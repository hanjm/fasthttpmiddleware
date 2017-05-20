package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"time"
	"go.uber.org/zap/zapcore"
)

// NewLogMiddleware return a middleware which log code(status code), time(response time), method(request method), path(request url ath), addr(remote addr).
// if status code is 2xx, the log level is info, otherwise warn.
// if your app is behind of Nginx, you you may meed to set xRealIp to True so that get a actual remoteAdr.
func NewLogMiddleware(logger *zap.Logger, xRealIp bool) Middleware {
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
			if ctx.Response.StatusCode()/100 == 2 {
				logger.Info("access", zap.Int("code", ctx.Response.StatusCode()), zap.Duration("time", time.Since(startTime)), zap.ByteString("method", ctx.Method()), zap.ByteString("path", ctx.Path()), addrField)
			} else {
				logger.Warn("access", zap.Int("code", ctx.Response.StatusCode()), zap.Duration("time", time.Since(startTime)), zap.ByteString("method", ctx.Method()), zap.ByteString("path", ctx.Path()), addrField)
			}
		}
	}
}
