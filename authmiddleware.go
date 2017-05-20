package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
)

type AuthFunc func(ctx *fasthttp.RequestCtx) bool

func NewAuthMiddleware(authFunc AuthFunc) Middleware {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if authFunc(ctx) {
				h(ctx)
			} else {
				ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			}
		}
	}
}
