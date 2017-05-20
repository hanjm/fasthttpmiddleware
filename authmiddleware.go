package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
)

type AuthFunc func(ctx *fasthttp.RequestCtx) bool

// NewAuthMiddleware accepts a customer auth function and then returns a middleware which
// if auth function returns false, it will term the HTTP request and response 403 status code
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
