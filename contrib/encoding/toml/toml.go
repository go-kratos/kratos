package toml

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for the toml compressor.
const Name = "toml"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with toml.
type codec struct{}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := toml.NewEncoder(buf)

	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	data := buf.Bytes()
	return data, nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func (c codec) Name() string {
	return Name
}
