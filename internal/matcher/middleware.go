package matcher

import (
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
)

// Matcher is a middleware matcher.
type Matcher interface {
	Use(ms ...middleware.Middleware)
	Add(selector string, ms ...middleware.Middleware)
	Match(operation string) []middleware.Middleware
}

// New new a middleware matcher.
func New() Matcher {
	return &matcher{
		matchs: make(map[string][]middleware.Middleware),
	}
}

type matcher struct {
	prefix   []string
	defaults []middleware.Middleware
	matchs   map[string][]middleware.Middleware
}

func (m *matcher) Use(ms ...middleware.Middleware) {
	m.defaults = ms
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
