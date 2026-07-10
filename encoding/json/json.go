package json

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v3/encoding"
)

// Name is the name registered for the json codec.
const Name = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v any) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	default:
		return json.Marshal(m)
	}
}

func (codec) Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	default:
		return json.Unmarshal(data, m)
	}
}

func (codec) Name() string {
	return Name
}
