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

// BindVars bind map parameters to target.
func BindVars(vars map[string]string, target interface{}) error {
	values := make(map[string][]string, len(vars))
	for k, v := range vars {
		values[k] = []string{v}
	}
	if msg, ok := target.(proto.Message); ok {
		return mapProto(msg, values)
	}
	return mapForm(target, values)
}

// BindValue bind map parameters to target.
// Deprecated: use BindVars instead.
func BindValue(vars map[string]string, target interface{}) error {
	return BindVars(vars, target)
}
