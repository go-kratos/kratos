package util

// CamelCase convert to Camel-Case
func CamelCase(name string) string {
	if len(name) == 0 {
		return name
	}
	if name[0] == '_' {
		return name
	}
	input := []byte(name)
	ret := make([]byte, 0, len(input))
	toUpper := func(c byte) byte {
		if 'a' <= c && c <= 'z' {
			c -= 32
		}
		return c
	}
	ret = append(ret, toUpper(input[0]))
	underline := false
	for _, c := range input[1:] {
		switch {
		case c == '_':
			underline = true
		default:
			if underline {
				c = toUpper(c)
			}
			ret = append(ret, c)
			underline = false
		}
	}
	return string(ret)
}
