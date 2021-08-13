package selector

import (
	"regexp"
	"strings"
)

type OperationMatcher interface {
	Match(operation string) bool
}

type operationMatcherFunc func(operation string) bool

func (f operationMatcherFunc) Match(operation string) bool {
	return f(operation)
}

type pathMather string

func (m pathMather) Match(operation string) bool {
	return string(m) == operation
}

type prefixMather string

func (m prefixMather) Match(operation string) bool {
	return strings.HasPrefix(operation, string(m))
}

type regexMatcher struct {
	re *regexp.Regexp
}

func (m regexMatcher) Match(operation string) bool {
	return m.re.FindString(operation) == operation
}
