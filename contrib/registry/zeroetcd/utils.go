package etcd

import (
	"fmt"
	"strings"
)

const (
	indexOfKey = iota
	indexOfId
)

const(
	Delimiter = '/'
)
// TimeToLive is seconds to live in etcd.

func extract(etcdKey string, index int) (string, bool) {
	if index < 0 {
		return "", false
	}

	fields := strings.FieldsFunc(etcdKey, func(ch rune) bool {
		return ch == Delimiter
	})
	if index >= len(fields) {
		return "", false
	}

	return fields[index], true
}

func extractId(etcdKey string) (string, bool) {
	return extract(etcdKey, indexOfId)
}

func extractKey(etcdKey string) (string, bool) {
	return extract(etcdKey, indexOfKey)
}

func makeEtcdKey(key string, id int64) string {
	return fmt.Sprintf("%s%c%d", key, Delimiter, id)
}

