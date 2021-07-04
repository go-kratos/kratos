package form

import (
	"encoding/base64"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/schema"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Name is the name registered for the form codec.
const Name = "x-www-form-urlencoded"

var invalidValue = reflect.Value{}

func init() {
	// encoder
	encoder := schema.NewEncoder()
	encoder.SetAliasTag("json")
	encoder.RegisterEncoder(&timestamppb.Timestamp{}, func(v reflect.Value) string {
		r := v.Interface().(*timestamppb.Timestamp)
		return r.AsTime().Format(time.RFC3339Nano)
	})
	encoder.RegisterEncoder(&durationpb.Duration{}, func(v reflect.Value) string {
		r := v.Interface().(*durationpb.Duration)
		return r.AsDuration().String()
	})
	encoder.RegisterEncoder(&wrapperspb.BytesValue{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.BytesValue)
		return base64.StdEncoding.EncodeToString(r.Value)
	})
	encoder.RegisterEncoder(&wrapperspb.DoubleValue{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.DoubleValue)
		return strconv.FormatFloat(r.Value, 'E', -1, 64)
	})
	encoder.RegisterEncoder(&wrapperspb.FloatValue{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.FloatValue)
		return strconv.FormatFloat(float64(r.Value), 'E', -1, 32)
	})
	encoder.RegisterEncoder(&wrapperspb.Int64Value{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.Int64Value)
		return strconv.FormatInt(r.Value, 10)
	})
	encoder.RegisterEncoder(&wrapperspb.Int32Value{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.Int32Value)
		return strconv.FormatInt(int64(r.Value), 10)
	})
	encoder.RegisterEncoder(&wrapperspb.UInt64Value{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.UInt64Value)
		return strconv.FormatUint(r.Value, 10)
	})
	encoder.RegisterEncoder(&wrapperspb.UInt32Value{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.UInt32Value)
		return strconv.FormatUint(uint64(r.Value), 10)
	})
	encoder.RegisterEncoder(&wrapperspb.BoolValue{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.BoolValue)
		return strconv.FormatBool(r.Value)
	})
	encoder.RegisterEncoder(&wrapperspb.StringValue{}, func(v reflect.Value) string {
		r := v.Interface().(*wrapperspb.StringValue)
		return r.Value
	})
	// decoder
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.RegisterConverter(&timestamppb.Timestamp{}, func(v string) reflect.Value {
		r, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(timestamppb.New(r))
	})
	decoder.RegisterConverter(&durationpb.Duration{}, func(v string) reflect.Value {
		r, err := time.ParseDuration(v)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(durationpb.New(r))
	})
	decoder.RegisterConverter(&wrapperspb.DoubleValue{}, func(v string) reflect.Value {
		r, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.Double(r))
	})
	decoder.RegisterConverter(&wrapperspb.FloatValue{}, func(v string) reflect.Value {
		r, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.Float(float32(r)))
	})
	decoder.RegisterConverter(&wrapperspb.Int64Value{}, func(v string) reflect.Value {
		r, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.Int32(int32(r)))
	})
	decoder.RegisterConverter(&wrapperspb.UInt64Value{}, func(v string) reflect.Value {
		r, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.UInt64(r))
	})
	decoder.RegisterConverter(&wrapperspb.UInt32Value{}, func(v string) reflect.Value {
		r, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.UInt32(uint32(r)))
	})
	decoder.RegisterConverter(&wrapperspb.BoolValue{}, func(v string) reflect.Value {
		r, err := strconv.ParseBool(v)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.Bool(r))
	})
	decoder.RegisterConverter(&wrapperspb.StringValue{}, func(v string) reflect.Value {
		return reflect.ValueOf(wrapperspb.String(v))
	})
	decoder.RegisterConverter(&wrapperspb.StringValue{}, func(v string) reflect.Value {
		r, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return invalidValue
		}
		return reflect.ValueOf(wrapperspb.Bytes(r))
	})
	decoder.RegisterConverter(&field_mask.FieldMask{}, func(v string) reflect.Value {
		fm := &field_mask.FieldMask{}
		fm.Paths = append(fm.Paths, strings.Split(v, ",")...)
		return reflect.ValueOf(fm)
	})
	encoding.RegisterCodec(codec{encoder: encoder, decoder: decoder})
}

type codec struct {
	encoder *schema.Encoder
	decoder *schema.Decoder
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	var vs = url.Values{}
	if err := c.encoder.Encode(v, vs); err != nil {
		return nil, err
	}
	return []byte(vs.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	if err := c.decoder.Decode(v, vs); err != nil {
		return err
	}
	return nil
}

func (codec) Name() string {
	return Name
}
