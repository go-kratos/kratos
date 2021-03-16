package json

import (
	"encoding/json"
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for the json codec.
const Name = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() && rv.CanInterface() {
			rv.Set(reflect.New(rv.Type().Elem()))
			v = rv.Interface()
		}
		rv = rv.Elem()
	}
	return json.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
