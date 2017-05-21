package fasthttpmiddleware

import (
	"bytes"
	"github.com/buaazp/fasthttprouter"
	"github.com/hanjm/zaplog"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestNewNormalOnion(t *testing.T) {
	exampleAuthFunc := func(ctx *fasthttp.RequestCtx) bool {
		if bytes.HasPrefix(ctx.Path(), []byte("/protect")) {
			return false
		} else {
			return true
		}
	}
	logger := zaplog.NewNoCallerLogger(false)
	mo := NewNormalMiddlewareOnion(exampleAuthFunc, logger)
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("hello")
	}
	panicHandler := func(ctx *fasthttp.RequestCtx) {
		panic("test panic")
	}
	router := fasthttprouter.New()
	router.GET("/", mo.Apply(requestHandler))
	router.GET("/protect", mo.Apply(requestHandler))
	router.GET("/panic", mo.Apply(panicHandler))
	go func() {
		fasthttp.ListenAndServe(":8000", router.Handler)
	}()
	time.Sleep(time.Second)
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
	time.Sleep(time.Second)
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
