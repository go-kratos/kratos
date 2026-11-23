package cbor

import (
	"github.com/fxamacker/cbor/v2"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for cbor compressor.
const Name = "cbor"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with cbor.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return cbor.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return cbor.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
