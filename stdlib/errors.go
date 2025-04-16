package stdlib

import (
	"github.com/diiyw/z"
)

func wrapError(err error) z.Object {
	if err == nil {
		return z.TrueValue
	}
	return &z.Error{Value: &z.String{Value: err.Error()}}
}
