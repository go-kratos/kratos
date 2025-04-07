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

func (e *env) Load() (kv []*config.KeyValue, err error) {
	return e.load(os.Environ()), nil
}

func (e *env) load(envs []string) []*config.KeyValue {
	var kv []*config.KeyValue
	for _, env := range envs {
		var k, v string
		subs := strings.SplitN(env, "=", 2) //nolint:mnd
		k = subs[0]
		if len(subs) > 1 {
			v = subs[1]
		}

		if len(e.prefixes) > 0 {
			p, ok := matchPrefix(e.prefixes, k)
			if !ok || len(p) == len(k) {
				continue
			}
			// trim prefix
			k = strings.TrimPrefix(k, p)
			k = strings.TrimPrefix(k, "_")
		}

		if len(k) != 0 {
			kv = append(kv, &config.KeyValue{
				Key:   k,
				Value: []byte(v),
			})
		}
	}
	return kv
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
