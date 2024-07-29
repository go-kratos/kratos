package matcher

import (
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
)

// Matcher is a middleware matcher.
type Matcher interface {
	Use(ms ...middleware.Middleware)
	UseStream(ms ...middleware.StreamMiddleware)
	Add(selector string, ms ...middleware.Middleware)
	Match(operation string) []middleware.Middleware
	AddStream(selector string, ms ...middleware.StreamMiddleware)
	MatchStream(operation string) []middleware.StreamMiddleware
}

// New new a middleware matcher.
func New() Matcher {
	return &matcher{
		matchs:       make(map[string][]middleware.Middleware),
		streamMatchs: make(map[string][]middleware.StreamMiddleware),
	}
}

type matcher struct {
	prefix         []string
	streamPrefix   []string
	defaults       []middleware.Middleware
	streamDefaults []middleware.StreamMiddleware
	matchs         map[string][]middleware.Middleware
	streamMatchs   map[string][]middleware.StreamMiddleware
}

func (m *matcher) Use(ms ...middleware.Middleware) {
	m.defaults = ms
}

func (m *matcher) UseStream(ms ...middleware.StreamMiddleware) {
	m.streamDefaults = ms
}

func (m *matcher) Add(selector string, ms ...middleware.Middleware) {
	if strings.HasSuffix(selector, "*") {
		selector = strings.TrimSuffix(selector, "*")
		m.prefix = append(m.prefix, selector)
		// sort the prefix:
		//  - /foo/bar
		//  - /foo
		sort.Slice(m.prefix, func(i, j int) bool {
			return m.prefix[i] > m.prefix[j]
		})
	}
	m.matchs[selector] = ms
}

func (m *matcher) AddStream(selector string, ms ...middleware.StreamMiddleware) {
	if strings.HasSuffix(selector, "*") {
		selector = strings.TrimSuffix(selector, "*")
		m.streamPrefix = append(m.streamPrefix, selector)
		// sort the prefix:
		//  - /foo/bar
		//  - /foo
		sort.Slice(m.streamPrefix, func(i, j int) bool {
			return m.streamPrefix[i] > m.streamPrefix[j]
		})
	}
	m.streamMatchs[selector] = ms
}

func (m *matcher) Match(operation string) []middleware.Middleware {
	ms := make([]middleware.Middleware, 0, len(m.defaults))
	if len(m.defaults) > 0 {
		ms = append(ms, m.defaults...)
	}
	if next, ok := m.matchs[operation]; ok {
		return append(ms, next...)
	}
	for _, prefix := range m.prefix {
		if strings.HasPrefix(operation, prefix) {
			return append(ms, m.matchs[prefix]...)
		}
	}
	return ms
}

func (m *matcher) MatchStream(operation string) []middleware.StreamMiddleware {
	ms := make([]middleware.StreamMiddleware, 0, len(m.streamDefaults))
	if len(m.streamDefaults) > 0 {
		ms = append(ms, m.streamDefaults...)
	}
	if next, ok := m.streamMatchs[operation]; ok {
		return append(ms, next...)
	}
	for _, prefix := range m.streamPrefix {
		if strings.HasPrefix(operation, prefix) {
			return append(ms, m.streamMatchs[prefix]...)
		}
	}
	return ms
}
