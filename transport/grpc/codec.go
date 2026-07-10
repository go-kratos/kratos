package grpc

import (
	"fmt"

	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"

	enc "github.com/go-kratos/kratos/v3/encoding"
	"github.com/go-kratos/kratos/v3/encoding/protojson"
)

const jsonName = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return enc.GetCodec(protojson.Name).Marshal(vv)
}

func (codec) Unmarshal(data []byte, v any) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return enc.GetCodec(protojson.Name).Unmarshal(data, vv)
}

func (codec) Name() string {
	return jsonName
}
