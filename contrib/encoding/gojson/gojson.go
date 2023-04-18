package gojson

import (
	"github.com/goccy/go-json"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for the goccy/go-json compressor.
const Name = "gojson"

func init() {
	encoding.RegisterCodec(GoJson{})
}

// GoJson is a Codec implementation with goccy/go-json.
type GoJson struct {
}

func (GoJson) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (GoJson) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (GoJson) Name() string {
	return Name
}
