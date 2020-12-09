package json

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/encoding"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the json codec.
const Name = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	m, ok := v.(proto.Message)
	if ok {
		return protojson.Marshal(m)
	}
	return json.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	m, ok := v.(proto.Message)
	if ok {
		return protojson.Unmarshal(data, m)
	}
	return json.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
