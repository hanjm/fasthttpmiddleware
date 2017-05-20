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
