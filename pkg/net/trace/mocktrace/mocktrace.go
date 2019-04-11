package mocktrace

import (
	"github.com/bilibili/kratos/pkg/net/trace"
)

// MockTrace .
type MockTrace struct {
	Spans []*MockSpan
}

// New .
func (m *MockTrace) New(operationName string, opts ...trace.Option) trace.Trace {
	span := &MockSpan{OperationName: operationName, MockTrace: m}
	m.Spans = append(m.Spans, span)
	return span
}

// Inject .
func (m *MockTrace) Inject(t trace.Trace, format interface{}, carrier interface{}) error {
	return nil
}

// Extract .
func (m *MockTrace) Extract(format interface{}, carrier interface{}) (trace.Trace, error) {
	return &MockSpan{}, nil
}

// MockSpan .
type MockSpan struct {
	*MockTrace
	OperationName string
	FinishErr     error
	Finished      bool
	Tags          []trace.Tag
	Logs          []trace.LogField
}

// Fork .
func (m *MockSpan) Fork(serviceName string, operationName string) trace.Trace {
	span := &MockSpan{OperationName: operationName, MockTrace: m.MockTrace}
	m.Spans = append(m.Spans, span)
	return span
}

// Follow .
func (m *MockSpan) Follow(serviceName string, operationName string) trace.Trace {
	span := &MockSpan{OperationName: operationName, MockTrace: m.MockTrace}
	m.Spans = append(m.Spans, span)
	return span
}

// Finish .
func (m *MockSpan) Finish(perr *error) {
	if perr != nil {
		m.FinishErr = *perr
	}
	m.Finished = true
}

// SetTag .
func (m *MockSpan) SetTag(tags ...trace.Tag) trace.Trace {
	m.Tags = append(m.Tags, tags...)
	return m
}

// SetLog .
func (m *MockSpan) SetLog(logs ...trace.LogField) trace.Trace {
	m.Logs = append(m.Logs, logs...)
	return m
}

// Visit .
func (m *MockSpan) Visit(fn func(k, v string)) {}

// SetTitle .
func (m *MockSpan) SetTitle(title string) {
	m.OperationName = title
}

// TraceID .
func (m *MockSpan) TraceID() string {
	return ""
}
