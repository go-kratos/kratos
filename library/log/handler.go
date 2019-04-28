package log

import (
	"context"
	"time"

	pkgerr "github.com/pkg/errors"
)

const (
	_timeFormat = "2006-01-02T15:04:05.999999"

	// log level defined in level.go.
	_levelValue = "level_value"
	//  log level name: INFO, WARN...
	_level = "level"
	// log time.
	_time = "time"
	// request path.
	// _title = "title"
	// log file.
	_source = "source"
	// common log filed.
	_log = "log"
	// app name.
	_appID = "app_id"
	// container ID.
	_instanceID = "instance_id"
	// uniq ID from trace.
	_tid = "traceid"
	// request time.
	// _ts = "ts"
	// requester.
	_caller = "caller"
	// container environment: prod, pre, uat, fat.
	_deplyEnv = "env"
	// container area.
	_zone = "zone"
	// mirror flag
	_mirror = "mirror"
	// color.
	_color = "color"
	// cluster.
	_cluster = "cluster"
)

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	// Log handle log
	// variadic D is k-v struct represent log content
	Log(context.Context, Level, ...D)

	// SetFormat set render format on log output
	// see StdoutHandler.SetFormat for detail
	SetFormat(string)

	// Close handler
	Close() error
}

func newHandlers(filters []string, handlers ...Handler) *Handlers {
	set := make(map[string]struct{})
	for _, k := range filters {
		set[k] = struct{}{}
	}
	return &Handlers{filters: set, handlers: handlers}
}

// Handlers a bundle for hander with filter function.
type Handlers struct {
	filters  map[string]struct{}
	handlers []Handler
}

// Log handlers logging.
func (hs Handlers) Log(c context.Context, lv Level, d ...D) {
	var fn string
	for i := range d {
		if _, ok := hs.filters[d[i].Key]; ok {
			d[i].Value = "***"
		}
		if d[i].Key == _source {
			fn = d[i].Value.(string)
		}
	}
	if fn == "" {
		d = append(d, KV(_source, funcName(4)))
	}
	d = append(d, KV(_time, time.Now()), KV(_levelValue, int(lv)), KV(_level, lv.String()))
	errIncr(lv, fn)
	for _, h := range hs.handlers {
		h.Log(c, lv, d...)
	}
}

// Close close resource.
func (hs Handlers) Close() (err error) {
	for _, h := range hs.handlers {
		if e := h.Close(); e != nil {
			err = pkgerr.WithStack(e)
		}
	}
	return
}

// SetFormat .
func (hs Handlers) SetFormat(format string) {
	for _, h := range hs.handlers {
		h.SetFormat(format)
	}
}
