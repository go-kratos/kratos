package form

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const fieldSeparator = "."

var errInvalidFormatMapKey = errors.New("invalid formatting for map key")

// DecodeValues decode url value into proto message.
func DecodeValues(msg proto.Message, values url.Values) error {
	for key, values := range values {
		if err := populateFieldValues(msg.ProtoReflect(), strings.Split(key, "."), values); err != nil {
			return err
		}
	}
	return nil
}

func populateFieldValues(v protoreflect.Message, fieldPath []string, values []string) error {
	if len(fieldPath) < 1 {
		return errors.New("no field path")
	}
	if len(values) < 1 {
		return errors.New("no value provided")
	}

	var fd protoreflect.FieldDescriptor
	for i, fieldName := range fieldPath {
		if fd = getFieldDescriptor(v, fieldName); fd == nil {
			// ignore unexpected field.
			return nil
		}
		if fd.IsMap() && len(fieldPath) == 2 {
			return populateMapField(fd, v.Mutable(fd).Map(), fieldPath, values)
		}
		if i == len(fieldPath)-1 {
			break
		}
		if fd.Message() == nil || fd.Cardinality() == protoreflect.Repeated {
			if fd.IsMap() && len(fieldPath) > 1 {
				// post subfield
				return populateMapField(fd, v.Mutable(fd).Map(), []string{fieldPath[1]}, values)
			}
			return fmt.Errorf("invalid path: %q is not a message", fieldName)
		}

		v = v.Mutable(fd).Message()
	}
	if of := fd.ContainingOneof(); of != nil {
		if f := v.WhichOneof(of); f != nil {
			return fmt.Errorf("field already set for oneof %q", of.FullName().Name())
		}
	}
	switch {
	case fd.IsList():
		return populateRepeatedField(fd, v.Mutable(fd).List(), values)
	case fd.IsMap():
		return populateMapField(fd, v.Mutable(fd).Map(), fieldPath, values)
	}
	if len(values) > 1 {
		return fmt.Errorf("too many values for field %q: %s", fd.FullName().Name(), strings.Join(values, ", "))
	}
	return populateField(fd, v, values[0])
}

func getFieldDescriptor(v protoreflect.Message, fieldName string) protoreflect.FieldDescriptor {
	var (
		fields = v.Descriptor().Fields()
		fd     = getDescriptorByFieldAndName(fields, fieldName)
	)
	if fd == nil {
		switch {
		case v.Descriptor().FullName() == structMessageFullname:
			fd = fields.ByNumber(structFieldsFieldNumber)
		case len(fieldName) > 2 && strings.HasSuffix(fieldName, "[]"):
			fd = getDescriptorByFieldAndName(fields, strings.TrimSuffix(fieldName, "[]"))
		default:
			// If the type is map, you get the string "map[kratos]", where "map" is a field of proto and "kratos" is a key of map
			// Use symbol . for separating fields/structs. (eg. structfield.field)
			// ref: https://github.com/go-playground/form
			field, _, err := parseURLQueryMapKey(fieldName)
			if err != nil {
				break
			}
			fd = getDescriptorByFieldAndName(fields, field)
		}
	}
	return fd
}

func getDescriptorByFieldAndName(fields protoreflect.FieldDescriptors, fieldName string) protoreflect.FieldDescriptor {
	var fd protoreflect.FieldDescriptor
	if fd = fields.ByName(protoreflect.Name(fieldName)); fd == nil {
		fd = fields.ByJSONName(fieldName)
	}
	return fd
}

func populateField(fd protoreflect.FieldDescriptor, v protoreflect.Message, value string) error {
	if value == "" {
		return nil
	}
	val, err := parseField(fd, value)
	if err != nil {
		return fmt.Errorf("parsing field %q: %w", fd.FullName().Name(), err)
	}
	v.Set(fd, val)
	return nil
}

func populateRepeatedField(fd protoreflect.FieldDescriptor, list protoreflect.List, values []string) error {
	for _, value := range values {
		v, err := parseField(fd, value)
		if err != nil {
			return fmt.Errorf("parsing list %q: %w", fd.FullName().Name(), err)
		}
		list.Append(v)
	}
	return nil
}

func populateMapField(fd protoreflect.FieldDescriptor, mp protoreflect.Map, fieldPath []string, values []string) error {
	_, keyName, err := parseURLQueryMapKey(strings.Join(fieldPath, fieldSeparator))
	if err != nil {
		return err
	}
	key, err := parseField(fd.MapKey(), keyName)
	if err != nil {
		return fmt.Errorf("parsing map key %q: %w", fd.FullName().Name(), err)
	}
	value, err := parseField(fd.MapValue(), values[len(values)-1])
	if err != nil {
		return fmt.Errorf("parsing map value %q: %w", fd.FullName().Name(), err)
	}
	mp.Set(key.MapKey(), value)
	return nil
}

func parseField(fd protoreflect.FieldDescriptor, value string) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBool(v), nil
	case protoreflect.EnumKind:
		enum, err := protoregistry.GlobalTypes.FindEnumByName(fd.Enum().FullName())
		switch {
		case errors.Is(err, protoregistry.NotFound):
			return protoreflect.Value{}, fmt.Errorf("enum %q is not registered", fd.Enum().FullName())
		case err != nil:
			return protoreflect.Value{}, fmt.Errorf("failed to look up enum: %w", err)
		}
		v := enum.Descriptor().Values().ByName(protoreflect.Name(value))
		if v == nil {
			i, err := strconv.ParseInt(value, 10, 32) //nolint:mnd
			if err != nil {
				return protoreflect.Value{}, fmt.Errorf("%q is not a valid value", value)
			}
			v = enum.Descriptor().Values().ByNumber(protoreflect.EnumNumber(i))
			if v == nil {
				return protoreflect.Value{}, fmt.Errorf("%q is not a valid value", value)
			}
		}
		return protoreflect.ValueOfEnum(v.Number()), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(value, 10, 32) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(value, 10, 64) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt64(v), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(value, 10, 32) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(value, 10, 64) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint64(v), nil
	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(value, 32) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil
	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(value, 64) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat64(v), nil
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(value), nil
	case protoreflect.BytesKind:
		v, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBytes(v), nil
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return parseMessage(fd.Message(), value)
	default:
		panic(fmt.Sprintf("unknown field kind: %v", fd.Kind()))
	}
}

func parseMessage(md protoreflect.MessageDescriptor, value string) (protoreflect.Value, error) {
	var msg proto.Message
	switch md.FullName() {
	case "google.protobuf.Timestamp":
		if value == nullStr {
			break
		}
		t, err := time.ParseInLocation(time.RFC3339Nano, value, time.Local)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = timestamppb.New(t)
	case "google.protobuf.Duration":
		if value == nullStr {
			break
		}
		d, err := time.ParseDuration(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = durationpb.New(d)
	case "google.protobuf.DoubleValue":
		v, err := strconv.ParseFloat(value, 64) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.Double(v)
	case "google.protobuf.FloatValue":
		v, err := strconv.ParseFloat(value, 32) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.Float(float32(v))
	case "google.protobuf.Int64Value":
		v, err := strconv.ParseInt(value, 10, 64) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.Int64(v)
	case "google.protobuf.Int32Value":
		v, err := strconv.ParseInt(value, 10, 32) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.Int32(int32(v))
	case "google.protobuf.UInt64Value":
		v, err := strconv.ParseUint(value, 10, 64) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.UInt64(v)
	case "google.protobuf.UInt32Value":
		v, err := strconv.ParseUint(value, 10, 32) //nolint:mnd
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.UInt32(uint32(v))
	case "google.protobuf.BoolValue":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = wrapperspb.Bool(v)
	case "google.protobuf.StringValue":
		msg = wrapperspb.String(value)
	case "google.protobuf.BytesValue":
		v, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			if v, err = base64.URLEncoding.DecodeString(value); err != nil {
				return protoreflect.Value{}, err
			}
		}
		msg = wrapperspb.Bytes(v)
	case "google.protobuf.FieldMask":
		fm := &fieldmaskpb.FieldMask{}
		for _, fv := range strings.Split(value, ",") {
			fm.Paths = append(fm.Paths, jsonSnakeCase(fv))
		}
		msg = fm
	case "google.protobuf.Value":
		fm, err := structpb.NewValue(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = fm
	case "google.protobuf.Struct":
		var v structpb.Struct
		if err := protojson.Unmarshal([]byte(value), &v); err != nil {
			return protoreflect.Value{}, err
		}
		msg = &v
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported message type: %q", string(md.FullName()))
	}
	return protoreflect.ValueOfMessage(msg.ProtoReflect()), nil
}

// jsonSnakeCase converts a camelCase identifier to a snake_case identifier,
// according to the protobuf JSON specification.
// references: https://github.com/protocolbuffers/protobuf-go/blob/master/encoding/protojson/well_known_types.go#L864
func jsonSnakeCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if isASCIIUpper(c) {
			b = append(b, '_')
			c += 'a' - 'A' // convert to lowercase
		}
		b = append(b, c)
	}
	return string(b)
}

func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// parseURLQueryMapKey parse the url.Values the field name and key name of the value map type key
// for example: convert "map[key]" to "map" and "key"
func parseURLQueryMapKey(key string) (string, string, error) {
	var (
		startIndex = strings.IndexByte(key, '[')
		endIndex   = strings.IndexByte(key, ']')
	)
	if startIndex < 0 {
		//nolint:mnd
		values := strings.Split(key, fieldSeparator)
		//nolint:mnd
		if len(values) != 2 {
			return "", "", errInvalidFormatMapKey
		}
		return values[0], values[1], nil
	}
	if startIndex <= 0 || startIndex >= endIndex || len(key) != endIndex+1 {
		return "", "", errInvalidFormatMapKey
	}
	return key[:startIndex], key[startIndex+1 : endIndex], nil
}
