package memcache

import (
	"context"
)

// MockErr for unit test.
type MockErr struct {
	Error error
}

var _ Conn = MockErr{}

// MockWith return a mock conn.
func MockWith(err error) MockErr {
	return MockErr{Error: err}
}

// Err .
func (m MockErr) Err() error { return m.Error }

// Close .
func (m MockErr) Close() error { return m.Error }

// Add .
func (m MockErr) Add(item *Item) error { return m.Error }

// Set .
func (m MockErr) Set(item *Item) error { return m.Error }

// Replace .
func (m MockErr) Replace(item *Item) error { return m.Error }

// CompareAndSwap .
func (m MockErr) CompareAndSwap(item *Item) error { return m.Error }

// Get .
func (m MockErr) Get(key string) (*Item, error) { return nil, m.Error }

// GetMulti .
func (m MockErr) GetMulti(keys []string) (map[string]*Item, error) { return nil, m.Error }

// Touch .
func (m MockErr) Touch(key string, timeout int32) error { return m.Error }

// Delete .
func (m MockErr) Delete(key string) error { return m.Error }

// Increment .
func (m MockErr) Increment(key string, delta uint64) (uint64, error) { return 0, m.Error }

// Decrement .
func (m MockErr) Decrement(key string, delta uint64) (uint64, error) { return 0, m.Error }

// Scan .
func (m MockErr) Scan(item *Item, v interface{}) error { return m.Error }

// WithContext .
func (m MockErr) WithContext(ctx context.Context) Conn { return m }
