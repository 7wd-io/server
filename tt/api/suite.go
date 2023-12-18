package api

import (
	srv "7wd.io/http"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
)

type Suite struct {
	App     *fiber.App
	Srv     *httptest.Server
	baseUrl string
	ss      *suite.Suite
}

func (dst *Suite) SetupSuite(o SuiteOptions) {
	dst.App = srv.New()
	dst.ss = o.Suite
	dst.baseUrl = o.BaseUrl

	o.Svc.Bind(dst.App)
}

func (dst *Suite) TearDownSuite() {
	//defer dst.Srv.Close()
}

func (dst *Suite) SetupTest() {
	// mute
}

func (dst *Suite) TearDownTest() {
	// mute
}

func (dst *Suite) GET(path string) *Req {
	return &Req{
		method:  "GET",
		path:    dst.baseUrl + path,
		params:  map[string]interface{}{},
		headers: http.Header{},
		app:     dst.App,
		ss:      dst.ss,
	}
}

func (dst *Suite) POST(path string) *Req {
	return &Req{
		method:  "POST",
		path:    dst.baseUrl + path,
		params:  map[string]interface{}{},
		headers: http.Header{},
		app:     dst.App,
		ss:      dst.ss,
	}
}

type SuiteOptions struct {
	Svc     Binder
	BaseUrl string
	Suite   *suite.Suite
}
