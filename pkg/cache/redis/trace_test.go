package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/net/trace"
	"github.com/stretchr/testify/assert"
)

const testTraceSlowLogThreshold = time.Duration(250 * time.Millisecond)

type mockTrace struct {
	tags          []trace.Tag
	logs          []trace.LogField
	perr          *error
	operationName string
	finished      bool
}

func (m *mockTrace) Fork(serviceName string, operationName string) trace.Trace {
	m.operationName = operationName
	return m
}
func (m *mockTrace) Follow(serviceName string, operationName string) trace.Trace {
	panic("not implemented")
}
func (m *mockTrace) Finish(err *error) {
	m.perr = err
	m.finished = true
}
func (m *mockTrace) SetTag(tags ...trace.Tag) trace.Trace {
	m.tags = append(m.tags, tags...)
	return m
}
func (m *mockTrace) SetLog(logs ...trace.LogField) trace.Trace {
	m.logs = append(m.logs, logs...)
	return m
}
func (m *mockTrace) Visit(fn func(k, v string)) {}
func (m *mockTrace) SetTitle(title string)      {}
func (m *mockTrace) TraceID() string            { return "" }

type mockConn struct{}

func (c *mockConn) Close() error { return nil }
func (c *mockConn) Err() error   { return nil }
func (c *mockConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return nil, nil
}
func (c *mockConn) Send(commandName string, args ...interface{}) error { return nil }
func (c *mockConn) Flush() error                                       { return nil }
func (c *mockConn) Receive() (reply interface{}, err error)            { return nil, nil }
func (c *mockConn) WithContext(context.Context) Conn                   { return c }

func TestTraceDo(t *testing.T) {
	tr := &mockTrace{}
	ctx := trace.NewContext(context.Background(), tr)
	tc := &traceConn{Conn: &mockConn{}, slowLogThreshold: testTraceSlowLogThreshold}
	conn := tc.WithContext(ctx)

	conn.Do("GET", "test")

	assert.Equal(t, "Redis:GET", tr.operationName)
	assert.NotEmpty(t, tr.tags)
	assert.True(t, tr.finished)
}

func TestTraceDoErr(t *testing.T) {
	tr := &mockTrace{}
	ctx := trace.NewContext(context.Background(), tr)
	tc := &traceConn{Conn: MockErr{Error: fmt.Errorf("hhhhhhh")},
		slowLogThreshold: testTraceSlowLogThreshold}
	conn := tc.WithContext(ctx)

	conn.Do("GET", "test")

	assert.Equal(t, "Redis:GET", tr.operationName)
	assert.True(t, tr.finished)
	assert.NotNil(t, *tr.perr)
}

func TestTracePipeline(t *testing.T) {
	tr := &mockTrace{}
	ctx := trace.NewContext(context.Background(), tr)
	tc := &traceConn{Conn: &mockConn{}, slowLogThreshold: testTraceSlowLogThreshold}
	conn := tc.WithContext(ctx)

	N := 2
	for i := 0; i < N; i++ {
		conn.Send("GET", "hello, world")
	}
	conn.Flush()
	for i := 0; i < N; i++ {
		conn.Receive()
	}

	assert.Equal(t, "Redis:Pipeline", tr.operationName)
	assert.NotEmpty(t, tr.tags)
	assert.NotEmpty(t, tr.logs)
	assert.True(t, tr.finished)
}

func TestTracePipelineErr(t *testing.T) {
	tr := &mockTrace{}
	ctx := trace.NewContext(context.Background(), tr)
	tc := &traceConn{Conn: MockErr{Error: fmt.Errorf("hahah")},
		slowLogThreshold: testTraceSlowLogThreshold}
	conn := tc.WithContext(ctx)

	N := 2
	for i := 0; i < N; i++ {
		conn.Send("GET", "hello, world")
	}
	conn.Flush()
	for i := 0; i < N; i++ {
		conn.Receive()
	}

	assert.Equal(t, "Redis:Pipeline", tr.operationName)
	assert.NotEmpty(t, tr.tags)
	assert.NotEmpty(t, tr.logs)
	assert.True(t, tr.finished)
	var isError bool
	for _, tag := range tr.tags {
		if tag.Key == "error" {
			isError = true
		}
	}
	assert.True(t, isError)
}

func TestSendStatement(t *testing.T) {
	tr := &mockTrace{}
	ctx := trace.NewContext(context.Background(), tr)
	tc := &traceConn{Conn: MockErr{Error: fmt.Errorf("hahah")},
		slowLogThreshold: testTraceSlowLogThreshold}
	conn := tc.WithContext(ctx)
	conn.Send("SET", "hello", "test")
	conn.Flush()
	conn.Receive()

	assert.Equal(t, "Redis:Pipeline", tr.operationName)
	assert.NotEmpty(t, tr.tags)
	assert.NotEmpty(t, tr.logs)
	assert.Equal(t, "event", tr.logs[0].Key)
	assert.Equal(t, "Send", tr.logs[0].Value)
	assert.Equal(t, "db.statement", tr.logs[1].Key)
	assert.Equal(t, "SET hello", tr.logs[1].Value)
	assert.True(t, tr.finished)
	var isError bool
	for _, tag := range tr.tags {
		if tag.Key == "error" {
			isError = true
		}
	}
	assert.True(t, isError)
}

func TestDoStatement(t *testing.T) {
	tr := &mockTrace{}
	ctx := trace.NewContext(context.Background(), tr)
	tc := &traceConn{Conn: MockErr{Error: fmt.Errorf("hahah")},
		slowLogThreshold: testTraceSlowLogThreshold}
	conn := tc.WithContext(ctx)
	conn.Do("SET", "hello", "test")

	assert.Equal(t, "Redis:SET", tr.operationName)
	assert.Equal(t, "SET hello", tr.tags[len(tr.tags)-1].Value)
	assert.True(t, tr.finished)
}

func BenchmarkTraceConn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c, err := DialDefaultServer()
		if err != nil {
			b.Fatal(err)
		}
		t := &traceConn{
			Conn:             c,
			connTags:         []trace.Tag{trace.TagString(trace.TagPeerAddress, "abc")},
			slowLogThreshold: time.Duration(1 * time.Second),
		}
		c2 := t.WithContext(context.TODO())
		if _, err := c2.Do("PING"); err != nil {
			b.Fatal(err)
		}
		c2.Close()
	}
}
