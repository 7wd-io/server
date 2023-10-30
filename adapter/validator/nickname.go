package validator

import (
	vv "github.com/go-playground/validator/v10"
	"regexp"
)

const (
	nicknameMinLength = 3
	nicknameMaxLength = 15
)

var reNickname = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9]*$`)

func nickname(fl vv.FieldLevel) bool {
	value := fl.Field().String()
	l := len([]rune(value))

	if l < nicknameMinLength || l > nicknameMaxLength {
		return false
	}

	return reNickname.MatchString(value)
}
