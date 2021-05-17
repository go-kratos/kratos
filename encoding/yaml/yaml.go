package yaml

import (
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"
	"gopkg.in/yaml.v2"
)

// Name is the name registered for the json codec.
const Name = "yaml"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return yaml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
