package binding

import (
	"bytes"
	"encoding/xml"
	"net/http"

	"github.com/pkg/errors"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return errors.WithStack(err)
	}
	return validate(obj)
}

func (xmlBinding) BindBody(body []byte, obj interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(obj); err != nil {
		return errors.WithStack(err)
	}
	return validate(obj)
}
