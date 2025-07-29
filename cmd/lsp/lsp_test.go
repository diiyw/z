package lsp

import (
	"testing"
)

func TestOnDefinition(t *testing.T) {
	OnDefinition([]byte(`{"code":"fmt := import(\"fmt\")","offset":19}`))
}
