package jsonlint

import (
	"bytes"
	"testing"
)

var testdata = `
{
	"a1": "a1",
	"b2": 
}
`

var testdataok = `
{
	"hello": "world"
}
`

func TestJsonLint(t *testing.T) {
	lint := jsonlint{}
	r := bytes.NewBufferString(testdata)
	lintErr := lint.Lint(r)
	if lintErr == nil {
		t.Fatalf("expect lintErr != nil")
	}
	t.Logf("%s", lintErr.Error())
}

func TestJsonLintOk(t *testing.T) {
	lint := jsonlint{}
	r := bytes.NewBufferString(testdataok)
	lintErr := lint.Lint(r)
	if lintErr != nil {
		t.Error(lintErr)
	}
}
