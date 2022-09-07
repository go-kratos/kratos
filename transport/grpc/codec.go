package grpc

import (
	"fmt"

	enc "github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return enc.GetCodec(json.Name).Marshal(vv)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return enc.GetCodec(json.Name).Unmarshal(data, vv)
}

func (codec) Name() string {
	return json.Name
}
