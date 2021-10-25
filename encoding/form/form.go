package form

import (
	"net/url"
	"reflect"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/go-kratos/kratos/v2/encoding"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

const (
	// Name is form codec name
	Name = "x-www-form-urlencoded"
	// Null value string
	nullStr = "null"
)

func init() {
	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	encoder := form.NewEncoder()
	encoder.SetTagName("json")
	encoding.RegisterCodec(codec{encoder: encoder, decoder: decoder})
}

type codec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	var vs url.Values
	var err error
	if m, ok := v.(proto.Message); ok {
		vs, err = EncodeMap(m)
		if err != nil {
			return nil, err
		}
	} else {
		vs, err = c.encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if m, ok := v.(proto.Message); ok {
		return MapProto(m, vs)
	} else if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
		return MapProto(m, vs)
	} else if _, ok := v.(*map[string]string); ok {
		// form body bind map<string,string>.
		vd := make(map[string]string)
		for key, values := range vs {
			vd[key] = values[0]
		}
		rv.Set(reflect.ValueOf(vd))
		return nil
	} else if _, ok := v.(*map[string]*structpb.Value); ok {
		// form body bind map<string, google.protobuf.Value>
		vd := make(map[string]*structpb.Value)
		for key, values := range vs {
			value, err := structpb.NewValue(values[0])
			if err != nil {
				return err
			}
			vd[key] = value
		}
		rv.Set(reflect.ValueOf(vd))
		return nil
	}

	return c.decoder.Decode(v, vs)
}

func (codec) Name() string {
	return Name
}
