package form

import (
	"net/url"
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

// Name is form codec name
const Name = "x-www-form-urlencoded"

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
	}

	return c.decoder.Decode(v, vs)
}

func (codec) Name() string {
	return Name
}
