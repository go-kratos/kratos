package gock

import "net/http"

// MatchersHeader exposes an slice of HTTP header specific mock matchers.
var MatchersHeader = []MatchFunc{
	MatchMethod,
	MatchScheme,
	MatchHost,
	MatchPath,
	MatchHeaders,
	MatchQueryParams,
}

// MatchersBody exposes an slice of HTTP body specific built-in mock matchers.
var MatchersBody = []MatchFunc{
	MatchBody,
}

// Matchers stores all the built-in mock matchers.
var Matchers = append(MatchersHeader, MatchersBody...)

// DefaultMatcher stores the default Matcher instance used to match mocks.
var DefaultMatcher = NewMatcher()

// MatchFunc represents the required function
// interface implemented by matchers.
type MatchFunc func(*http.Request, *Request) (bool, error)

// Matcher represents the required interface implemented by mock matchers.
type Matcher interface {
	// Get returns a slice of registered function matchers.
	Get() []MatchFunc

	// Add adds a new matcher function.
	Add(MatchFunc)

	// Set sets the matchers functions stack.
	Set([]MatchFunc)

	// Flush flushes the current matchers function stack.
	Flush()

	// Match matches the given http.Request with a mock Request.
	Match(*http.Request, *Request) (bool, error)
}

// MockMatcher implements a mock matcher
type MockMatcher struct {
	Matchers []MatchFunc
}

// NewMatcher creates a new mock matcher
// using the default matcher functions.
func NewMatcher() *MockMatcher {
	return &MockMatcher{Matchers: Matchers}
}

// NewBasicMatcher creates a new matcher with header only mock matchers.
func NewBasicMatcher() *MockMatcher {
	return &MockMatcher{Matchers: MatchersHeader}
}

// NewEmptyMatcher creates a new empty matcher with out default amtchers.
func NewEmptyMatcher() *MockMatcher {
	return &MockMatcher{Matchers: []MatchFunc{}}
}

// Get returns a slice of registered function matchers.
func (m *MockMatcher) Get() []MatchFunc {
	return m.Matchers
}

// Add adds a new function matcher.
func (m *MockMatcher) Add(fn MatchFunc) {
	m.Matchers = append(m.Matchers, fn)
}

// Set sets a new stack of matchers functions.
func (m *MockMatcher) Set(stack []MatchFunc) {
	m.Matchers = stack
}

// Flush flushes the current matcher
func (m *MockMatcher) Flush() {
	m.Matchers = []MatchFunc{}
}

// Match matches the given http.Request with a mock request
// returning true in case that the request matches, otherwise false.
func (m *MockMatcher) Match(req *http.Request, ereq *Request) (bool, error) {
	for _, matcher := range m.Matchers {
		matches, err := matcher(req, ereq)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

// MatchMock is a helper function that matches the given http.Request
// in the list of registered mocks, returning it if matches or error if it fails.
func MatchMock(req *http.Request) (Mock, error) {
	for _, mock := range GetAll() {
		matches, err := mock.Match(req)
		if err != nil {
			return nil, err
		}
		if matches {
			return mock, nil
		}
	}
	return nil, nil
}
