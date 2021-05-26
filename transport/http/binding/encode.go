package binding

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// EncodePath binds proto message to url path
func EncodePath(method string, pathPattern string, msg proto.Message) string {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return pathPattern
	}
	reg := regexp.MustCompile(`/{[.\w]+}`)
	if reg == nil {
		return pathPattern
	}
	pathParams := make(map[string]struct{}, 0)
	path := reg.ReplaceAllStringFunc(pathPattern, func(in string) string {
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
	if method == "GET" {
		u, err := encodeQuery(msg)
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
		if fd = fields.ByName(protoreflect.Name(fieldName)); fd == nil {
			fd = fields.ByJSONName(fieldName)
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

func encodeQuery(msg proto.Message) (url.Values, error) {
	if msg == nil || (reflect.ValueOf(msg).Kind() == reflect.Ptr && reflect.ValueOf(msg).IsNil()) {
		return url.Values{}, nil
	}
	u := make(url.Values)
	err := encodeByField(u, "", msg.ProtoReflect())
	if err != nil {
		return nil, err
	}
	return u, nil
}

func encodeByField(u url.Values, path string, v protoreflect.Message) error {
	for i := 0; i < v.Descriptor().Fields().Len(); i++ {
		fd := v.Descriptor().Fields().Get(i)
		var key string = fd.TextName()
		var newPath string
		if path == "" {
			newPath = key
		} else {
			newPath = path + "." + key
		}
		if fd.HasJSONName() {
			key = fd.JSONName()
		}
		switch {
		case fd.IsList():
			list, err := encodeRepeatedField(fd, v.Get(fd).List())
			if err != nil {
				return err
			}
			u[newPath] = list
		case fd.IsMap():
			m, err := encodeMapField(fd, v.Get(fd).Map())
			if err != nil {
				return err
			}
			for k, value := range m {
				u[fmt.Sprintf("%s[%s]", newPath, k)] = []string{value}
			}
		case (fd.Kind() == protoreflect.MessageKind) || (fd.Kind() == protoreflect.GroupKind):
			value, err := encodeMessage(fd.Message(), v.Get(fd))
			if err == nil {
				u[newPath] = []string{value}
				return nil
			}
			return encodeByField(u, newPath, v.Get(fd).Message())
		default:
			value, err := encodeField(fd, v.Get(fd))
			if err != nil {
				return err
			}
			u[newPath] = []string{value}
		}
	}
	return nil
}

func encodeMessageField(fieldDescriptor protoreflect.FieldDescriptor, msgValue protoreflect.Message) (string, error) {
	return encodeField(fieldDescriptor, msgValue.Get(fieldDescriptor))
}

func encodeRepeatedField(fieldDescriptor protoreflect.FieldDescriptor, list protoreflect.List) ([]string, error) {
	var values []string
	for i := 0; i < list.Len(); i++ {
		value, err := encodeField(fieldDescriptor, list.Get(i))
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func encodeMapField(fieldDescriptor protoreflect.FieldDescriptor, mp protoreflect.Map) (map[string]string, error) {
	m := make(map[string]string)
	mp.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		key, err := encodeField(fieldDescriptor.MapValue(), k.Value())
		if err != nil {
			return false
		}
		value, err := encodeField(fieldDescriptor.MapValue(), v)
		if err != nil {
			return false
		}
		m[key] = value
		return true
	})

	return m, nil
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
