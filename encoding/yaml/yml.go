package yaml

import "github.com/go-kratos/kratos/v2/encoding"

const NameAlias = "yml"

type codecAlias struct {
	codec
}

func init() {
	encoding.RegisterCodec(codecAlias{})
}

func (codecAlias) Name() string {
	return NameAlias
}
