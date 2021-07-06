package binding

import (
	"net/http"
	"net/url"

	"github.com/go-kratos/kratos/v2/encoding/form"

	gform "github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

var decoder = gform.NewDecoder()

func init() {
	decoder.SetTagName("json")
}

// BindQuery bind vars parameters to target.
func BindQuery(vars url.Values, target interface{}) error {
	if msg, ok := target.(proto.Message); ok {
		return form.MapProto(msg, vars)
	}

	return decoder.Decode(target, vars)
}

// BindForm bind form parameters to target.
func BindForm(req *http.Request, target interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if msg, ok := target.(proto.Message); ok {
		return form.MapProto(msg, req.Form)
	}
	return decoder.Decode(target, req.Form)
}
