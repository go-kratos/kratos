package binding

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding/form"

	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

var reg = regexp.MustCompile(`/{[\\.\w]+}`)

// EncodeURL encode proto message to url path.
func EncodeURL(pathTemplate string, msg proto.Message, needQuery bool) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}
	pathParams := make(map[string]struct{})
	path := reg.ReplaceAllStringFunc(pathTemplate, func(in string) string {
		// it's unreachable because the reg means that must have more than one char in {}
		//if len(in) < 4 { //nolint:gomnd // **  explain the 4 number here :-) **
		//	return in
		//}
		key := in[2 : len(in)-1]
		vars := strings.Split(key, ".")
		value, err := getValueByField(msg.ProtoReflect(), vars)
		if err != nil {
			return in
		}
		pathParams[key] = struct{}{}
		return "/" + value
	})
	if !needQuery {
		if query := form.EncodeFieldMask(msg.ProtoReflect()); query != "" {
			return path + "?" + query
		}
		return path
	}
	u, err := form.EncodeValues(msg)
	if err == nil && len(u) > 0 {
		for key := range pathParams {
			delete(u, key)
		}
		query := u.Encode()
		if query != "" {
			path += "?" + query
		}
	}
	return path
}

func getValueByField(v protoreflect.Message, fieldPath []string) (string, error) {
	var fd protoreflect.FieldDescriptor
	for i, fieldName := range fieldPath {
		fields := v.Descriptor().Fields()
		if fd = fields.ByJSONName(fieldName); fd == nil {
			fd = fields.ByName(protoreflect.Name(fieldName))
			if fd == nil {
				return "", fmt.Errorf("field path not found: %q", fieldName)
			}
		}
		if i == len(fieldPath)-1 {
			break
		}
		if fd.Message() == nil || fd.Cardinality() == protoreflect.Repeated {
			return "", fmt.Errorf("invalid path: %q is not a message", fieldName)
		}
		v = v.Get(fd).Message()
	}
	return form.EncodeField(fd, v.Get(fd))
}
