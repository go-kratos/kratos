package httpstream

import (
	"net/http"
	"sync"
	"regexp"
	"bytes"
	"strconv"
	"context"
	"go-common/app/service/ops/log-agent/event"
)

type HttpStream struct {
	c          *Config
	logstreams map[chan *event.ProcessorEvent]*filterRule
	l          sync.Mutex
}
type filterRule struct {
	maxLines   int
	appId      string
	instanceId string
	reg        *regexp.Regexp
}

var LogSourceChan = make(chan *event.ProcessorEvent)

// initlogStream init log stream
func NewHttpStream(config *Config) (httpStream *HttpStream, err error) {
	h := new(HttpStream)
	if err := config.ConfigValidate(); err != nil {
		return nil, err
	}
	h.c = config
	h.logstreams = make(map[chan *event.ProcessorEvent]*filterRule)
	http.HandleFunc("/logs", h.LogStreamer())
	go h.route()
	go http.ListenAndServe(h.c.Addr, nil)
	return h, nil
}

// route 把日志路由到所有注册的logstream
func (s *HttpStream) route() {
	for buf := range LogSourceChan {
		for logstream, _ := range s.logstreams {
			logstream <- buf
		}
	}
}

// LogStreamer 接收请求
func (s *HttpStream) LogStreamer() func(w http.ResponseWriter, req *http.Request) {
	logsHandler := func(w http.ResponseWriter, req *http.Request) {
		logstream := make(chan *event.ProcessorEvent)
		f := new(filterRule)
		// parse params
		params := req.URL.Query()
		if appId, ok := params["app_id"]; ok {
			f.appId = appId[0]
		} else {
			w.Write(append([]byte("必须指定app_id"), '\n'))
			return
		}
		if reg, ok := params["regexp"]; ok {
			if filterReg, err := regexp.Compile(reg[0]); err == nil {
				f.reg = filterReg
			} else {
				w.Write(append([]byte("正则表达式格式错误"), '\n'))
				return
			}
		}
		if instanceId, ok := params["instance_id"]; ok {
			f.instanceId = instanceId[0]
		}

		if maxLines, ok := params["max_lines"]; ok {
			if n, err := strconv.Atoi(maxLines[0]); err == nil {
				f.maxLines = n
			} else {
				w.Write(append([]byte("max_lines格式错误"), '\n'))
				return
			}
		}

		s.add(logstream, f)
		go func() {
			select {
			case <-req.Context().Done():
				s.remove(logstream)
			}
		}()
		defer s.httpStreamer(req.Context(), w, req, logstream, f)
	}
	return logsHandler
}

// httpStreamer 过滤并输出日志
func (s *HttpStream) httpStreamer(ctx context.Context, w http.ResponseWriter, req *http.Request, logstream chan *event.ProcessorEvent, f *filterRule) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if f.instanceId != "" {
		//w.Write([]byte(fmt.Sprintf(
		//	"\x1b[1;31m%s\x1b[0m\n", "注意：caster中，只有(1)app_id满足服务树三级格式 (2)日志包含instance_id且值为实例名称 的情况下，日志才能输出\n")))
		w.Write([]byte("注意：caster中，只有(1)app_id满足服务树三级格式 (2)日志包含instance_id且值为实例名称 的情况下，日志才能输出                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     \n"))
		w.(http.Flusher).Flush()
	}
	c := 0
	for {
		select {
		case e := <-logstream:
			if f.appId != "" && string(e.AppId) != f.appId {
				continue
			}
			if f.reg != nil && !f.reg.Match(e.Bytes()) {
				continue
			}
			if f.instanceId != "" && !bytes.Contains(e.Bytes(), []byte(f.instanceId)) {
				continue
			}
			if f.maxLines != 0 && c >= f.maxLines {
				s.remove(logstream)
				return
			}
			c += 1
			// TODO event recycle
			w.Write(append(e.Bytes(), '\n'))
			w.(http.Flusher).Flush()
		case <-ctx.Done():
			return
		}
	}
}

// add 注册logstream
func (s *HttpStream) add(logstream chan *event.ProcessorEvent, f *filterRule) {
	s.l.Lock()
	defer s.l.Unlock()
	s.logstreams[logstream] = f
}

// remove 注销logstream
func (s *HttpStream) remove(logstream chan *event.ProcessorEvent) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.logstreams, logstream)
}
