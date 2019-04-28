package lint_test

import (
	"bytes"
	"testing"

	"go-common/app/admin/main/config/pkg/lint"
	_ "go-common/app/admin/main/config/pkg/lint/json"
	_ "go-common/app/admin/main/config/pkg/lint/toml"
)

func TestLint(t *testing.T) {
	jsonRead := bytes.NewBufferString(`{"hello": "world", "a1":"ab"}`)
	err := lint.Lint("json", jsonRead)
	if err != nil {
		t.Errorf("%v", err)
	}
	tomlRead := bytes.NewBufferString(`[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true`)
	err = lint.Lint("toml", tomlRead)
	if err != nil {
		t.Errorf("%v", err)
	}
	err = lint.Lint("test", tomlRead)
	if err != lint.ErrLintNotExists {
		t.Errorf("%v", err)
	}
}
