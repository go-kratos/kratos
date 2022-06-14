package apollo

import (
	"github.com/apolloconfig/agollo/v4/constant"
	"github.com/apolloconfig/agollo/v4/extension"
)

type jsonExtParser struct{}

func (parser jsonExtParser) Parse(configContent interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"content": configContent}, nil
}

type yamlExtParser struct{}

func (parser yamlExtParser) Parse(configContent interface{}) (out map[string]interface{}, err error) {
	return map[string]interface{}{"content": configContent}, nil
}

func init() {
	// add json/yaml/yml format
	extension.AddFormatParser(constant.JSON, &jsonExtParser{})
	extension.AddFormatParser(constant.YAML, &yamlExtParser{})
	extension.AddFormatParser(constant.YML, &yamlExtParser{})
}
