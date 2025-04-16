package stdlib

import (
	"github.com/d5/tengo/v2"
)

func wrapError(err error) z.Object {
	if err == nil {
		return z.TrueValue
	}
	return &z.Error{Value: &z.String{Value: err.Error()}}
}
