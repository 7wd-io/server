package http

import (
	"7wd.io/rr"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
)

type Req struct {
	method  string
	path    string
	params  map[string]interface{}
	headers http.Header
	app     *fiber.App
	asserts []func(res *http.Response)
	ss      *suite.Suite
}

func (dst *Req) WithParam(key string, value interface{}) *Req {
	dst.params[key] = value

	return dst
}

func (dst *Req) WithParams(p map[string]interface{}) *Req {
	for k, v := range p {
		dst.params[k] = v
	}

	return dst
}

func (dst *Req) WithToken(t string) *Req {
	dst.headers.Set("Authorization", "Bearer "+t)

	return dst
}

//func (dst *Req) WithAutoPassport() *Req {
//	p := &domain.Passport{
//		Id:       1,
//		Nickname: "autoUser",
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
//			Subject:   "autoUser",
//		},
//	}
//
//	return dst.WithPassport(p)
//}
//
//func (dst *Req) WithPassport(p *domain.Passport) *Req {
//	token, _ := tokenf.Token(p)
//
//	dst.headers.Set("Authorization", "Bearer "+token)
//
//	return dst
//}

func (dst *Req) WithAssertErr(expected error) *Req {
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

func (dst *Req) WithAssertStatus(status int) *Req {
	dst.asserts = append(dst.asserts, func(res *http.Response) {
		dst.ss.Equal(status, res.StatusCode)
	})

	return dst
}

func (dst *Req) WithAssertStatusOk() *Req {
	dst.WithAssertStatus(http.StatusOK)

	return dst
}

func (dst *Req) WithAssertStatusCreated() *Req {
	dst.WithAssertStatus(http.StatusCreated)

	return dst
}

// Send always last called. All asserts before
func (dst *Req) Send() {
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

func (dst *Req) toHttpReq() *http.Request {
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
