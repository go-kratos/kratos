package trace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bilibili/kratos/pkg/net/trace/jaegerutils"
	"github.com/bilibili/kratos/pkg/net/trace/jaegerutils/thrift"
	"github.com/bilibili/kratos/pkg/net/trace/jaegerutils/thrift-gen/jaeger"
	protogen "github.com/bilibili/kratos/pkg/net/trace/proto"
)

const (
	// DefaultUDPSpanServerHost is the default host to send the spans to, via UDP
	DefaultUDPSpanServerHost = "localhost"

	// DefaultUDPSpanServerPort is the default port to send the spans to, via UDP
	DefaultUDPSpanServerPort = 6831

	// UDPPacketMaxLength is the max size of UDP packet we want to send, synced with jaeger-agent
	UDPPacketMaxLength = 65000

	emitBatchOverhead        = 30
	defaultMaxTagValueLength = 1024

	defaultHTTPTimeout   = time.Second * 5
	defaultHTTPBatchSize = 100
)

var errSpanTooLarge = errors.New("Span is too large")

func stringify(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprintf("%+v", value)
}

// TimeToMicrosecondsSinceEpochInt64 converts Go time.Time to a long
// representing time since epoch in microseconds, which is used expected
// in the Jaeger spans encoded as Thrift.
func TimeToMicrosecondsSinceEpochInt64(t time.Time) int64 {
	// ^^^ Passing time.Time by value is faster than passing a pointer!
	// BenchmarkTimeByValue-8	2000000000	         1.37 ns/op
	// BenchmarkTimeByPtr-8  	2000000000	         1.98 ns/op

	return t.UnixNano() / 1000
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

func buildTag(tag *Tag, maxTagValueLength int) *jaeger.Tag {
	jTag := &jaeger.Tag{Key: tag.Key}
	switch value := tag.Value.(type) {
	case string:
		vStr := truncateString(value, maxTagValueLength)
		jTag.VStr = &vStr
		jTag.VType = jaeger.TagType_STRING
	case []byte:
		if len(value) > maxTagValueLength {
			value = value[:maxTagValueLength]
		}
		jTag.VBinary = value
		jTag.VType = jaeger.TagType_BINARY
	case int:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case uint:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case int8:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case uint8:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case int16:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case uint16:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case int32:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case uint32:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case int64:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case uint64:
		vLong := int64(value)
		jTag.VLong = &vLong
		jTag.VType = jaeger.TagType_LONG
	case float32:
		vDouble := float64(value)
		jTag.VDouble = &vDouble
		jTag.VType = jaeger.TagType_DOUBLE
	case float64:
		vDouble := float64(value)
		jTag.VDouble = &vDouble
		jTag.VType = jaeger.TagType_DOUBLE
	case bool:
		vBool := value
		jTag.VBool = &vBool
		jTag.VType = jaeger.TagType_BOOL
	default:
		vStr := truncateString(stringify(value), maxTagValueLength)
		jTag.VStr = &vStr
		jTag.VType = jaeger.TagType_STRING
	}
	return jTag
}

func buildTags(tags []Tag, maxTagValueLength int) []*jaeger.Tag {
	jTags := make([]*jaeger.Tag, 0, len(tags))
	for _, tag := range tags {
		jTag := buildTag(&tag, maxTagValueLength)
		jTags = append(jTags, jTag)
	}
	return jTags
}

// ConvertLogsToJaegerTags converts log Fields into jaeger tags.
func ConvertLogsToJaegerTags(logFields []*protogen.Field) []*jaeger.Tag {
	fields := make([]*jaeger.Tag, 0, len(logFields))
	for _, field := range logFields {
		vStr := string(field.Value)
		fields = append(fields, &jaeger.Tag{Key: field.Key, VStr: &vStr, VType: jaeger.TagType_STRING})
	}
	return fields
}

func buildLogs(logs []*protogen.Log) []*jaeger.Log {
	jLogs := make([]*jaeger.Log, 0, len(logs))
	for _, log := range logs {
		jLog := &jaeger.Log{
			Timestamp: log.Timestamp / 1000,
			Fields:    ConvertLogsToJaegerTags(log.Fields),
		}
		jLogs = append(jLogs, jLog)
	}
	return jLogs
}

// BuildJaegerThrift builds jaeger span based on internal span.
func BuildJaegerThrift(sp *Span) *jaeger.Span {
	duration := sp.duration.Nanoseconds() / int64(time.Microsecond)
	startTime := TimeToMicrosecondsSinceEpochInt64(sp.startTime)
	jaegerSpan := &jaeger.Span{
		TraceIdLow:    int64(sp.context.TraceID),
		TraceIdHigh:   0,
		SpanId:        int64(sp.context.SpanID),
		ParentSpanId:  int64(sp.context.ParentID),
		OperationName: sp.operationName,
		Flags:         int32(sp.context.Flags),
		StartTime:     startTime,
		Duration:      duration,
		Tags:          buildTags(sp.tags, defaultMaxTagValueLength),
		Logs:          buildLogs(sp.logs),
	}
	if sp.context.ParentID != 0 {
		jaegerSpan.References = []*jaeger.SpanRef{&jaeger.SpanRef{
			RefType:     jaeger.SpanRefType_CHILD_OF,
			TraceIdLow:  int64(sp.context.ParentID),
			TraceIdHigh: 0,
			SpanId:      int64(sp.context.SpanID),
		}}
	}
	return jaegerSpan
}

func buildJaegerProcessThrift(dp *dapper) *jaeger.Process {
	process := &jaeger.Process{
		ServiceName: dp.serviceName,
		Tags:        buildTags(dp.tags, defaultMaxTagValueLength),
	}
	return process
}

type jaegerUDPReport struct {
	mx              sync.Mutex
	client          *jaegerutils.AgentClientUDP
	process         *jaeger.Process
	maxPacketSize   int                   // max size of datagram in bytes
	maxSpanBytes    int                   // max number of bytes to record spans (excluding envelope) in the datagram
	byteBufferSize  int                   // current number of span bytes accumulated in the buffer
	spanBuffer      []*jaeger.Span        // spans buffered before a flush
	thriftBuffer    *thrift.TMemoryBuffer // buffer used to calculate byte size of a span
	thriftProtocol  thrift.TProtocol
	processByteSize int
}

func (j *jaegerUDPReport) calcSizeOfSerializedThrift(thriftStruct thrift.TStruct) int {
	j.thriftBuffer.Reset()
	thriftStruct.Write(j.thriftProtocol)
	return j.thriftBuffer.Len()
}

// NewJaegerUDPReport report trace to jaeger use udp.
func NewJaegerUDPReport(hostPort string, maxPacketSize int) (reporter, error) {
	if len(hostPort) == 0 {
		hostPort = fmt.Sprintf("%s:%d", DefaultUDPSpanServerHost, DefaultUDPSpanServerPort)
	}
	if maxPacketSize == 0 {
		maxPacketSize = UDPPacketMaxLength
	}
	protocolFactory := thrift.NewTCompactProtocolFactory()
	// Each span is first written to thriftBuffer to determine its size in bytes.
	thriftBuffer := thrift.NewTMemoryBufferLen(maxPacketSize)
	thriftProtocol := protocolFactory.GetProtocol(thriftBuffer)

	client, err := jaegerutils.NewAgentClientUDP(hostPort, maxPacketSize)
	if err != nil {
		return nil, err
	}
	report := &jaegerUDPReport{
		client:         client,
		maxSpanBytes:   maxPacketSize - emitBatchOverhead,
		thriftBuffer:   thriftBuffer,
		thriftProtocol: thriftProtocol}
	return report, nil
}

func (j *jaegerUDPReport) resetBuffers() {
	for i := range j.spanBuffer {
		j.spanBuffer[i] = nil
	}
	j.spanBuffer = j.spanBuffer[:0]
	j.byteBufferSize = j.processByteSize
}

func (j *jaegerUDPReport) Flush() error {
	n := len(j.spanBuffer)
	if n == 0 {
		return nil
	}
	err := j.client.EmitBatch(&jaeger.Batch{Process: j.process, Spans: j.spanBuffer})
	j.resetBuffers()
	return err
}

func (j *jaegerUDPReport) WriteSpan(sp *Span) error {
	j.mx.Lock()
	defer j.mx.Unlock()
	if j.process == nil {
		j.process = buildJaegerProcessThrift(sp.dapper)
		j.processByteSize = j.calcSizeOfSerializedThrift(j.process)
		j.byteBufferSize += j.processByteSize
	}
	jSpan := BuildJaegerThrift(sp)
	spanSize := j.calcSizeOfSerializedThrift(jSpan)
	if spanSize > j.maxSpanBytes {
		return errSpanTooLarge
	}

	j.byteBufferSize += spanSize
	if j.byteBufferSize <= j.maxSpanBytes {
		j.spanBuffer = append(j.spanBuffer, jSpan)
		if j.byteBufferSize < j.maxSpanBytes {
			return nil
		}
		return j.Flush()
	}
	// the latest span did not fit in the buffer
	err := j.Flush()
	j.spanBuffer = append(j.spanBuffer, jSpan)
	j.byteBufferSize = spanSize + j.processByteSize
	return err
}

func (j *jaegerUDPReport) Close() error {
	j.Flush()
	return j.client.Close()
}

// NewJaegerHTTPReport report trace to jaeger use http protocol.
func NewJaegerHTTPReport(entrypoint string, batchSize int) (reporter, error) {
	// TODO: support multi entrypoint and custom path.
	if !strings.HasPrefix(entrypoint, "http://") {
		entrypoint = fmt.Sprintf("http://%s/api/traces", entrypoint)
	}
	if batchSize == 0 {
		batchSize = defaultHTTPBatchSize
	}
	httpReport := &jaegerHTTPReport{
		entrypoint: entrypoint,
		client:     &http.Client{Timeout: defaultHTTPTimeout},
		batchSize:  defaultHTTPBatchSize,
		spans:      make([]*jaeger.Span, 0, defaultHTTPBatchSize),
		queue:      make(chan *bytes.Buffer, 1),
	}
	httpReport.wg.Add(1)
	go httpReport.daemon()
	return httpReport, nil
}

type jaegerHTTPReport struct {
	mx         sync.Mutex
	wg         sync.WaitGroup
	entrypoint string
	client     *http.Client
	batchSize  int
	process    *jaeger.Process
	spans      []*jaeger.Span
	queue      chan *bytes.Buffer
}

func (j *jaegerHTTPReport) WriteSpan(sp *Span) error {
	j.mx.Lock()
	defer j.mx.Unlock()
	if j.process == nil {
		j.process = buildJaegerProcessThrift(sp.dapper)
	}
	jSpan := BuildJaegerThrift(sp)
	j.spans = append(j.spans, jSpan)
	if len(j.spans) > j.batchSize {
		return j.flush()
	}
	return nil
}

func (j *jaegerHTTPReport) flush() error {
	batch := jaeger.Batch{
		Process: j.process,
		Spans:   j.spans,
	}
	batchSize := len(j.spans)
	body, err := serializeThrift(&batch)
	// reset spans don't care serialize error.
	j.spans = j.spans[:0]
	if err != nil {
		return err
	}
	select {
	case j.queue <- body:
	default:
		jaegerReportDroppedCounter.Add(float64(batchSize), "jaeger_http_report")
	}
	return nil
}

func (j *jaegerHTTPReport) daemon() {
	for body := range j.queue {
		if err := j.send(body); err != nil {
			jaegerReportErrorCounter.Add(1, "jaeger_http_report")
			_sdkLogger.Printf("WARN report trace data error %s, you can ignore this error.", err)
		}
	}
	j.wg.Done()
}

func (j *jaegerHTTPReport) send(body *bytes.Buffer) error {
	req, err := http.NewRequest(http.MethodPost, j.entrypoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-thrift")

	resp, err := j.client.Do(req)
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

func (j *jaegerHTTPReport) Close() error {
	j.flush()
	close(j.queue)
	j.wg.Wait()
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
