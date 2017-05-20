## a funny middleware onion for [fasthttp](github.com/valyala/fasthttp). inspired by [Alice](https://github.com/justinas/alice)

### Example

```go
package main

import (
	"bytes"
	"fasthttpmiddleware"
	"github.com/buaazp/fasthttprouter"
	"github.com/hanjm/fasthttpmiddleware"
	"github.com/hanjm/zaplog"
	"github.com/valyala/fasthttp"
)

func exampleAuthFunc(ctx *fasthttp.RequestCtx) bool {
	if bytes.HasPrefix(ctx.Path(), []byte("/protect")) {
		return false
	} else {
		return true
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("hello")
}

func panicHandler(ctx *fasthttp.RequestCtx) {
	panic("test panic")
}

func main() {
	logger := zaplog.NewNoCallerLogger(false)
	middleware := fasthttpmiddleware.NewNormalMiddlewareOnion(exampleAuthFunc, logger)
	router := fasthttprouter.New()
	router.GET("/", middleware.Apply(requestHandler))
	router.GET("/protect", middleware.Apply(requestHandler))
	router.GET("/panic", middleware.Apply(panicHandler))
	fasthttp.ListenAndServe(":8000", router.Handler)
}
```

### Document

```go
type AuthFunc func(ctx *fasthttp.RequestCtx) bool

type Middleware func(h fasthttp.RequestHandler) fasthttp.RequestHandler

func NewAuthMiddleware(authFunc AuthFunc) Middleware

func NewLogMiddleware(logger *zap.Logger) Middleware

func NewRecoverMiddleware(logger *zap.Logger) Middleware

type MiddlewareOnion struct {
    // contains filtered or unexported fields
}
    MiddlewareOnion represent the middlewares like a onion, the bigger index
    of middleware in MiddlewareOnion.layers locate at outside

func New(middlewares ...Middleware) MiddlewareOnion
    New return a middleware onion with given middlewares

func NewNormalMiddlewareOnion(authFunc AuthFunc, logger *zap.Logger) MiddlewareOnion
    NewNormalMiddlewareOnion return a normal middleware onion. recover ->
    auth -> log type AuthFunc func(ctx *fasthttp.RequestCtx) bool

func (o MiddlewareOnion) Append(middlewares ...Middleware) []Middleware
    Append will copy all middleware layers to newLayers, then append
    middlewares in to newLayers

func (o MiddlewareOnion) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler

func (o MiddlewareOnion) Extend(middlewares ...Middleware)
    Extend will then append middlewares in to MiddlewareOnion.layers



	

```

