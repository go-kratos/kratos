package http

import (
	"net/http"
	"net/url"

	"github.com/go-kratos/kratos/v3/encoding"
	"github.com/go-kratos/kratos/v3/encoding/form"
	"github.com/go-kratos/kratos/v3/errors"
)

func bindQuery(vars url.Values, target any) error {
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(vars.Encode()), target); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}

func bindForm(req *http.Request, target any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := encoding.GetCodec(form.Name).Unmarshal([]byte(req.Form.Encode()), target); err != nil {
		return errors.BadRequest("CODEC", err.Error())
	}
	return nil
}
