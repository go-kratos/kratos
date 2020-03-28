package trace

import (
	"encoding/binary"
	errs "errors"
	"fmt"
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	protogen "github.com/go-kratos/kratos/pkg/net/trace/proto"
)

const protoVersion1 int32 = 1

var (
	errSpanVersion = errs.New("trace: marshal not support version")
)

func marshalSpan(sp *Span, version int32) ([]byte, error) {
	if version == protoVersion1 {
		return marshalSpanV1(sp)
	}
	return nil, errSpanVersion
}

func marshalSpanV1(sp *Span) ([]byte, error) {
	protoSpan := new(protogen.Span)
	protoSpan.Version = protoVersion1
	protoSpan.ServiceName = sp.dapper.serviceName
	protoSpan.OperationName = sp.operationName
	protoSpan.TraceId = sp.context.TraceID
	protoSpan.SpanId = sp.context.SpanID
	protoSpan.ParentId = sp.context.ParentID
	protoSpan.SamplingProbability = sp.context.Probability
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
