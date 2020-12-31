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

func newResolver(s source.Source, p map[string]parser.Parser) (*resolver, error) {
	r := &resolver{
		source:  s,
		parsers: p,
	}
	return r, r.load()
}

func (r *resolver) reload(kv *source.KeyValue) error {
	parser, ok := r.parsers[kv.Format]
	if !ok {
		return fmt.Errorf("unsupported parsing formats: %s", kv.Format)
	}
	var v interface{}
	if err := parser.Unmarshal(kv.Value, &v); err != nil {
		return err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	jv := &atomicValue{}
	jv.raw.Store(raw)
	r.values.Store(kv.Key, jv)
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

func (r *resolver) extractValue(values map[string]interface{}, path string) Value {
	keys := strings.Split(path, ".")
	for idx, key := range keys {
		v, ok := values[key]
		if !ok {
			return nil
		}
		if idx == len(keys)-1 {
			jv := &atomicValue{}
			jv.raw.Store(v)
			return jv
		}
		if values, ok = v.(map[string]interface{}); !ok {
			return nil
		}
	}
	return nil
}

func (r *resolver) Resolve(path string) (ret Value) {
	r.values.Range(func(k, v interface{}) bool {
		if values, err := v.(Value).Map(); err == nil {
			if next := r.extractValue(values, path); next != nil {
				ret = next
				return false
			}
		}
		return true
	})
	return
}
