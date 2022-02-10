package form

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// EncodeValues encode a message into url values.
func EncodeValues(msg proto.Message) (url.Values, error) {
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
		var key string
		var newPath string
		if fd.HasJSONName() {
			key = fd.JSONName()
		} else {
			key = fd.TextName()
		}
		if path == "" {
			newPath = key
		} else {
			newPath = path + "." + key
		}

		if of := fd.ContainingOneof(); of != nil {
			if f := v.WhichOneof(of); f != nil {
				if f != fd {
					continue
				}
			}
			continue
		}
		switch {
		case fd.IsList():
			if v.Get(fd).List().Len() > 0 {
				list, err := encodeRepeatedField(fd, v.Get(fd).List())
				if err != nil {
					return err
				}
				u[newPath] = list
			}
		case fd.IsMap():
			if v.Get(fd).Map().Len() > 0 {
				m, err := encodeMapField(fd, v.Get(fd).Map())
				if err != nil {
					return err
				}
				for k, value := range m {
					u[fmt.Sprintf("%s[%s]", newPath, k)] = []string{value}
				}
			}
		case (fd.Kind() == protoreflect.MessageKind) || (fd.Kind() == protoreflect.GroupKind):
			value, err := encodeMessage(fd.Message(), v.Get(fd))
			if err == nil {
				u[newPath] = []string{value}
				continue
			}
			err = encodeByField(u, newPath, v.Get(fd).Message())
			if err != nil {
				return err
			}
		default:
			value, err := EncodeField(fd, v.Get(fd))
			if err != nil {
				return err
			}
			u[newPath] = []string{value}
		}
	}

	return nil
}

func encodeRepeatedField(fieldDescriptor protoreflect.FieldDescriptor, list protoreflect.List) ([]string, error) {
	var values []string
	for i := 0; i < list.Len(); i++ {
		value, err := EncodeField(fieldDescriptor, list.Get(i))
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
		key, err := EncodeField(fieldDescriptor.MapValue(), k.Value())
		if err != nil {
			return false
		}
		value, err := EncodeField(fieldDescriptor.MapValue(), v)
		if err != nil {
			return false
		}
		m[key] = value
		return true
	})

	return m, nil
}

// EncodeField encode proto message filed
func EncodeField(fieldDescriptor protoreflect.FieldDescriptor, value protoreflect.Value) (string, error) {
	switch fieldDescriptor.Kind() {
	case protoreflect.BoolKind:
		return strconv.FormatBool(value.Bool()), nil
	case protoreflect.EnumKind:
		if fieldDescriptor.Enum().FullName() == "google.protobuf.NullValue" {
			return nullStr, nil
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

// encodeMessage marshals the fields in the given protoreflect.Message.
// If the typeURL is non-empty, then a synthetic "@type" field is injected
// containing the URL as the value.
func encodeMessage(msgDescriptor protoreflect.MessageDescriptor, value protoreflect.Value) (string, error) {
	switch msgDescriptor.FullName() {
	case timestampMessageFullname:
		return marshalTimestamp(value.Message())
	case durationMessageFullname:
		return marshalDuration(value.Message())
	case bytesMessageFullname:
		return marshalBytes(value.Message())
	case "google.protobuf.DoubleValue", "google.protobuf.FloatValue", "google.protobuf.Int64Value", "google.protobuf.Int32Value",
		"google.protobuf.UInt64Value", "google.protobuf.UInt32Value", "google.protobuf.BoolValue", "google.protobuf.StringValue":
		fd := msgDescriptor.Fields()
		v := value.Message().Get(fd.ByName(protoreflect.Name("value")))
		return fmt.Sprintf("%v", v.Interface()), nil
	case "google.protobuf.FieldMask":
		m, ok := value.Message().Interface().(*field_mask.FieldMask)
		if !ok {
			return "", nil
		}
		for i, v := range m.Paths {
			m.Paths[i] = jsonCamelCase(v)
		}
		return strings.Join(m.Paths, ","), nil
	default:
		return "", fmt.Errorf("unsupported message type: %q", string(msgDescriptor.FullName()))
	}
}

// JSONCamelCase converts a snake_case identifier to a camelCase identifier,
// according to the protobuf JSON specification.
// references: https://github.com/protocolbuffers/protobuf-go/blob/master/encoding/protojson/well_known_types.go#L842
func jsonCamelCase(s string) string {
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if c != '_' {
			if wasUnderscore && isASCIILower(c) {
				c -= 'a' - 'A' // convert to uppercase
			}
			b = append(b, c)
		}
		wasUnderscore = c == '_'
	}
	return string(b)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}
