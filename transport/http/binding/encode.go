package binding

import (
	"reflect"
	"regexp"

	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v2/encoding/form"
)

var reg = regexp.MustCompile(`{[\\.\w]+}`)

// EncodeURL encode proto message to url path.
func EncodeURL(pathTemplate string, msg interface{}, needQuery bool) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}
	queryParams, _ := form.EncodeValues(msg)
	pathParams := make(map[string]struct{})
	path := reg.ReplaceAllStringFunc(pathTemplate, func(in string) string {
		// it's unreachable because the reg means that must have more than one char in {}
		// if len(in) < 4 { //nolint:mnd // **  explain the 4 number here :-) **
		//	return in
		// }
		key := in[1 : len(in)-1]
		pathParams[key] = struct{}{}
		return queryParams.Get(key)
	})
	if !needQuery {
		if v, ok := msg.(proto.Message); ok {
			if query := form.EncodeFieldMask(v.ProtoReflect()); query != "" {
				return path + "?" + query
			}
		}
		return path
	}
	if len(queryParams) > 0 {
		for key := range pathParams {
			delete(queryParams, key)
		}
		if query := queryParams.Encode(); query != "" {
			path += "?" + query
		}
	}
	return path
}
