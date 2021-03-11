package binding

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

// BindForm bind form parameters to target.
func BindForm(req *http.Request, target interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if msg, ok := target.(proto.Message); ok {
		return mapProto(msg, req.Form)
	}
	return mapForm(target, req.Form)
}
