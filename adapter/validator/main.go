package validator

import vv "github.com/go-playground/validator/v10"

var v *validator

func init() {
	_v := vv.New()
	_ = _v.RegisterValidation("nickname", nickname)

	v = &validator{v: _v}
}
