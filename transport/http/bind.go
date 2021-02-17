package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// BindVars parses url parameters.
func BindVars(req *http.Request, msg proto.Message) error {
	for key, value := range Vars(req) {
		if err := populateFieldValues(msg.ProtoReflect(), strings.Split(key, "."), []string{value}); err != nil {
			return err
		}
	}
	return nil
}

// BindForm parses form parameters.
func BindForm(req *http.Request, msg proto.Message) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	for key, values := range req.Form {
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
		fields := v.Descriptor().Fields()

		if fd = fields.ByName(protoreflect.Name(fieldName)); fd == nil {
			fd = fields.ByJSONName(fieldName)
			if fd == nil {
				log.Printf("field not found in %q: %q\n", v.Descriptor().FullName(), strings.Join(fieldPath, "."))
				return nil
			}
		}

		if i == len(fieldPath)-1 {
			break
		}

		if fd.Message() == nil || fd.Cardinality() == protoreflect.Repeated {
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
		return populateMapField(fd, v.Mutable(fd).Map(), values)
	}
	if len(values) > 1 {
		return fmt.Errorf("too many values for field %q: %s", fd.FullName().Name(), strings.Join(values, ", "))
	}
	return populateField(fd, v, values[0])
}

func populateField(fd protoreflect.FieldDescriptor, v protoreflect.Message, value string) error {
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

func populateMapField(fd protoreflect.FieldDescriptor, mp protoreflect.Map, values []string) error {
	if len(values) != 2 {
		return fmt.Errorf("more than one value provided for key %q in map %q", values[0], fd.FullName())
	}
	key, err := parseField(fd.MapKey(), values[0])
	if err != nil {
		return fmt.Errorf("parsing map key %q: %w", fd.FullName().Name(), err)
	}
	value, err := parseField(fd.MapValue(), values[1])
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
			i, err := strconv.Atoi(value)
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
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfInt64(v), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfUint64(v), nil
	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil
	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(value, 64)
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
		if value == "null" {
			break
		}
		t, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg, err = ptypes.TimestampProto(t)
		if err != nil {
			return protoreflect.Value{}, err
		}
	case "google.protobuf.Duration":
		if value == "null" {
			break
		}
		d, err := time.ParseDuration(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = ptypes.DurationProto(d)
	case "google.protobuf.DoubleValue":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.DoubleValue{Value: v}
	case "google.protobuf.FloatValue":
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.FloatValue{Value: float32(v)}
	case "google.protobuf.Int64Value":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.Int64Value{Value: v}
	case "google.protobuf.Int32Value":
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.Int32Value{Value: int32(v)}
	case "google.protobuf.UInt64Value":
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.UInt64Value{Value: v}
	case "google.protobuf.UInt32Value":
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.UInt32Value{Value: uint32(v)}
	case "google.protobuf.BoolValue":
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.BoolValue{Value: v}
	case "google.protobuf.StringValue":
		msg = &wrappers.StringValue{Value: value}
	case "google.protobuf.BytesValue":
		v, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		msg = &wrappers.BytesValue{Value: v}
	case "google.protobuf.FieldMask":
		fm := &field_mask.FieldMask{}
		fm.Paths = append(fm.Paths, strings.Split(value, ",")...)
		msg = fm
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported message type: %q", string(md.FullName()))
	}
	return protoreflect.ValueOfMessage(msg.ProtoReflect()), nil
}
