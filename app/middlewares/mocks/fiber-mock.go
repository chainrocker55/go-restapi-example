package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func NewMockFiberCtx(method, path string, headers map[string]string, body []byte, params map[string]string) *fiber.Ctx {
	app := fiber.New()
	req := fasthttp.AcquireRequest()
	// resp := fasthttp.AcquireResponse()

	req.Header.SetMethod(method)

	for key, value := range params {
		req.Header.Set(key, value)
	}
	
	req.SetRequestURI(path)
	if body != nil {
		req.SetBody(body)
	}

	// Create a new fasthttp.RequestCtx
	fctx := &fasthttp.RequestCtx{
	}
	fctx.Init(req, nil, nil)

	// Set path parameters
	for key, value := range params {
		fctx.URI().SetQueryString(key + "=" + value)
	}

	// Create a new fiber.Ctx
	ctx := app.AcquireCtx(fctx)

	return ctx
}
