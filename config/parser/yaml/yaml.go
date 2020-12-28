package yaml

import (
	"github.com/ghodss/yaml"
	"github.com/go-kratos/kratos/v2/config/parser"
)

type yamlParser struct{}

// NewParser new a yaml parser.
func NewParser() parser.Parser {
	return &yamlParser{}
}

func (p *yamlParser) Format() string {
	return "yaml"
}

func (p *yamlParser) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (p *yamlParser) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
