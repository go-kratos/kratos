package binding

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return errors.WithStack(err)
	}
	return validate(obj)
}

func (jsonBinding) BindBody(body []byte, obj interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
