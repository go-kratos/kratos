package text

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/config/parser"
)

var _ parser.Parser = (*textParser)(nil)

type textParser struct{}

// NewParser new a json parser.
func NewParser() parser.Parser {
	return &textParser{}
}

func (p *textParser) Format() string {
	return "text"
}

func (p *textParser) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (p *textParser) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
