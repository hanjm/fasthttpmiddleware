package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Middleware func(h fasthttp.RequestHandler) fasthttp.RequestHandler

// MiddlewareOnion represent the middleware like an onion,
// the bigger index of middleware in MiddlewareOnion.layers locate at outside
type MiddlewareOnion struct {
	layers []Middleware
}

// MiddlewareOnion returns a middleware onion with given middlewares
func NewMiddlewareOnion(middlewares ...Middleware) MiddlewareOnion {
	return MiddlewareOnion{append([]Middleware{}, middlewares...)}
}

// NewNormalMiddlewareOnion returns a normal middleware onion. recover -> auth -> log.
// the type of AuthFunc is "func(ctx *fasthttp.RequestCtx) bool"
func NewNormalMiddlewareOnion(authFunc AuthFunc, logger *zap.Logger) MiddlewareOnion {
	return MiddlewareOnion{[]Middleware{
		NewLogMiddleware(logger, true),
		NewAuthMiddleware(authFunc),
		NewRecoverMiddleware(logger),
	}}
}

func (o MiddlewareOnion) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	for i := len(o.layers) - 1; i > -1; i-- {
		h = o.layers[i](h)
	}
	return h
}

// Append copy all middleware layers to newLayers, then append middlewares to newLayers, then return a new middleware onion.
func (o MiddlewareOnion) Append(middlewares ...Middleware) MiddlewareOnion {
	return MiddlewareOnion{append(o.layers, middlewares...)}
}
