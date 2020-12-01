package jaeger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	ja "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/thrift"
	j "github.com/uber/jaeger-client-go/thrift-gen/jaeger"
)

// Default timeout for http request in seconds
const defaultHTTPTimeout = time.Second * 5

// HTTPTransport implements Transport by forwarding spans to a http server.
type HTTPTransport struct {
	url             string
	client          *http.Client
	batchSize       int
	spans           []*j.Span
	process         *j.Process
	httpCredentials *HTTPBasicAuthCredentials
	headers         map[string]string
}

// HTTPBasicAuthCredentials stores credentials for HTTP basic auth.
type HTTPBasicAuthCredentials struct {
	username string
	password string
}

// HTTPOption sets a parameter for the HttpCollector
type HTTPOption func(c *HTTPTransport)

// HTTPTimeout sets maximum timeout for http request.
func HTTPTimeout(duration time.Duration) HTTPOption {
	return func(c *HTTPTransport) { c.client.Timeout = duration }
}

// HTTPBatchSize sets the maximum batch size, after which a collect will be
// triggered. The default batch size is 100 spans.
func HTTPBatchSize(n int) HTTPOption {
	return func(c *HTTPTransport) { c.batchSize = n }
}

// HTTPBasicAuth sets the credentials required to perform HTTP basic auth
func HTTPBasicAuth(username string, password string) HTTPOption {
	return func(c *HTTPTransport) {
		c.httpCredentials = &HTTPBasicAuthCredentials{username: username, password: password}
	}
}

// HTTPRoundTripper configures the underlying Transport on the *http.Client
// that is used
func HTTPRoundTripper(transport http.RoundTripper) HTTPOption {
	return func(c *HTTPTransport) {
		c.client.Transport = transport
	}
}

// HTTPHeaders defines the HTTP headers that will be attached to the jaeger client's HTTP request
func HTTPHeaders(headers map[string]string) HTTPOption {
	return func(c *HTTPTransport) {
		c.headers = headers
	}
}

// NewHTTPTransport returns a new HTTP-backend transport. url should be an http
// url of the collector to handle POST request, typically something like:
//     http://hostname:14268/api/traces?format=jaeger.thrift
func NewHTTPTransport(url string, options ...HTTPOption) *HTTPTransport {
	c := &HTTPTransport{
		url:       url,
		client:    &http.Client{Timeout: defaultHTTPTimeout},
		batchSize: 100,
		spans:     []*j.Span{},
	}

	for _, option := range options {
		option(c)
	}
	return c
}

// Append implements Transport.
func (c *HTTPTransport) Append(span *Span) (int, error) {
	if c.process == nil {
		process := j.NewProcess()
		process.ServiceName = span.ServiceName()
		c.process = process
	}
	jSpan := BuildJaegerThrift(span)
	c.spans = append(c.spans, jSpan)
	if len(c.spans) >= c.batchSize {
		return c.Flush()
	}
	return 0, nil
}

// Flush implements Transport.
func (c *HTTPTransport) Flush() (int, error) {
	count := len(c.spans)
	if count == 0 {
		return 0, nil
	}
	err := c.send(c.spans)
	c.spans = c.spans[:0]
	return count, err
}

// Close implements Transport.
func (c *HTTPTransport) Close() error {
	return nil
}

func (c *HTTPTransport) send(spans []*j.Span) error {
	batch := &j.Batch{
		Spans:   spans,
		Process: c.process,
	}
	body, err := serializeThrift(batch)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-thrift")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	if c.httpCredentials != nil {
		req.SetBasicAuth(c.httpCredentials.username, c.httpCredentials.password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("error from collector: %d", resp.StatusCode)
	}
	return nil
}

func serializeThrift(obj thrift.TStruct) (*bytes.Buffer, error) {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolTransport(t)
	if err := obj.Write(p); err != nil {
		return nil, err
	}
	return t.Buffer, nil
}

func BuildJaegerThrift(span *Span) *j.Span {
	span.Lock()
	defer span.Unlock()
	startTime := span.startTime.UnixNano() / 1000
	duration := span.duration.Nanoseconds() / int64(time.Microsecond)
	jaegerSpan := &j.Span{
		TraceIdLow:    int64(span.context.traceID.Low),
		TraceIdHigh:   int64(span.context.traceID.High),
		SpanId:        int64(span.context.spanID),
		ParentSpanId:  int64(span.context.parentID),
		OperationName: span.operationName,
		Flags:         int32(span.context.samplingState.flags()),
		StartTime:     startTime,
		Duration:      duration,
		Tags:          buildTags(span.tags, 100),
		Logs:          buildLogs(span.logs),
		References:    buildReferences(span.references),
	}
	return jaegerSpan
}

func stringify(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprintf("%+v", value)
}

func truncateString(value string, maxLength int) string {
	// we ignore the problem of utf8 runes possibly being sliced in the middle,
	// as it is rather expensive to iterate through each tag just to find rune
	// boundaries.
	if len(value) > maxLength {
		return value[:maxLength]
	}
	return value
}

func buildTags(tags []Tag, maxTagValueLength int) []*j.Tag {
	jTags := make([]*j.Tag, 0, len(tags))
	for _, tag := range tags {
		jTag := buildTag(&tag, maxTagValueLength)
		jTags = append(jTags, jTag)
	}
	return jTags
}
func buildTag(tag *Tag, maxTagValueLength int) *j.Tag {
	jTag := &j.Tag{Key: tag.key}
	switch value := tag.value.(type) {
	case string:
		vStr := truncateString(value, maxTagValueLength)
		jTag.VStr = &vStr
		jTag.VType = j.TagType_STRING
	case []byte:
		if len(value) > maxTagValueLength {
			value = value[:maxTagValueLength]
		}
		jTag.VBinary = value
		jTag.VType = j.TagType_BINARY
	case int:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case uint:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case int8:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case uint8:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case int16:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case uint16:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case int32:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case uint32:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case int64:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case uint64:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = j.TagType_LONG
	case float32:
		vDouble := float64(value)
		jTag.VDouble = &vDouble
		jTag.VType = j.TagType_DOUBLE
	case float64:
		vDouble := float64(value)
		jTag.VDouble = &vDouble
		jTag.VType = j.TagType_DOUBLE
	case bool:
		vBool := value
		jTag.VBool = &vBool
		jTag.VType = j.TagType_BOOL
	default:
		vStr := truncateString(stringify(value), maxTagValueLength)
		jTag.VStr = &vStr
		jTag.VType = j.TagType_STRING
	}
	return jTag
}

func buildLogs(logs []opentracing.LogRecord) []*j.Log {
	jLogs := make([]*j.Log, 0, len(logs))
	for _, log := range logs {
		jLog := &j.Log{
			Timestamp: log.Timestamp.UnixNano() / 1000,
			Fields:    ja.ConvertLogsToJaegerTags(log.Fields),
		}
		jLogs = append(jLogs, jLog)
	}
	return jLogs
}

func buildReferences(references []Reference) []*j.SpanRef {
	retMe := make([]*j.SpanRef, 0, len(references))
	for _, ref := range references {
		if ref.Type == opentracing.ChildOfRef {
			retMe = append(retMe, spanRef(ref.Context, j.SpanRefType_CHILD_OF))
		} else if ref.Type == opentracing.FollowsFromRef {
			retMe = append(retMe, spanRef(ref.Context, j.SpanRefType_FOLLOWS_FROM))
		}
	}
	return retMe
}

func spanRef(ctx SpanContext, refType j.SpanRefType) *j.SpanRef {
	return &j.SpanRef{
		RefType:     refType,
		TraceIdLow:  int64(ctx.traceID.Low),
		TraceIdHigh: int64(ctx.traceID.High),
		SpanId:      int64(ctx.spanID),
	}
}
