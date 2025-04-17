package stdlib_test

import (
	"os"
	"testing"

	"github.com/diiyw/z"
	"github.com/diiyw/z/require"
)

func TestReadFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	require.NoError(t, err)
	_ = tf.Close()

	module(t, "os").call("read_file", tf.Name()).
		expect(&z.Bytes{Value: content})
}

func TestReadFileArgs(t *testing.T) {
	module(t, "os").call("read_file").expectError()
}
func TestFileStatArgs(t *testing.T) {
	module(t, "os").call("stat").expectError()
}

func TestFileStatFile(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog")
	tf, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tf.Name()) }()

	_, err = tf.Write(content)
	require.NoError(t, err)
	_ = tf.Close()

	stat, err := os.Stat(tf.Name())
	if err != nil {
		t.Logf("could not get tmp file stat: %s", err)
		return
	}

	module(t, "os").call("stat", tf.Name()).expect(&z.ImmutableMap{
		Value: map[string]z.Object{
			"name":      &z.String{Value: stat.Name()},
			"mtime":     &z.Time{Value: stat.ModTime()},
			"size":      &z.Int{Value: stat.Size()},
			"mode":      &z.Int{Value: int64(stat.Mode())},
			"directory": z.FalseValue,
		},
	})
}

func TestFileStatDir(t *testing.T) {
	td, err := os.MkdirTemp("", "test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(td) }()

	stat, err := os.Stat(td)
	require.NoError(t, err)

	module(t, "os").call("stat", td).expect(&z.ImmutableMap{
		Value: map[string]z.Object{
			"name":      &z.String{Value: stat.Name()},
			"mtime":     &z.Time{Value: stat.ModTime()},
			"size":      &z.Int{Value: stat.Size()},
			"mode":      &z.Int{Value: int64(stat.Mode())},
			"directory": z.TrueValue,
		},
	})
}

func TestOSExpandEnv(t *testing.T) {
	curMaxStringLen := z.MaxStringLen
	defer func() { z.MaxStringLen = curMaxStringLen }()
	z.MaxStringLen = 12

	_ = os.Setenv("Z", "FOO BAR")
	module(t, "os").call("expand_env", "$Z").expect("FOO BAR")

	_ = os.Setenv("Z", "FOO")
	module(t, "os").call("expand_env", "$Z $Z").expect("FOO FOO")

	_ = os.Setenv("Z", "123456789012")
	module(t, "os").call("expand_env", "$Z").expect("123456789012")

	_ = os.Setenv("Z", "1234567890123")
	module(t, "os").call("expand_env", "$Z").expectError()

	_ = os.Setenv("Z", "123456")
	module(t, "os").call("expand_env", "$Z$Z").expect("123456123456")

	_ = os.Setenv("Z", "123456")
	module(t, "os").call("expand_env", "${Z}${Z}").
		expect("123456123456")

	_ = os.Setenv("Z", "123456")
	module(t, "os").call("expand_env", "$Z $Z").expectError()

	_ = os.Setenv("Z", "123456")
	module(t, "os").call("expand_env", "${Z} ${Z}").expectError()
}
