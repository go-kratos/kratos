package redis

import (
	"context"
)

// MockErr for unit test.
type MockErr struct {
	Error error
}

// MockWith return a mock conn.
func MockWith(err error) MockErr {
	return MockErr{Error: err}
}

// Err .
func (m MockErr) Err() error { return m.Error }

// Close .
func (m MockErr) Close() error { return m.Error }

// Do .
func (m MockErr) Do(commandName string, args ...interface{}) (interface{}, error) { return nil, m.Error }

// Send .
func (m MockErr) Send(commandName string, args ...interface{}) error { return m.Error }

// Flush .
func (m MockErr) Flush() error { return m.Error }

// Receive .
func (m MockErr) Receive() (interface{}, error) { return nil, m.Error }

// WithContext .
func (m MockErr) WithContext(context.Context) Conn { return m }
