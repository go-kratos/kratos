package file

import "strings"

func format(name string) string {
	if idx := strings.LastIndexByte(name, '.'); idx >= 0 {
		return name[idx+1:]
	}
	return ""
}
