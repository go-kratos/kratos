package utils

import (
	"os"
	"unicode"
)

// LcFirst lower the first letter
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func IsDir(name string) bool {
	file, err := os.Open(name)

	if err != nil {
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}
	return fi.IsDir()
}
