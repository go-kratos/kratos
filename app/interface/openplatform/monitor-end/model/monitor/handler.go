package monitor

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	// app_id
	_appID = "app_id"
	// time format
	_timeFormat = "2006-01-02T15:04:05.999999"
	//  log level name: INFO, WARN...
	_level = "level"
	// log time.
	_time = "time"
	// uniq ID from trace.
	_tid = "traceid"
	// default chan size
	_defaultChan = 2048
	// default agent timeout
	_agentTimeout = xtime.Duration(20 * time.Millisecond)
	// default merge wait time
	_mergeWait = 1 * time.Second
	// max buffer size
	_maxBuffer = 10 * 1024 * 1024 // 10mb
)

var (
	_defaultMonitorConfig = &MonitorConfig{
		Proto: "unixgram",
		Addr:  "/var/run/lancer/collector.sock",
	}
	_defaultTaskIDs = map[string]string{
		env.DeployEnvFat1: "000069",
		env.DeployEnvUat:  "000069",
		env.DeployEnvPre:  "000161",
		env.DeployEnvProd: "000161",
	}
	// log separator
	_logSeparator = []byte("\u0001")
)

// MonitorConfig agent config.
type MonitorConfig struct {
	TaskID  string
	Buffer  int
	Proto   string         `dsn:"network"`
	Addr    string         `dsn:"address"`
	Chan    int            `dsn:"query.chan"`
	Timeout xtime.Duration `dsn:"query.timeout"`
}

// MonitorHandler .
type MonitorHandler struct {
	c      *MonitorConfig
	msgs   chan map[string]interface{}
	waiter sync.WaitGroup
	pool   sync.Pool
}

// NewMonitor a MonitorHandler.
func NewMonitor(c *MonitorConfig) (a *MonitorHandler) {
	if c == nil {
		c = _defaultMonitorConfig
	}
	if c.Buffer == 0 {
		c.Buffer = 1
	}
	if len(c.TaskID) == 0 {
		c.TaskID = _defaultTaskIDs[env.DeployEnv]
	}
	c.Timeout = _agentTimeout
	a = &MonitorHandler{c: c}
	a.pool.New = func() interface{} {
		return make(map[string]interface{}, 20)
	}
	a.msgs = make(chan map[string]interface{}, _defaultChan)
	a.waiter.Add(1)
	go a.writeproc()
	return
}

// Log log to udp statsd daemon.
func (h *MonitorHandler) Log(ctx context.Context, lv log.Level, appID string, args ...log.D) {
	if args == nil {
		return
	}
	d := h.data()
	for _, arg := range args {
		d[arg.Key] = arg.Value
	}
	if t, ok := trace.FromContext(ctx); ok {
		d[_tid] = fmt.Sprintf("%s", t)
	}
	d[_appID] = env.AppID + "." + appID
	d[_level] = lv.String()
	d[_time] = time.Now().Format(_timeFormat)
	select {
	case h.msgs <- d:
	default:
	}
}

// writeproc write data into connection.
func (h *MonitorHandler) writeproc() {
	var (
		buf   bytes.Buffer
		conn  net.Conn
		err   error
		count int
		quit  bool
	)
	defer h.waiter.Done()
	taskID := []byte(h.c.TaskID)
	tick := time.NewTicker(_mergeWait)
	enc := json.NewEncoder(&buf)
	for {
		select {
		case d := <-h.msgs:
			if d == nil {
				quit = true
				goto DUMP
			}
			if buf.Len() >= _maxBuffer {
				buf.Reset() // avoid oom
			}
			now := time.Now()
			buf.Write(taskID)
			buf.Write([]byte(fmt.Sprintf("%d", now.UnixNano()/1e6)))
			enc.Encode(d)
			h.free(d)
			if count++; count < h.c.Buffer {
				buf.Write(_logSeparator)
				continue
			}
		case <-tick.C:
		}
		if conn == nil || err != nil {
			if conn, err = net.DialTimeout(h.c.Proto, h.c.Addr, time.Duration(h.c.Timeout)); err != nil {
				log.Error("net.DialTimeout(%s:%s) error(%v)\n", h.c.Proto, h.c.Addr, err)
				continue
			}
		}
	DUMP:
		if conn != nil && buf.Len() > 0 {
			count = 0
			if _, err = conn.Write(buf.Bytes()); err != nil {
				log.Error("conn.Write(%d bytes) error(%v)\n", buf.Len(), err)
				conn.Close()
			} else {
				// only succeed reset buffer, let conn reconnect.
				log.Info("conn Write(%d bytes) data(%v)\n", buf.Len(), string(buf.Bytes()))
				buf.Reset()
			}
		}
		if quit {
			if conn != nil && err == nil {
				conn.Close()
			}
			return
		}
	}
}

func (h *MonitorHandler) data() map[string]interface{} {
	return h.pool.Get().(map[string]interface{})
}

func (h *MonitorHandler) free(d map[string]interface{}) {
	for k := range d {
		delete(d, k)
	}
	h.pool.Put(d)
}

// Close close the connection.
func (h *MonitorHandler) Close() (err error) {
	h.msgs <- nil
	h.waiter.Wait()
	return nil
}

// SetFormat .
func (h *MonitorHandler) SetFormat(string) {
	// discard setformat
}

// Info .
func (h *MonitorHandler) Info(ctx context.Context, appID string, args ...log.D) {
	h.Log(ctx, log.Level(1), appID, args...)
}

// Warn .
func (h *MonitorHandler) Warn(ctx context.Context, appID string, args ...log.D) {
	h.Log(ctx, log.Level(2), appID, args...)
}

// Error .
func (h *MonitorHandler) Error(ctx context.Context, appID string, args ...log.D) {
	h.Log(ctx, log.Level(3), appID, args...)
}

// CalCode handle codes.
func (a *Log) CalCode() {
	if a.HTTPCode == "" {
		return
	}
	a.Codes = a.HTTPCode
	var (
		codes Codes
		err   error
	)
	if err = json.Unmarshal([]byte(a.HTTPCode), &codes); err != nil {
		log.Warn("s.CalCode error(%+v), codes(%s)", err, a.HTTPCode)
		return
	}
	a.HTTPCode = fmt.Sprintf("%v", codes.HTTPCode)
	a.BusinessCode = fmt.Sprintf("%v", codes.HTTPBusinessCode)
	a.InnerCode = fmt.Sprintf("%v", codes.HTTPInnerCode)
	// 电商inner_code 覆盖 business_code
	if a.InnerCode != "-1" {
		// 电商code 1 转成 0
		if a.InnerCode == "1" {
			a.BusinessCode = "0"
		} else {
			a.BusinessCode = a.InnerCode
		}
	}
	if a.BusinessCode == "-1" {
		a.BusinessCode = "0"
	}
}
