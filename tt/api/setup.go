package api

import "7wd.io/adapter/token"

var tokenf token.F

func init() {
	tokenf = token.New("secret")
}
