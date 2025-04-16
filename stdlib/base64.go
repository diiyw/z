package stdlib

import (
	"encoding/base64"

	"github.com/diiyw/z"
)

var base64Module = map[string]z.Object{
	"encode": &z.UserFunction{
		Value: FuncAYRS(base64.StdEncoding.EncodeToString),
	},
	"decode": &z.UserFunction{
		Value: FuncASRYE(base64.StdEncoding.DecodeString),
	},
	"raw_encode": &z.UserFunction{
		Value: FuncAYRS(base64.RawStdEncoding.EncodeToString),
	},
	"raw_decode": &z.UserFunction{
		Value: FuncASRYE(base64.RawStdEncoding.DecodeString),
	},
	"url_encode": &z.UserFunction{
		Value: FuncAYRS(base64.URLEncoding.EncodeToString),
	},
	"url_decode": &z.UserFunction{
		Value: FuncASRYE(base64.URLEncoding.DecodeString),
	},
	"raw_url_encode": &z.UserFunction{
		Value: FuncAYRS(base64.RawURLEncoding.EncodeToString),
	},
	"raw_url_decode": &z.UserFunction{
		Value: FuncASRYE(base64.RawURLEncoding.DecodeString),
	},
}
