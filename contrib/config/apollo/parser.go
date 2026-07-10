package apollo

const contentKey = "content"

type jsonExtParser struct{}

func (parser jsonExtParser) Parse(configContent any) (map[string]any, error) {
	return map[string]any{contentKey: configContent}, nil
}

type yamlExtParser struct{}

func (parser yamlExtParser) Parse(configContent any) (map[string]any, error) {
	return map[string]any{contentKey: configContent}, nil
}
