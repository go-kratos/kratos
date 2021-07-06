package binding

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/encoding/form"

	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// EncodeURL encode proto message to url path.
func EncodeURL(pathTemplate string, msg proto.Message, needQuery bool) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathTemplate
	}
	reg := regexp.MustCompile(`/{[.\w]+}`)
	if reg == nil {
		return pathTemplate
	}
	pathParams := make(map[string]struct{}, 0)
	path := reg.ReplaceAllStringFunc(pathTemplate, func(in string) string {
		if len(in) < 4 {
			return in
		}
		key := in[2 : len(in)-1]
		vars := strings.Split(key, ".")
		if value, err := getValueByField(msg.ProtoReflect(), vars); err == nil {
			pathParams[key] = struct{}{}
			return "/" + value
		}
		return in
	})
	if needQuery {
		u, err := form.EncodeMap(msg)
		if err == nil && len(u) > 0 {
			for key := range pathParams {
				delete(u, key)
			}
			query := u.Encode()
			if query != "" {
				path += "?" + query
			}
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
	return encodeField(fd, v.Get(fd))
}

func encodeField(fieldDescriptor protoreflect.FieldDescriptor, value protoreflect.Value) (string, error) {
	switch fieldDescriptor.Kind() {
	case protoreflect.BoolKind:
		return strconv.FormatBool(value.Bool()), nil
	case protoreflect.EnumKind:
		if fieldDescriptor.Enum().FullName() == "google.protobuf.NullValue" {
			return "null", nil
		}
		desc := fieldDescriptor.Enum().Values().ByNumber(value.Enum())
		return string(desc.Name()), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.BytesKind:
		return base64.URLEncoding.EncodeToString(value.Bytes()), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return encodeMessage(fieldDescriptor.Message(), value)
	default:
		return fmt.Sprintf("%v", value.Interface()), nil
	}
}

// marshalMessage marshals the fields in the given protoreflect.Message.
// If the typeURL is non-empty, then a synthetic "@type" field is injected
// containing the URL as the value.
func encodeMessage(msgDescriptor protoreflect.MessageDescriptor, value protoreflect.Value) (string, error) {
	switch msgDescriptor.FullName() {
	case "google.protobuf.Timestamp":
		t, ok := value.Interface().(*timestamppb.Timestamp)
		if !ok {
			return "", nil
		}
		return t.AsTime().Format(time.RFC3339Nano), nil
	case "google.protobuf.Duration":
		d, ok := value.Interface().(*durationpb.Duration)
		if !ok {
			return "", nil
		}
		return d.AsDuration().String(), nil
	case "google.protobuf.BytesValue":
		b, ok := value.Interface().(*wrapperspb.BytesValue)
		if !ok {
			return "", nil
		}
		return base64.StdEncoding.EncodeToString(b.Value), nil
	case "google.protobuf.DoubleValue", "google.protobuf.FloatValue", "google.protobuf.Int64Value", "google.protobuf.Int32Value",
		"google.protobuf.UInt64Value", "google.protobuf.UInt32Value", "google.protobuf.BoolValue", "google.protobuf.StringValue":
		fd := msgDescriptor.Fields()
		v := value.Message().Get(fd.ByName(protoreflect.Name("value"))).Message()
		return fmt.Sprintf("%v", v.Interface()), nil
	case "google.protobuf.FieldMask":
		m, ok := value.Interface().(*field_mask.FieldMask)
		if !ok {
			return "", nil
		}
		return strings.Join(m.Paths, ","), nil
	default:
		return "", fmt.Errorf("unsupported message type: %q", string(msgDescriptor.FullName()))
	}
}
