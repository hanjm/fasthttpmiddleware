## A funny middleware onion for [fasthttp](github.com/valyala/fasthttp), inspired by [Alice](https://github.com/justinas/alice)
[![GoDoc](https://godoc.org/github.com/hanjm/fasthttpmiddleware?status.svg)](https://godoc.org/github.com/hanjm/fasthttpmiddleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/hanjm/fasthttpmiddleware)](https://goreportcard.com/report/github.com/hanjm/fasthttpmiddleware)
[![code-coverage](http://gocover.io/_badge/github.com/hanjm/fasthttpmiddleware)](http://gocover.io/github.com/hanjm/fasthttpmiddleware)

### Example

```go
package main

import (
	"bytes"
	"github.com/buaazp/fasthttprouter"
	"github.com/hanjm/fasthttpmiddleware"
	"github.com/hanjm/zaplog"
	"github.com/valyala/fasthttp"
)

func exampleAuthFunc(ctx *fasthttp.RequestCtx) bool {
	if bytes.HasPrefix(ctx.Path(), []byte("/protect")) {
		return false
	}
	return true
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("hello")
}

func panicHandler(ctx *fasthttp.RequestCtx) {
	panic("test panic")
}

func main() {
	logger := zaplog.NewNoCallerLogger(false)
	mo := fasthttpmiddleware.NewNormalMiddlewareOnion(exampleAuthFunc, logger)
	moWithoutAuth := fasthttpmiddleware.NewMiddlewareOnion(
		fasthttpmiddleware.NewLogMiddleware(logger, false),
		fasthttpmiddleware.NewRecoverMiddleware(logger),
	)
	router := fasthttprouter.New()
	router.GET("/", mo.Apply(requestHandler))
	router.GET("/protect", mo.Apply(requestHandler))
	router.GET("/panic", mo.Apply(panicHandler))
	router.GET("/noAuth", moWithoutAuth.Apply(requestHandler))
	fasthttp.ListenAndServe(":8000", router.Handler)
}
```

### Document

```go
type AuthFunc func(ctx *fasthttp.RequestCtx) bool
    AuthFunc is your custom auth function type

type Middleware func(h fasthttp.RequestHandler) fasthttp.RequestHandler
    Middleware is a function which receive a fasthttp.RequestHandler then
    return a fasthttp.RequestHandler.

func NewAuthMiddleware(authFunc AuthFunc) Middleware
    NewAuthMiddleware accepts a customer auth function and then returns a
    middleware which only accepts auth passed request. If auth function
    returns false, it will term the HTTP request and response 403 status
    code

func NewLogMiddleware(logger *zap.Logger, xRealIp bool) Middleware
    NewLogMiddleware returns a middleware which log code(status code),
    time(response time), method(request method), path(request URL ath),
    addr(remote address). if the status code is 2xx, the log level is info,
    otherwise, the log level is warn. if your app is behind of Nginx, you
    may meed to set xRealIp to True so that get an actual remote address.

func NewPrometheusMiddleware(bindAddr string, xRealIp bool, logger *zap.Logger) Middleware
    NewPrometheusMiddleware return a middleware which can be used by
    [prometheus](https://github.com/prometheus/prometheus) The prometheus is
    a monitoring system and time series database. Note: the returned
    middleware is contain the function of logmiddleware.

func NewRecoverMiddleware(logger *zap.Logger) Middleware
    NewRecoverMiddleware return a middleware which can let app recover from
    a panic in request handler. panic stack info will appear on the field
    named "stacktrace" in the log line

type MiddlewareOnion struct {
    // contains filtered or unexported fields
}
    MiddlewareOnion represent the middleware like an onion, the bigger index
    of middleware in MiddlewareOnion.layers locate at outside

func NewMiddlewareOnion(middlewares ...Middleware) MiddlewareOnion
    NewMiddlewareOnion returns a middleware onion with given middlewares

func NewNormalMiddlewareOnion(authFunc AuthFunc, xRealIp bool, logger *zap.Logger) MiddlewareOnion
    NewNormalMiddlewareOnion returns a normal middleware onion. recover ->
    auth -> log. the type of AuthFunc is "func(ctx *fasthttp.RequestCtx)
    bool". if your app is behind of Nginx, you may meed to set xRealIp to
    True so that get an actual remote address.

func (o MiddlewareOnion) Append(middlewares ...Middleware) MiddlewareOnion
    Append copy all middleware layers to newLayers, then append middlewares
    to newLayers, then return a new middleware onion.

func (o MiddlewareOnion) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler
    Apply apply the middleware onion to a fasthttp.RequestHandler



	

```

