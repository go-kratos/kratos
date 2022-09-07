package apollo

type jsonExtParser struct{}

func (parser jsonExtParser) Parse(configContent interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"content": configContent}, nil
}

type yamlExtParser struct{}

func (parser yamlExtParser) Parse(configContent interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"content": configContent}, nil
}
