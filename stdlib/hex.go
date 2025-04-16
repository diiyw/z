package stdlib

import (
	"encoding/hex"

	"github.com/d5/tengo/v2"
)

var hexModule = map[string]z.Object{
	"encode": &z.UserFunction{Value: FuncAYRS(hex.EncodeToString)},
	"decode": &z.UserFunction{Value: FuncASRYE(hex.DecodeString)},
}
