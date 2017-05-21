package fasthttpmiddleware

import (
	"bytes"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/hanjm/zaplog"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestNewNormalOnion(t *testing.T) {
	exampleAuthFunc := func(ctx *fasthttp.RequestCtx) bool {
		if bytes.HasPrefix(ctx.Path(), []byte("/protect")) {
			return false
		}
		return true
	}
	logger := zaplog.NewNoCallerLogger(false)
	mo := NewNormalMiddlewareOnion(exampleAuthFunc, false, logger)
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("hello")
	}
	panicHandler := func(ctx *fasthttp.RequestCtx) {
		panic("test panic")
	}
	router := fasthttprouter.New()
	router.GET("/", requestHandler)
	router.GET("/protect", requestHandler)
	router.GET("/panic", panicHandler)
	doneChan := make(chan struct{})
	go func() {
		fasthttp.ListenAndServe(":8000", mo.Apply(router.Handler))
	}()
	go func() {
		var resp []byte
		code, _, _ := fasthttp.Get(resp, "http://127.0.0.1:8000/")
		if code != 200 {
			t.Fatal("unexpected response")
		}
		code, _, _ = fasthttp.Get(resp, "http://127.0.0.1:8000/protect")
		if code != 403 {
			t.Fatal("unexpected response")
		}
		code, _, _ = fasthttp.Get(resp, "http://127.0.0.1:8000/panic")
		if code != 500 {
			t.Fatal("unexpected response")
		}
		doneChan <- struct{}{}
	}()
	<-doneChan
}

func TestNewPrometheusMiddleware(t *testing.T) {
	logger := zaplog.NewNoCallerLogger(false)
	mo := NewMiddlewareOnion(
		NewPrometheusMiddleware(":8001", false, logger),
		NewRecoverMiddleware(logger),
	)
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("hello")
	}
	router := fasthttprouter.New()
	router.GET("/", requestHandler)
	doneChan := make(chan struct{})
	go func() {
		fasthttp.ListenAndServe(":8002", mo.Apply(router.Handler))
	}()
	go func() {
		var resp []byte
		for i := 0; i < 10; i++ {
			fasthttp.Get(resp, "http://127.0.0.1:8002/")
		}
		fasthttp.Get(resp, "http://127.0.0.1:8001/metrics")
		fmt.Printf("%s", resp)
		doneChan <- struct{}{}
	}()
	<-doneChan
}

func TestMiddlewareOnion_Append(t *testing.T) {
	mo := NewMiddlewareOnion()
	if len(mo.layers) != 0 {
		t.Fatal(mo.layers)
	}
	logger := zaplog.NewNoCallerLogger(false)
	loggerMiddleware := NewLogMiddleware(logger, false)
	newMo := mo.Append(loggerMiddleware)
	if len(mo.layers) != 0 || len(newMo.layers) != 1 {
		t.Fatal(mo.layers, newMo.layers)
	}
}
