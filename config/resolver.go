package config

import (
	"github.com/go-kratos/kratos/v2/config/parser"
	"github.com/go-kratos/kratos/v2/config/provider"
)

// Resolver is config resolver.
type Resolver interface {
	Resolve(key string) (Value, bool)
}

type resolver struct {
	provider provider.Provider
	parser   map[string]parser.Parser
	cached   map[string]Value
}

func newResolver(provider provider.Provider, parser map[string]parser.Parser) Resolver {
	return &resolver{provider: provider, parser: parser}
}

func (r *resolver) load() {
	kvs, err := r.provider.Load()
	if err != nil {
		return
	}
}

func (r *resolver) Resolve(key string) (Value, bool) {
	if v, ok := r.cached[key]; ok {
		return v, true
	}
	return nil, false
}
