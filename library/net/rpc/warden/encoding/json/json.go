package codec

import (
	"bytes"
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

//Reference https://jbrandhorst.com/post/grpc-json/
func init() {
	encoding.RegisterCodec(JSON{
		Marshaler: jsonpb.Marshaler{
			EmitDefaults: true,
			OrigName:     true,
		},
	})
}

// JSON is impl of encoding.Codec
type JSON struct {
	jsonpb.Marshaler
	jsonpb.Unmarshaler
}

// Name is name of JSON
func (j JSON) Name() string {
	return "json"
}

// Marshal is json marshal
func (j JSON) Marshal(v interface{}) (out []byte, err error) {
	if pm, ok := v.(proto.Message); ok {
		b := new(bytes.Buffer)
		err := j.Marshaler.Marshal(b, pm)
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}
	return json.Marshal(v)
}

// Unmarshal is json unmarshal
func (j JSON) Unmarshal(data []byte, v interface{}) (err error) {
	if pm, ok := v.(proto.Message); ok {
		b := bytes.NewBuffer(data)
		return j.Unmarshaler.Unmarshal(b, pm)
	}
	return json.Unmarshal(data, v)
}
