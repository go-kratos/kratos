package apollo

import (
	"encoding/json"

	"github.com/apolloconfig/agollo/v4/constant"
	"github.com/apolloconfig/agollo/v4/extension"
)

type jsonExtParser struct{}

func (j jsonExtParser) Parse(configContent interface{}) (map[string]interface{}, error) {
	v, ok := configContent.(string)
	if !ok {
		return nil, nil
	}
	out := make(map[string]interface{}, 4)
	err := json.Unmarshal([]byte(v), &out)
	return out, err
}

func init() {
	// add json format
	extension.AddFormatParser(constant.JSON, &jsonExtParser{})
}
