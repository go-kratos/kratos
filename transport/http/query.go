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

// PopulateVars sets a value in a nested Protobuf structure.
func PopulateVars(msg proto.Message, req *http.Request, include, exclude []string) error {
	for key, value := range Vars(req) {
		if len(include) > 0 {
			if populateFilters(key, include) {
				continue
			}
		}
		if len(exclude) > 0 {
			if !populateFilters(key, exclude) {
				continue
			}
		}
		if err := populateFieldValueFromPath(msg.ProtoReflect(), strings.Split(key, "."), []string{value}); err != nil {
			return err
		}
	}
	return nil
}

// PopulateForm parses query parameters
func PopulateForm(msg proto.Message, req *http.Request, include, exclude []string) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	for key, values := range req.Form {
		if len(include) > 0 {
			if populateFilters(key, exclude) {
				continue
			}
		}
		if len(exclude) > 0 {
			if !populateFilters(key, exclude) {
				continue
			}
		}
		if err := populateFieldValueFromPath(msg.ProtoReflect(), strings.Split(key, "."), values); err != nil {
			return err
		}
	}
	return nil
}

func populateFieldValueFromPath(msgValue protoreflect.Message, fieldPath []string, values []string) error {
	if len(fieldPath) < 1 {
		return errors.New("no field path")
	}
	if len(values) < 1 {
		return errors.New("no value provided")
	}

	var fieldDescriptor protoreflect.FieldDescriptor
	for i, fieldName := range fieldPath {
		fields := msgValue.Descriptor().Fields()

		// Get field by name
		fieldDescriptor = fields.ByName(protoreflect.Name(fieldName))
		if fieldDescriptor == nil {
			fieldDescriptor = fields.ByJSONName(fieldName)
			if fieldDescriptor == nil {
				// We're not returning an error here because this could just be
				// an extra query parameter that isn't part of the request.
				log.Printf("field not found in %q: %q\n", msgValue.Descriptor().FullName(), strings.Join(fieldPath, "."))
				return nil
			}
		}

		// If this is the last element, we're done
		if i == len(fieldPath)-1 {
			break
		}

		// Only singular message fields are allowed
		if fieldDescriptor.Message() == nil || fieldDescriptor.Cardinality() == protoreflect.Repeated {
			return fmt.Errorf("invalid path: %q is not a message", fieldName)
		}

		// Get the nested message
		msgValue = msgValue.Mutable(fieldDescriptor).Message()
	}

	// Check if oneof already set
	if of := fieldDescriptor.ContainingOneof(); of != nil {
		if f := msgValue.WhichOneof(of); f != nil {
			return fmt.Errorf("field already set for oneof %q", of.FullName().Name())
		}
	}

	switch {
	case fieldDescriptor.IsList():
		return populateRepeatedField(fieldDescriptor, msgValue.Mutable(fieldDescriptor).List(), values)
	case fieldDescriptor.IsMap():
		return populateMapField(fieldDescriptor, msgValue.Mutable(fieldDescriptor).Map(), values)
	}

	if len(values) > 1 {
		return fmt.Errorf("too many values for field %q: %s", fieldDescriptor.FullName().Name(), strings.Join(values, ", "))
	}

	return populateField(fieldDescriptor, msgValue, values[0])
}

func populateField(fieldDescriptor protoreflect.FieldDescriptor, msgValue protoreflect.Message, value string) error {
	v, err := parseField(fieldDescriptor, value)
	if err != nil {
		return fmt.Errorf("parsing field %q: %w", fieldDescriptor.FullName().Name(), err)
	}

	msgValue.Set(fieldDescriptor, v)
	return nil
}

func populateRepeatedField(fieldDescriptor protoreflect.FieldDescriptor, list protoreflect.List, values []string) error {
	for _, value := range values {
		v, err := parseField(fieldDescriptor, value)
		if err != nil {
			return fmt.Errorf("parsing list %q: %w", fieldDescriptor.FullName().Name(), err)
		}
		list.Append(v)
	}

	return nil
}

func populateMapField(fieldDescriptor protoreflect.FieldDescriptor, mp protoreflect.Map, values []string) error {
	if len(values) != 2 {
		return fmt.Errorf("more than one value provided for key %q in map %q", values[0], fieldDescriptor.FullName())
	}

	key, err := parseField(fieldDescriptor.MapKey(), values[0])
	if err != nil {
		return fmt.Errorf("parsing map key %q: %w", fieldDescriptor.FullName().Name(), err)
	}

	value, err := parseField(fieldDescriptor.MapValue(), values[1])
	if err != nil {
		return fmt.Errorf("parsing map value %q: %w", fieldDescriptor.FullName().Name(), err)
	}

	mp.Set(key.MapKey(), value)

	return nil
}

func parseField(fieldDescriptor protoreflect.FieldDescriptor, value string) (protoreflect.Value, error) {
	switch fieldDescriptor.Kind() {
	case protoreflect.BoolKind:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return protoreflect.Value{}, err
		}
		return protoreflect.ValueOfBool(v), nil
	case protoreflect.EnumKind:
		enum, err := protoregistry.GlobalTypes.FindEnumByName(fieldDescriptor.Enum().FullName())
		switch {
		case errors.Is(err, protoregistry.NotFound):
			return protoreflect.Value{}, fmt.Errorf("enum %q is not registered", fieldDescriptor.Enum().FullName())
		case err != nil:
			return protoreflect.Value{}, fmt.Errorf("failed to look up enum: %w", err)
		}
		// Look for enum by name
		v := enum.Descriptor().Values().ByName(protoreflect.Name(value))
		if v == nil {
			i, err := strconv.Atoi(value)
			if err != nil {
				return protoreflect.Value{}, fmt.Errorf("%q is not a valid value", value)
			}
			// Look for enum by number
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
		return parseMessage(fieldDescriptor.Message(), value)
	default:
		panic(fmt.Sprintf("unknown field kind: %v", fieldDescriptor.Kind()))
	}
}

func parseMessage(msgDescriptor protoreflect.MessageDescriptor, value string) (protoreflect.Value, error) {
	var msg proto.Message
	switch msgDescriptor.FullName() {
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
		return protoreflect.Value{}, fmt.Errorf("unsupported message type: %q", string(msgDescriptor.FullName()))
	}

	return protoreflect.ValueOfMessage(msg.ProtoReflect()), nil
}

func populateFilters(key string, filters []string) bool {
	for _, s := range filters {
		if s == key {
			return false
		}
	}
	return true
}
