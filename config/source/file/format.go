package file

import "strings"

func format(name string) string {
	if ext := strings.Split(name, "."); len(ext) > 1 {
		return ext[len(ext)-1]
	}
	return "text"
}
