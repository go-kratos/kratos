package form

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
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

func encodeByField(u url.Values, path string, m protoreflect.Message) (finalErr error) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		var (
			key     string
			newPath string
		)
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
			if f := m.WhichOneof(of); f != nil {
				if f != fd {
					return true
				}
			}
		}
		switch {
		case fd.IsList():
			if v.List().Len() > 0 {
				list, err := encodeRepeatedField(fd, v.List())
				if err != nil {
					finalErr = err
					return false
				}
				u[newPath] = list
			}
		case fd.IsMap():
			if v.Map().Len() > 0 {
				m, err := encodeMapField(fd, v.Map())
				if err != nil {
					finalErr = err
					return false
				}
				for k, value := range m {
					u[fmt.Sprintf("%s[%s]", newPath, k)] = []string{value}
				}
			}
		case (fd.Kind() == protoreflect.MessageKind) || (fd.Kind() == protoreflect.GroupKind):
			value, err := encodeMessage(fd.Message(), v)
			if err == nil {
				u[newPath] = []string{value}
				return true
			}
			if err = encodeByField(u, newPath, v.Message()); err != nil {
				finalErr = err
				return false
			}
		default:
			value, err := EncodeField(fd, v)
			if err != nil {
				finalErr = err
				return false
			}
			u[newPath] = []string{value}
		}
		return true
	})
	return
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
		return fmt.Sprint(value.Interface()), nil
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
		return fmt.Sprint(v.Interface()), nil
	case fieldMaskFullName:
		m, ok := value.Message().Interface().(*fieldmaskpb.FieldMask)
		if !ok || m == nil {
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

// EncodeFieldMask return field mask name=paths
func EncodeFieldMask(m protoreflect.Message) (query string) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.Kind() == protoreflect.MessageKind {
			if msg := fd.Message(); msg.FullName() == fieldMaskFullName {
				value, err := encodeMessage(msg, v)
				if err != nil {
					return false
				}
				if fd.HasJSONName() {
					query = fd.JSONName() + "=" + value
				} else {
					query = fd.TextName() + "=" + value
				}
				return false
			}
		}
		return true
	})
	return
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
