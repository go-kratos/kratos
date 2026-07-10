package protojson

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/go-kratos/kratos/v3/encoding"
)

// Name is the name registered for the protojson codec.
const Name = "protojson"

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protojson.
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	m, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return MarshalOptions.Marshal(m)
}

func (codec) Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	m, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return UnmarshalOptions.Unmarshal(data, m)
}

func (codec) Name() string {
	return Name
}
