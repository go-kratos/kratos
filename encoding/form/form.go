package form

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/gorilla/schema"
	"net/url"
)

const Name = "x-www-form-urlencoded"

func init() {
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")
	encoder := schema.NewEncoder()
	encoder.SetAliasTag("json")
	encoding.RegisterCodec(codec{encoder: encoder, decoder: decoder})
}

type codec struct {
	encoder *schema.Encoder
	decoder *schema.Decoder
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	var vs = url.Values{}
	if err := c.encoder.Encode(v, vs); err != nil {
		return nil, err
	}
	return []byte(vs.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	if err := c.decoder.Decode(v, vs); err != nil {
		return err
	}
	return nil
}

func (codec) Name() string {
	return Name
}
