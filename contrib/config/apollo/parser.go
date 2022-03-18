package apollo

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/apolloconfig/agollo/v4/constant"
	"github.com/apolloconfig/agollo/v4/extension"
)

type jsonExtParser struct{}

func (parser jsonExtParser) Parse(configContent interface{}) (map[string]interface{}, error) {
	v, ok := configContent.(string)
	if !ok {
		return nil, nil
	}
	out := make(map[string]interface{}, 4)
	err := json.Unmarshal([]byte(v), &out)
	return out, err
}

type yamlExtParser struct{}

func (parser yamlExtParser) Parse(configContent interface{}) (out map[string]interface{}, err error) {
	v, ok := configContent.(string)
	if !ok {
		return nil, nil
	}
	out = make(map[string]interface{}, 4)
	err = yaml.Unmarshal([]byte(v), &out)
	return
}

func init() {
	// add json/yaml/yml format
	extension.AddFormatParser(constant.JSON, &jsonExtParser{})
	extension.AddFormatParser(constant.YAML, &yamlExtParser{})
	extension.AddFormatParser(constant.YML, &yamlExtParser{})
}
