package fasthttpmiddleware

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Middleware func(h fasthttp.RequestHandler) fasthttp.RequestHandler

// MiddlewareOnion represent the middlewares like a onion,
// the bigger index of middleware in MiddlewareOnion.layers locate at outside
type MiddlewareOnion struct {
	layers []Middleware
}

// New return a middleware onion with given middlewares
func New(middlewares ...Middleware) MiddlewareOnion {
	return MiddlewareOnion{append([]Middleware{}, middlewares...)}
}

// NewNormalMiddlewareOnion return a normal middleware onion. recover -> auth -> log.
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

// Append copy all middleware layers to newLayers, then append middlewares in to newLayers
func (o MiddlewareOnion) Append(middlewares ...Middleware) []Middleware {
	newLayers := make([]Middleware, 0, len(o.layers)+len(middlewares))
	newLayers = append(newLayers, o.layers...)
	newLayers = append(newLayers, middlewares...)
	return newLayers
}

// Extend append middlewares to MiddlewareOnion.layers
func (o MiddlewareOnion) Extend(middlewares ...Middleware) {
	o.layers = append(o.layers, middlewares...)
}
