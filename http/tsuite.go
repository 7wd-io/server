package http

import (
	"7wd.io/di"
	"7wd.io/domain"
	"7wd.io/rr"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"time"
)

type tsuite struct {
	suite.Suite
	baseUrl string
}

func (dst *tsuite) SetupSuite() {
	c := di.MustNew()

	gameSvc := domain.NewGameService(
		c.Clock,
		c.Repo.Room,
		c.Repo.Game,
		c.Repo.GameClock,
		c.Repo.User,
		c.Dispatcher,
	)

	h := NewGame(
		gameSvc,
		domain.NewPlayAgainService(
			c.PlayAgain,
			c.Dispatcher,
			c.Repo.User,
			c.Repo.Room,
			gameSvc,
			c.Repo.Game,
		),
	)

	srv := New()

	h.Bind(srv)
}

func (dst *tsuite) TearDownSuite() {
	// mute
}

func (dst *tsuite) TearDownTest() {
	// mute
}

func (dst *tsuite) GET(path string) *TReq {
	return &TReq{
		method:  "GET",
		path:    dst.baseUrl + path,
		params:  map[string]interface{}{},
		headers: http.Header{},
		//app:     dst.App,
		//ss:      dst.ss,
	}
}

func (dst *tsuite) POST(path string) *TReq {
	return &TReq{
		method:  "POST",
		path:    dst.baseUrl + path,
		params:  map[string]interface{}{},
		headers: http.Header{},
		//app:     dst.App,
		//ss:      dst.ss,
	}
}

type TReq struct {
	method  string
	path    string
	params  map[string]interface{}
	headers http.Header
	app     *fiber.App
	asserts []func(res *http.Response)
	ss      *suite.Suite
}

func (dst *TReq) WithParam(key string, value interface{}) *TReq {
	dst.params[key] = value

	return dst
}

func (dst *TReq) WithParams(p map[string]interface{}) *TReq {
	for k, v := range p {
		dst.params[k] = v
	}

	return dst
}

func (dst *TReq) WithAutoPassport() *TReq {
	p := &domain.Passport{
		Id:       1,
		Nickname: "autoUser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Subject:   "autoUser",
		},
	}

	return dst.WithPassport(p)
}

func (dst *TReq) WithPassport(p *domain.Passport) *TReq {
	//token, _ := tokenf.Token(p)
	token := "secret"

	dst.headers.Set("Authorization", "Bearer "+token)

	return dst
}

func (dst *TReq) WithAssertErr(expected error) *TReq {
	dst.asserts = append(dst.asserts, func(res *http.Response) {
		var actual rr.AppError

		if err := json.NewDecoder(res.Body).Decode(&actual); err != nil {
			dst.ss.FailNow("WithAssertErr unmarshal failed")
		}

		dst.ss.Equal(http.StatusBadRequest, res.StatusCode)
		dst.ss.Equal(expected, actual)
	})

	return dst
}

func (dst *TReq) WithAssertStatus(status int) *TReq {
	dst.asserts = append(dst.asserts, func(res *http.Response) {
		dst.ss.Equal(status, res.StatusCode)
	})

	return dst
}

func (dst *TReq) WithAssertStatusOk() *TReq {
	dst.WithAssertStatus(http.StatusOK)

	return dst
}

func (dst *TReq) WithAssertStatusCreated() *TReq {
	dst.WithAssertStatus(http.StatusCreated)

	return dst
}

// Send always last called. All asserts before
func (dst *TReq) Send() {
	//res, err := dst.app.Test(dst.toHttpReq(), -1)
	res, err := dst.app.Test(dst.toHttpReq())

	dst.ss.NoError(err)

	if res == nil {
		dst.ss.FailNow("response nil")
	} else {
		defer res.Body.Close()
	}

	for _, assert := range dst.asserts {
		assert(res)
	}
}

func (dst *TReq) toHttpReq() *http.Request {
	var r *http.Request

	switch dst.method {
	case "GET":
		r = httptest.NewRequest(
			dst.method,
			dst.path,
			nil,
		)

		q := r.URL.Query()

		for k, v := range dst.params {
			q.Add(k, v.(string))
		}

		r.URL.RawQuery = q.Encode()

	default:
		b, _ := json.Marshal(dst.params)

		r = httptest.NewRequest(
			dst.method,
			dst.path,
			bytes.NewBuffer(b),
		)
	}

	//for _, h := range dst.headers {
	//
	//}

	r.Header = dst.headers
	r.Header.Set("Content-Type", "application/json")

	return r
}
