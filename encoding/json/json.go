package json

import (
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
	return protojson.Marshal(v.(proto.Message))
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return protojson.Unmarshal(data, v.(proto.Message))
}

func (codec) Name() string {
	return Name
}
