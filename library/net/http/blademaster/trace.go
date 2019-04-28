package blademaster

import (
	"io"
	"net/http"
	"net/http/httptrace"
	"strconv"

	"go-common/library/net/metadata"
	"go-common/library/net/trace"
)

const _defaultComponentName = "net/http"

// Trace is trace middleware
func Trace() HandlerFunc {
	return func(c *Context) {
		// handle http request
		// get derived trace from http request header
		t, err := trace.Extract(trace.HTTPFormat, c.Request.Header)
		if err != nil {
			var opts []trace.Option
			if ok, _ := strconv.ParseBool(trace.BiliTraceDebug); ok {
				opts = append(opts, trace.EnableDebug())
			}
			t = trace.New(c.Request.URL.Path, opts...)
		}
		t.SetTitle(c.Request.URL.Path)
		t.SetTag(trace.String(trace.TagComponent, _defaultComponentName))
		t.SetTag(trace.String(trace.TagHTTPMethod, c.Request.Method))
		t.SetTag(trace.String(trace.TagHTTPURL, c.Request.URL.String()))
		t.SetTag(trace.String(trace.TagSpanKind, "server"))
		t.SetTag(trace.String("caller", metadata.String(c.Context, metadata.Caller)))
		c.Context = trace.NewContext(c.Context, t)
		c.Next()
		t.Finish(&c.Error)
	}
}

type closeTracker struct {
	io.ReadCloser
	tr trace.Trace
}

func (c *closeTracker) Close() error {
	err := c.ReadCloser.Close()
	c.tr.SetLog(trace.Log(trace.LogEvent, "ClosedBody"))
	c.tr.Finish(&err)
	return err
}

// NewTraceTracesport NewTraceTracesport
func NewTraceTracesport(rt http.RoundTripper, peerService string, internalTags ...trace.Tag) *TraceTransport {
	return &TraceTransport{RoundTripper: rt, peerService: peerService, internalTags: internalTags}
}

// TraceTransport wraps a RoundTripper. If a request is being traced with
// Tracer, Transport will inject the current span into the headers,
// and set HTTP related tags on the span.
type TraceTransport struct {
	peerService  string
	internalTags []trace.Tag
	// The actual RoundTripper to use for the request. A nil
	// RoundTripper defaults to http.DefaultTransport.
	http.RoundTripper
}

// RoundTrip implements the RoundTripper interface
func (t *TraceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}
	tr, ok := trace.FromContext(req.Context())
	if !ok {
		return rt.RoundTrip(req)
	}
	operationName := "HTTP:" + req.Method
	// fork new trace
	tr = tr.Fork("", operationName)

	tr.SetTag(trace.TagString(trace.TagComponent, _defaultComponentName))
	tr.SetTag(trace.TagString(trace.TagHTTPMethod, req.Method))
	tr.SetTag(trace.TagString(trace.TagHTTPURL, req.URL.String()))
	tr.SetTag(trace.TagString(trace.TagSpanKind, "client"))
	if t.peerService != "" {
		tr.SetTag(trace.TagString(trace.TagPeerService, t.peerService))
	}
	tr.SetTag(t.internalTags...)

	// inject trace to http header
	trace.Inject(tr, trace.HTTPFormat, req.Header)

	// FIXME: uncomment after trace sdk is goroutinue safe
	// ct := clientTracer{tr: tr}
	// req = req.WithContext(httptrace.WithClientTrace(req.Context(), ct.clientTrace()))
	resp, err := rt.RoundTrip(req)

	if err != nil {
		tr.SetTag(trace.TagBool(trace.TagError, true))
		tr.Finish(&err)
		return resp, err
	}

	// TODO: get ecode
	tr.SetTag(trace.TagInt64(trace.TagHTTPStatusCode, int64(resp.StatusCode)))

	if req.Method == "HEAD" {
		tr.Finish(nil)
	} else {
		resp.Body = &closeTracker{resp.Body, tr}
	}
	return resp, err
}

type clientTracer struct {
	tr trace.Trace
}

func (h *clientTracer) clientTrace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn:              h.getConn,
		GotConn:              h.gotConn,
		PutIdleConn:          h.putIdleConn,
		GotFirstResponseByte: h.gotFirstResponseByte,
		Got100Continue:       h.got100Continue,
		DNSStart:             h.dnsStart,
		DNSDone:              h.dnsDone,
		ConnectStart:         h.connectStart,
		ConnectDone:          h.connectDone,
		WroteHeaders:         h.wroteHeaders,
		Wait100Continue:      h.wait100Continue,
		WroteRequest:         h.wroteRequest,
	}
}

func (h *clientTracer) getConn(hostPort string) {
	// ext.HTTPUrl.Set(h.sp, hostPort)
	h.tr.SetLog(trace.Log(trace.LogEvent, "GetConn"))
}

func (h *clientTracer) gotConn(info httptrace.GotConnInfo) {
	h.tr.SetTag(trace.TagBool("net/http.reused", info.Reused))
	h.tr.SetTag(trace.TagBool("net/http.was_idle", info.WasIdle))
	h.tr.SetLog(trace.Log(trace.LogEvent, "GotConn"))
}

func (h *clientTracer) putIdleConn(error) {
	h.tr.SetLog(trace.Log(trace.LogEvent, "PutIdleConn"))
}

func (h *clientTracer) gotFirstResponseByte() {
	h.tr.SetLog(trace.Log(trace.LogEvent, "GotFirstResponseByte"))
}

func (h *clientTracer) got100Continue() {
	h.tr.SetLog(trace.Log(trace.LogEvent, "Got100Continue"))
}

func (h *clientTracer) dnsStart(info httptrace.DNSStartInfo) {
	h.tr.SetLog(
		trace.Log(trace.LogEvent, "DNSStart"),
		trace.Log("host", info.Host),
	)
}

func (h *clientTracer) dnsDone(info httptrace.DNSDoneInfo) {
	fields := []trace.LogField{trace.Log(trace.LogEvent, "DNSDone")}
	for _, addr := range info.Addrs {
		fields = append(fields, trace.Log("addr", addr.String()))
	}
	if info.Err != nil {
		// TODO: support log error object
		fields = append(fields, trace.Log(trace.LogErrorObject, info.Err.Error()))
	}
	h.tr.SetLog(fields...)
}

func (h *clientTracer) connectStart(network, addr string) {
	h.tr.SetLog(
		trace.Log(trace.LogEvent, "ConnectStart"),
		trace.Log("network", network),
		trace.Log("addr", addr),
	)
}

func (h *clientTracer) connectDone(network, addr string, err error) {
	if err != nil {
		h.tr.SetLog(
			trace.Log("message", "ConnectDone"),
			trace.Log("network", network),
			trace.Log("addr", addr),
			trace.Log(trace.LogEvent, "error"),
			// TODO: support log error object
			trace.Log(trace.LogErrorObject, err.Error()),
		)
	} else {
		h.tr.SetLog(
			trace.Log(trace.LogEvent, "ConnectDone"),
			trace.Log("network", network),
			trace.Log("addr", addr),
		)
	}
}

func (h *clientTracer) wroteHeaders() {
	h.tr.SetLog(trace.Log("event", "WroteHeaders"))
}

func (h *clientTracer) wait100Continue() {
	h.tr.SetLog(trace.Log("event", "Wait100Continue"))
}

func (h *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		h.tr.SetLog(
			trace.Log("message", "WroteRequest"),
			trace.Log("event", "error"),
			// TODO: support log error object
			trace.Log(trace.LogErrorObject, info.Err.Error()),
		)
		h.tr.SetTag(trace.TagBool(trace.TagError, true))
	} else {
		h.tr.SetLog(trace.Log("event", "WroteRequest"))
	}
}
