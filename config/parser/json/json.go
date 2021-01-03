package json

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/config/parser"
)

var _ parser.Parser = (*jsonParser)(nil)

type jsonParser struct{}

// NewParser new a json parser.
func NewParser() parser.Parser {
	return &jsonParser{}
}

func (p *jsonParser) Format() string {
	return "json"
}

func (p *jsonParser) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (p *jsonParser) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
