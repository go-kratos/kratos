package json

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/config/parser"
)

type jsonParser struct{}

// NewParser new a json parser.
func NewParser() parser.Parser {
	return &jsonParser{}
}

func (j *jsonParser) Format() string {
	return "json"
}

func (j *jsonParser) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonParser) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
