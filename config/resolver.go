package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/provider"
)

// Resolver is config resolver.
type Resolver interface {
	Resolve(key string) (Value, bool)
}

type resolver struct {
	provider provider.Provider
	parsers  map[string]parser.Parser
	values   map[string]jsonValue
}

func newResolver(provider provider.Provider, parsers map[string]parser.Parser) (Resolver, error) {
	r := &resolver{
		provider: provider,
		parsers:  parsers,
		values:   make(map[string]jsonValue),
	}
	return r, r.load()
}

func (r *resolver) load() error {
	kvs, err := r.provider.Load()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
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
		jv := jsonValue{}
		if err := json.Unmarshal(data, &jv.raw); err != nil {
			return err
		}
		r.values[kv.Key] = jv
	}
	return nil
}

func (r *resolver) Resolve(key string) (Value, bool) {
	path := strings.Split(key, ".")
	for _, v := range r.values {
		if val := v.GetPath(path...); val.raw != nil {
			return &jsonValue{raw: val.raw}, true
		}
	}
	return nil, false
}
