package env

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
)

type env struct {
	prefixes []string
}

func NewSource(prefixes ...string) config.Source {
	return &env{prefixes: prefixes}
}

func (e *env) Load() (kvs []*config.KeyValue, err error) {
	return e.load(os.Environ()), nil
}

func (e *env) load(envs []string) []*config.KeyValue {
	var kvs []*config.KeyValue
	for _, env := range envs {
		k, v, _ := strings.Cut(env, "=")
		if k == "" {
			continue
		}
		if len(e.prefixes) > 0 {
			prefix, ok := matchPrefix(e.prefixes, k)
			if !ok || k == prefix {
				continue
			}
			k = strings.TrimPrefix(k, prefix)
			k = strings.TrimPrefix(k, "_")
		}
		if k != "" {
			kvs = append(kvs, &config.KeyValue{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	return kvs
}

func (e *env) Watch() (config.Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func matchPrefix(prefixes []string, s string) (string, bool) {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}
