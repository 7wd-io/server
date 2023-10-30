package token

import (
	"7wd.io/domain"
)

func New(secret string) F {
	return F{secret: secret}
}

type F struct {
	secret string
}

func (dst F) Token(passport *domain.Passport) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, passport).SignedString([]byte(dst.secret))
}
