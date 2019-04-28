package log

import (
	"context"
	"fmt"
	stdlog "log"
	"net"
	"strconv"
	"sync"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log/internal"
	"go-common/library/net/metadata"
	"go-common/library/net/trace"
	xtime "go-common/library/time"
)

const (
	_agentTimeout = xtime.Duration(20 * time.Millisecond)
	_mergeWait    = 1 * time.Second
	_maxBuffer    = 10 * 1024 * 1024 // 10mb
	_defaultChan  = 2048

	_defaultAgentConfig = "unixpacket:///var/run/lancer/collector_tcp.sock?timeout=100ms&chan=1024"
)

var (
	_logSeparator = []byte("\u0001")

	_defaultTaskIDs = map[string]string{
		env.DeployEnvFat1: "000069",
		env.DeployEnvUat:  "000069",
		env.DeployEnvPre:  "000161",
		env.DeployEnvProd: "000161",
	}
)

// AgentHandler agent struct.
type AgentHandler struct {
	c         *AgentConfig
	msgs      chan []core.Field
	waiter    sync.WaitGroup
	pool      sync.Pool
	enc       core.Encoder
	batchSend bool
	filters   map[string]struct{}
}

// AgentConfig agent config.
type AgentConfig struct {
	TaskID  string
	Buffer  int
	Proto   string         `dsn:"network"`
	Addr    string         `dsn:"address"`
	Chan    int            `dsn:"query.chan"`
	Timeout xtime.Duration `dsn:"query.timeout"`
}

// NewAgent a Agent.
func NewAgent(ac *AgentConfig) (a *AgentHandler) {
	if ac == nil {
		ac = parseDSN(_agentDSN)
	}
	if len(ac.TaskID) == 0 {
		ac.TaskID = _defaultTaskIDs[env.DeployEnv]
	}
	a = &AgentHandler{
		c: ac,
		enc: core.NewJSONEncoder(core.EncoderConfig{
			EncodeTime:     core.EpochTimeEncoder,
			EncodeDuration: core.SecondsDurationEncoder,
		}, core.NewBuffer(0)),
	}
	a.pool.New = func() interface{} {
		return make([]core.Field, 0, 16)
	}
	if ac.Chan == 0 {
		ac.Chan = _defaultChan
	}
	a.msgs = make(chan []core.Field, ac.Chan)
	if ac.Timeout == 0 {
		ac.Timeout = _agentTimeout
	}
	if ac.Buffer == 0 {
		ac.Buffer = 100
	}
	a.waiter.Add(1)

	// set fixed k/v into enc buffer
	KV(_appID, c.Family).AddTo(a.enc)
	KV(_deplyEnv, env.DeployEnv).AddTo(a.enc)
	KV(_instanceID, c.Host).AddTo(a.enc)
	KV(_zone, env.Zone).AddTo(a.enc)

	if a.c.Proto == "unixpacket" {
		a.batchSend = true
	}

	go a.writeproc()
	return
}

func (h *AgentHandler) data() []core.Field {
	return h.pool.Get().([]core.Field)
}

func (h *AgentHandler) free(f []core.Field) {
	f = f[0:0]
	h.pool.Put(f)
}

// Log log to udp statsd daemon.
func (h *AgentHandler) Log(ctx context.Context, lv Level, args ...D) {
	if args == nil {
		return
	}
	f := h.data()
	for i := range args {
		f = append(f, args[i])
	}
	if t, ok := trace.FromContext(ctx); ok {
		if s, ok := t.(fmt.Stringer); ok {
			f = append(f, KV(_tid, s.String()))
		} else {
			f = append(f, KV(_tid, fmt.Sprintf("%s", t)))
		}
	}
	if caller := metadata.String(ctx, metadata.Caller); caller != "" {
		f = append(f, KV(_caller, caller))
	}
	if color := metadata.String(ctx, metadata.Color); color != "" {
		f = append(f, KV(_color, color))
	}
	if cluster := metadata.String(ctx, metadata.Cluster); cluster != "" {
		f = append(f, KV(_cluster, cluster))
	}
	if metadata.Bool(ctx, metadata.Mirror) {
		f = append(f, KV(_mirror, true))
	}
	select {
	case h.msgs <- f:
	default:
	}
}

// writeproc write data into connection.
func (h *AgentHandler) writeproc() {
	var (
		conn  net.Conn
		err   error
		count int
		quit  bool
	)
	buf := core.NewBuffer(2048)

	defer h.waiter.Done()
	taskID := []byte(h.c.TaskID)
	tick := time.NewTicker(_mergeWait)
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
			buf.Write([]byte(strconv.FormatInt(now.UnixNano()/1e6, 10)))
			h.enc.Encode(buf, d...)
			h.free(d)
			if h.batchSend {
				buf.Write(_logSeparator)
				if count++; count < h.c.Buffer && buf.Len() < _maxBuffer {
					continue
				}
			}
		case <-tick.C:
		}
		if conn == nil || err != nil {
			if conn, err = net.DialTimeout(h.c.Proto, h.c.Addr, time.Duration(h.c.Timeout)); err != nil {
				stdlog.Printf("net.DialTimeout(%s:%s) error(%v)\n", h.c.Proto, h.c.Addr, err)
				continue
			}
		}
	DUMP:
		if conn != nil && buf.Len() > 0 {
			count = 0
			if _, err = conn.Write(buf.Bytes()); err != nil {
				stdlog.Printf("conn.Write(%d bytes) error(%v)\n", buf.Len(), err)
				conn.Close()
			} else {
				// only succeed reset buffer, let conn reconnect.
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

// Close close the connection.
func (h *AgentHandler) Close() (err error) {
	h.msgs <- nil
	h.waiter.Wait()
	return nil
}

// SetFormat .
func (h *AgentHandler) SetFormat(string) {
	// discard setformat
}
