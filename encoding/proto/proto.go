package proto

import (
	"github.com/go-kratos/kratos/v2/encoding"

	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the proto codec.
const Name = "proto"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (codec) Name() string {
	return Name
}
