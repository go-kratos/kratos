package selector

type notMatcher struct {
	subMatcher OperationMatcher
}

func (m notMatcher) Match(operation string) bool {
	return !m.subMatcher.Match(operation)
}

type orMatcher struct {
	subMatchers []OperationMatcher
}

func (m orMatcher) Match(operation string) bool {
	for _, sm := range m.subMatchers {
		if ok := sm.Match(operation); ok {
			return true
		}
	}
	return false
}

type orMatcherBuilder struct {
	subMatchers []OperationMatcher
}

func (b *orMatcherBuilder) build() OperationMatcher {
	if len(b.subMatchers) == 1 {
		return b.subMatchers[0]
	}
	return orMatcher{subMatchers: b.subMatchers}
}

func (b *orMatcherBuilder) push(m ...OperationMatcher) {
	b.subMatchers = append(b.subMatchers, m...)
}

type andMatcher struct {
	subMatchers []OperationMatcher
}

func (m andMatcher) Match(operation string) bool {
	for _, sm := range m.subMatchers {
		if ok := sm.Match(operation); !ok {
			return false
		}
	}
	return true
}
