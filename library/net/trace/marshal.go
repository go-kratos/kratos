package trace

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	protogen "go-common/library/net/trace/proto"
)

const protoVersion2 int32 = 2

func marshalSpan(sp *span, version int32) ([]byte, error) {
	if version == protoVersion2 {
		return marshalSpanV2(sp)
	}
	return marshalSpanV1(sp)
}

func marshalSpanV2(sp *span) ([]byte, error) {
	protoSpan := new(protogen.Span)
	protoSpan.Version = protoVersion2
	protoSpan.ServiceName = sp.dapper.serviceName
	protoSpan.OperationName = sp.operationName
	protoSpan.TraceId = sp.context.traceID
	protoSpan.SpanId = sp.context.spanID
	protoSpan.ParentId = sp.context.parentID
	protoSpan.SamplingProbability = sp.context.probability
	protoSpan.StartTime = &timestamp.Timestamp{
		Seconds: sp.startTime.Unix(),
		Nanos:   int32(sp.startTime.Nanosecond()),
	}
	protoSpan.Duration = &duration.Duration{
		Seconds: int64(sp.duration / time.Second),
		Nanos:   int32(sp.duration % time.Second),
	}
	protoSpan.Tags = make([]*protogen.Tag, len(sp.tags))
	for i := range sp.tags {
		protoSpan.Tags[i] = toProtoTag(sp.tags[i])
	}
	protoSpan.Logs = sp.logs
	return proto.Marshal(protoSpan)
}

func toLeagcyTag(tag Tag) *protogen.Tag {
	ptag := &protogen.Tag{Key: tag.Key}
	switch value := tag.Value.(type) {
	case string:
		ptag.Kind = protogen.Tag_STRING
		ptag.Value = []byte(value)
	case int:
		ptag.Kind = protogen.Tag_INT
		ptag.Value = []byte(strconv.FormatInt(int64(value), 10))
	case int32:
		ptag.Kind = protogen.Tag_INT
		ptag.Value = []byte(strconv.FormatInt(int64(value), 10))
	case int64:
		ptag.Kind = protogen.Tag_INT
		ptag.Value = []byte(strconv.FormatInt(value, 10))
	case bool:
		ptag.Kind = protogen.Tag_BOOL
		ptag.Value = []byte(strconv.FormatBool(value))
	case float32:
		ptag.Kind = protogen.Tag_BOOL
		ptag.Value = []byte(strconv.FormatFloat(float64(value), 'E', -1, 64))
	case float64:
		ptag.Kind = protogen.Tag_BOOL
		ptag.Value = []byte(strconv.FormatFloat(value, 'E', -1, 64))
	default:
		ptag.Kind = protogen.Tag_STRING
		ptag.Value = []byte((fmt.Sprintf("%v", tag.Value)))
	}
	return ptag
}

func toLeagcyLog(logs []*protogen.Log) []*protogen.Log {
	for _, log := range logs {
		if len(log.Fields) == 0 {
			continue
		}
		log.Key = log.Fields[0].Key
		log.Value = log.Fields[0].Value
	}
	return logs
}

func marshalSpanV1(sp *span) ([]byte, error) {
	protoSpan := new(protogen.Span)
	protoSpan.ServiceName = sp.dapper.serviceName
	protoSpan.OperationName = sp.operationName
	protoSpan.TraceId = sp.context.traceID
	protoSpan.SpanId = sp.context.spanID
	protoSpan.ParentId = sp.context.parentID
	protoSpan.SamplingProbability = sp.context.probability

	protoSpan.StartAt = sp.startTime.UnixNano()
	protoSpan.FinishAt = sp.startTime.UnixNano() + int64(sp.duration)

	protoSpan.Tags = make([]*protogen.Tag, len(sp.tags))
	for i := range sp.tags {
		protoSpan.Tags[i] = toLeagcyTag(sp.tags[i])
	}
	protoSpan.Logs = toLeagcyLog(sp.logs)
	return proto.Marshal(protoSpan)
}

func serializeInt64(v int64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(v))
	return data
}

func serializeFloat64(v float64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, math.Float64bits(v))
	return data
}

func serializeBool(v bool) []byte {
	data := make([]byte, 1)
	if v {
		data[0] = byte(1)
	} else {
		data[0] = byte(0)
	}
	return data
}

func toProtoTag(tag Tag) *protogen.Tag {
	ptag := &protogen.Tag{Key: tag.Key}
	switch value := tag.Value.(type) {
	case string:
		ptag.Kind = protogen.Tag_STRING
		ptag.Value = []byte(value)
	case int:
		ptag.Kind = protogen.Tag_INT
		ptag.Value = serializeInt64(int64(value))
	case int32:
		ptag.Kind = protogen.Tag_INT
		ptag.Value = serializeInt64(int64(value))
	case int64:
		ptag.Kind = protogen.Tag_INT
		ptag.Value = serializeInt64(value)
	case bool:
		ptag.Kind = protogen.Tag_BOOL
		ptag.Value = serializeBool(value)
	case float32:
		ptag.Kind = protogen.Tag_BOOL
		ptag.Value = serializeFloat64(float64(value))
	case float64:
		ptag.Kind = protogen.Tag_BOOL
		ptag.Value = serializeFloat64(value)
	default:
		ptag.Kind = protogen.Tag_STRING
		ptag.Value = []byte((fmt.Sprintf("%v", tag.Value)))
	}
	return ptag
}
