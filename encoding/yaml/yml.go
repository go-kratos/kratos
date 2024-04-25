package yaml

import (
	"github.com/go-kratos/kratos/v2/encoding"
)

func init() {
	encoding.RegisterCodec(ymlCodec{})
}

// ymlCodec is a Codec inheriting from the default yaml codec
// to support multiple name registrations
type ymlCodec struct {
	codec
}

func (ymlCodec) Name() string {
	return "yml"
}
