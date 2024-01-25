package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
)

type S struct {
	App *fiber.App
	Srv *httptest.Server
	ss  *suite.Suite
}

func (dst *S) SetupSuite(o SuiteOptions) {
	dst.App = New()
	dst.ss = o.Suite

	o.Svc.Bind(dst.App)
}

func (dst *S) TearDownSuite() {
	//defer dst.Srv.Close()
}

func (dst *S) SetupTest() {
	// mute
}

func (dst *S) TearDownTest() {
	// mute
}

func (dst *S) GET(path string) *Req {
	return &Req{
		method:  "GET",
		path:    path,
		params:  map[string]interface{}{},
		headers: http.Header{},
		app:     dst.App,
		ss:      dst.ss,
	}
}

func (dst *S) POST(path string) *Req {
	return &Req{
		method:  "POST",
		path:    path,
		params:  map[string]interface{}{},
		headers: http.Header{},
		app:     dst.App,
		ss:      dst.ss,
	}
}

type SuiteOptions struct {
	Svc   Binder
	Suite *suite.Suite
}
