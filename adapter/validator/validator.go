package validator

import vv "github.com/go-playground/validator/v10"

type validator struct {
	v *vv.Validate
}

func (dst *validator) Validate(s interface{}) error {
	return dst.v.Struct(s)
}
