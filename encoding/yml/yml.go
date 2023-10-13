package yml

import (
	"gopkg.in/yaml.v3"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for the yml codec.
const Name = "yml"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with yml.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
