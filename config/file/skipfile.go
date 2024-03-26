package file

import (
	"strings"
)

func isSkipFile(name string) bool {
	return strings.HasPrefix(name, ".")
}
