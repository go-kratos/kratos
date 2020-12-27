package config

import (
	"strings"

	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/provider"

	simplejson "github.com/bitly/go-simplejson"
)

// Resolver is config resolver.
type Resolver interface {
	Resolve(key string) (Value, bool)
}

type resolver struct {
	provider provider.Provider
	parser   map[string]parser.Parser
	kvs      map[string]*simplejson.Json
}

func newResolver(provider provider.Provider, parser map[string]parser.Parser) Resolver {
	r := &resolver{
		provider: provider,
		parser:   parser,
		kvs:      make(map[string]*simplejson.Json),
	}
	r.load()
	return r
}

func (r *resolver) load() error {
	kvs, err := r.provider.Load()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		// TODO parser to json
		raw, err := simplejson.NewJson(kv.Value)
		if err != nil {
			return err
		}
		r.kvs[kv.Key] = raw
	}
	return nil
}

func (r *resolver) Resolve(key string) (Value, bool) {
	path := strings.Split(key, ".")
	for _, v := range r.kvs {
		if raw := v.GetPath(path...); raw.Interface() != nil {
			return &jsonValue{raw: raw}, true
		}
	}
	return nil, false
}
