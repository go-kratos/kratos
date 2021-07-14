package env

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
)

type env struct {
	prefixs []string
}

func NewSource(prefixs ...string) config.Source {
	return &env{prefixs: prefixs}
}

func (e *env) Load() (kv []*config.KeyValue, err error) {
	for _, envstr := range os.Environ() {
		var k, v string
		subs := strings.SplitN(envstr, "=", 2)
		k = subs[0]
		if len(subs) > 1 {
			v = subs[1]
		}

		if len(e.prefixs) > 0 {
			p, ok := matchPrefix(e.prefixs, envstr)
			if !ok {
				continue
			}
			// trim prefix
			k = k[len(p):]
			if k[0] == '_' {
				k = k[1:]
			}
		}

		kv = append(kv, &config.KeyValue{
			Key:   k,
			Value: []byte(v),
		})
	}
	return
}

func (e *env) Watch() (config.Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func matchPrefix(prefixs []string, s string) (string, bool) {
	for _, p := range prefixs {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}
