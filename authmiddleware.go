package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
)

type AuthFunc func(ctx *fasthttp.RequestCtx) bool

// NewAuthMiddleware accept a customer auth function and then return a middleware which
// if auth function return false, it will term the http request and response 403 status code
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
