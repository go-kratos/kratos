package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/trace"
)

const (
	_traceComponentName = "library/cache/redis"
	_tracePeerService   = "redis"
	_traceSpanKind      = "client"
)

var _internalTags = []trace.Tag{
	trace.TagString(trace.TagSpanKind, _traceSpanKind),
	trace.TagString(trace.TagComponent, _traceComponentName),
	trace.TagString(trace.TagPeerService, _tracePeerService),
}

type traceConn struct {
	// tr for pipeline, if tr != nil meaning on pipeline
	tr trace.Trace
	// connTag include e.g. ip,port
	connTags []trace.Tag

	ctx context.Context

	// origin redis conn
	Conn
	pending int
	// TODO: split slow log from trace.
	slowLogThreshold time.Duration
}

func (t *traceConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	statement := getStatement(commandName, args...)
	defer t.slowLog(statement, time.Now())
	// NOTE: ignored empty commandName
	// current sdk will Do empty command after pipeline finished
	if t.tr == nil || commandName == "" {
		return t.Conn.Do(commandName, args...)
	}
	tr := t.tr.Fork("", "Redis:"+commandName)
	tr.SetTag(_internalTags...)
	tr.SetTag(t.connTags...)
	tr.SetTag(trace.TagString(trace.TagDBStatement, statement))
	reply, err = t.Conn.Do(commandName, args...)
	tr.Finish(&err)
	return
}

func (t *traceConn) Send(commandName string, args ...interface{}) (err error) {
	statement := getStatement(commandName, args...)
	defer t.slowLog(statement, time.Now())
	t.pending++
	if t.tr == nil {
		return t.Conn.Send(commandName, args...)
	}
	if t.pending == 1 {
		t.tr = t.tr.Fork("", "Redis:Pipeline")
		t.tr.SetTag(_internalTags...)
		t.tr.SetTag(t.connTags...)
	}
	t.tr.SetLog(
		trace.Log(trace.LogEvent, "Send"),
		trace.Log("db.statement", statement),
	)
	if err = t.Conn.Send(commandName, args...); err != nil {
		t.tr.SetTag(trace.TagBool(trace.TagError, true))
		t.tr.SetLog(
			trace.Log(trace.LogEvent, "Send Fail"),
			trace.Log(trace.LogMessage, err.Error()),
		)
	}
	return err
}

func (t *traceConn) Flush() error {
	defer t.slowLog("Flush", time.Now())
	if t.tr == nil {
		return t.Conn.Flush()
	}
	t.tr.SetLog(trace.Log(trace.LogEvent, "Flush"))
	err := t.Conn.Flush()
	if err != nil {
		t.tr.SetTag(trace.TagBool(trace.TagError, true))
		t.tr.SetLog(
			trace.Log(trace.LogEvent, "Flush Fail"),
			trace.Log(trace.LogMessage, err.Error()),
		)
	}
	return err
}

func (t *traceConn) Receive() (reply interface{}, err error) {
	defer t.slowLog("Receive", time.Now())
	if t.tr == nil {
		return t.Conn.Receive()
	}
	t.tr.SetLog(trace.Log(trace.LogEvent, "Receive"))
	reply, err = t.Conn.Receive()
	if err != nil {
		t.tr.SetTag(trace.TagBool(trace.TagError, true))
		t.tr.SetLog(
			trace.Log(trace.LogEvent, "Receive Fail"),
			trace.Log(trace.LogMessage, err.Error()),
		)
	}
	if t.pending > 0 {
		t.pending--
	}
	if t.pending == 0 {
		t.tr.Finish(nil)
		t.tr = nil
	}
	return reply, err
}

func (t *traceConn) WithContext(ctx context.Context) Conn {
	t.Conn = t.Conn.WithContext(ctx)
	t.tr, _ = trace.FromContext(ctx)
	return t
}

func (t *traceConn) slowLog(statement string, now time.Time) {
	du := time.Since(now)
	if du > t.slowLogThreshold {
		log.Warn("%s slow log statement: %s time: %v", _tracePeerService, statement, du)
	}
}

func getStatement(commandName string, args ...interface{}) (res string) {
	res = commandName
	if len(args) > 0 {
		res = fmt.Sprintf("%s %v", commandName, args[0])
	}
	return
}
