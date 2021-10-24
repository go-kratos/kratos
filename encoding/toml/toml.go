package toml

import (
	"github.com/go-kratos/kratos/v2/encoding"

	"github.com/pelletier/go-toml"
)

const Name = "toml"

func init() {
	encoding.RegisterCodec(codec{})
}

type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
