package dao

import (
	"fmt"
	"strconv"

	"github.com/dgryski/go-farm"
)

func rangeKey(prefix string, start, end int64) (string, string) {
	return prefix + strconv.FormatInt(start, 10), prefix + strconv.FormatInt(end, 10)
}

func keyPrefix(serviceName, operationName string) string {
	serviceNameHash := farm.Hash32([]byte(serviceName))
	operationNameHash := farm.Hash32([]byte(operationName))
	return fmt.Sprintf("%x%x", serviceNameHash, operationNameHash)
}
