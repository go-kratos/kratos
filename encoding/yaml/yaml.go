package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for the yaml codec.
const Name = "yaml"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with yaml.
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
