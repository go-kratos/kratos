package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/source"
)

type resolver struct {
	source  source.Source
	parsers map[string]parser.Parser
	values  sync.Map
}

func newResolver(s source.Source, opts options) (*resolver, error) {
	r := &resolver{
		source:  s,
		parsers: make(map[string]parser.Parser),
	}
	for _, parser := range opts.parsers {
		r.parsers[parser.Format()] = parser
	}
	return r, r.load()
}

func (r *resolver) reload(kv *source.KeyValue) error {
	parser, ok := r.parsers[kv.Format]
	if !ok {
		return fmt.Errorf("unsupported parsing formats: %s", kv.Format)
	}
	var m interface{}
	if err := parser.Unmarshal(kv.Value, &m); err != nil {
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	av := &atomicValue{}
	av.Store(v)
	r.values.Store(kv.Key, av)
	return nil
}

func (r *resolver) load() error {
	kvs, err := r.source.Load()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		if err := r.reload(kv); err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) extractValue(values interface{}, path string) Value {
	next, ok := values.(map[string]interface{})
	if !ok {
		return nil
	}
	keys := strings.Split(path, ".")
	for idx, key := range keys {
		v, ok := next[key]
		if !ok {
			return nil
		}
		if idx == len(keys)-1 {
			av := &atomicValue{}
			av.Store(v)
			return av
		}
		if next, ok = v.(map[string]interface{}); !ok {
			return nil
		}
	}
	return nil
}

func (r *resolver) Resolve(path string) (ret Value) {
	r.values.Range(func(k, v interface{}) bool {
		if values := v.(Value).Load(); values != nil {
			if next := r.extractValue(values, path); next != nil {
				ret = next
				return false
			}
		}
		return true
	})
	return
}
