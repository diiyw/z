package json_test

import (
	gojson "encoding/json"
	"testing"

	"github.com/diiyw/z"
	"github.com/diiyw/z/require"
	"github.com/diiyw/z/stdlib/json"
)

type ARR = []any
type MAP = map[string]any

func TestJSON(t *testing.T) {
	testJSONEncodeDecode(t, nil)

	testJSONEncodeDecode(t, 0)
	testJSONEncodeDecode(t, 1)
	testJSONEncodeDecode(t, -1)
	testJSONEncodeDecode(t, 1984)
	testJSONEncodeDecode(t, -1984)

	testJSONEncodeDecode(t, 0.0)
	testJSONEncodeDecode(t, 1.0)
	testJSONEncodeDecode(t, -1.0)
	testJSONEncodeDecode(t, 19.84)
	testJSONEncodeDecode(t, -19.84)

	testJSONEncodeDecode(t, "")
	testJSONEncodeDecode(t, "foo")
	testJSONEncodeDecode(t, "foo bar")
	testJSONEncodeDecode(t, "foo \"bar\"")
	// See: https://github.com/diiyw/z/issues/268
	testJSONEncodeDecode(t, "1\u001C04")
	testJSONEncodeDecode(t, "çığöşü")
	testJSONEncodeDecode(t, "ç1\u001C04IĞÖŞÜ")
	testJSONEncodeDecode(t, "错误测试")

	testJSONEncodeDecode(t, true)
	testJSONEncodeDecode(t, false)

	testJSONEncodeDecode(t, ARR{})
	testJSONEncodeDecode(t, ARR{0})
	testJSONEncodeDecode(t, ARR{false})
	testJSONEncodeDecode(t, ARR{1, 2, 3,
		"four", false})
	testJSONEncodeDecode(t, ARR{1, 2, 3,
		"four", false, MAP{"a": 0, "b": "bee", "bool": true}})

	testJSONEncodeDecode(t, MAP{})
	testJSONEncodeDecode(t, MAP{"a": 0})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee"})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee", "bool": true})

	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee",
		"arr": ARR{1, 2, 3, "four"}})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee",
		"arr": ARR{1, 2, 3, MAP{"a": false, "b": 109.4}}})

	testJSONEncodeDecode(t, MAP{"id1": 7075984636689534001, "id2": 7075984636689534002})
	testJSONEncodeDecode(t, ARR{1e3, 1e7})
}

func TestDecode(t *testing.T) {
	testDecodeError(t, `{`)
	testDecodeError(t, `}`)
	testDecodeError(t, `{}a`)
	testDecodeError(t, `{{}`)
	testDecodeError(t, `{}}`)
	testDecodeError(t, `[`)
	testDecodeError(t, `]`)
	testDecodeError(t, `[]a`)
	testDecodeError(t, `[[]`)
	testDecodeError(t, `[]]`)
	testDecodeError(t, `"`)
	testDecodeError(t, `"abc`)
	testDecodeError(t, `abc"`)
	testDecodeError(t, `.123`)
	testDecodeError(t, `123.`)
	testDecodeError(t, `1.2.3`)
	testDecodeError(t, `'a'`)
	testDecodeError(t, `true, false`)
	testDecodeError(t, `{"a:"b"}`)
	testDecodeError(t, `{a":"b"}`)
	testDecodeError(t, `{"a":"b":"c"}`)
}

func testDecodeError(t *testing.T, input string) {
	_, err := json.Decode([]byte(input))
	require.Error(t, err)
}

func testJSONEncodeDecode(t *testing.T, v any) {
	o, err := z.FromInterface(v)
	require.NoError(t, err)

	b, err := json.Encode(o)
	require.NoError(t, err)

	a, err := json.Decode(b)
	require.NoError(t, err, string(b))

	vj, err := gojson.Marshal(v)
	require.NoError(t, err)

	aj, err := gojson.Marshal(z.ToInterface(a))
	require.NoError(t, err)

	require.Equal(t, vj, aj)
}
