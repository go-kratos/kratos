package jsonlint

import (
	"encoding/json"
	"io"

	"go-common/app/admin/main/config/pkg/lint"
)

const filetype = "json"

type jsonlint struct{}

func (jsonlint) Lint(r io.Reader) lint.Error {
	var v interface{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&v)
	if err == nil {
		return nil
	}
	return lint.Error{{Line: -1, Message: err.Error()}}
}

func init() {
	lint.RegisterLinter(filetype, jsonlint{})
}
