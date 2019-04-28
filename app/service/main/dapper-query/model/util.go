package model

import (
	"encoding/binary"
	"math"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	protogen "go-common/library/net/trace/proto"
)

const protoVersion2 int32 = 2

// FromProtoSpan convert protogen.Span to model.Span
func FromProtoSpan(data []byte) (*Span, error) {
	protoSpan := new(protogen.Span)
	if err := proto.Unmarshal(data, protoSpan); err != nil {
		return nil, err
	}
	var span *Span
	var err error
	if protoSpan.Version != protoVersion2 {
		span, err = fromProtoSpanLeagcy(protoSpan)
	} else {
		span, err = fromProtoSpanInternal(protoSpan)
	}
	return span, err
}

func convertLeagcyTag(protoTag *protogen.Tag) Tag {
	tag := Tag{Key: protoTag.Key}
	switch protoTag.Kind {
	case protogen.Tag_STRING:
		tag.Kind = TagString
		tag.Value = string(protoTag.Value)
	case protogen.Tag_INT:
		tag.Kind = TagInt
		tag.Value, _ = strconv.ParseInt(string(protoTag.Value), 10, 64)
	case protogen.Tag_BOOL:
		tag.Kind = TagBool
		tag.Value, _ = strconv.ParseBool(string(protoTag.Value))
	case protogen.Tag_FLOAT:
		tag.Kind = TagFloat
		tag.Value, _ = strconv.ParseFloat(string(protoTag.Value), 64)
	}
	return tag
}

func convertLeagcyLog(protoLog *protogen.Log) Log {
	log := Log{Timestamp: protoLog.Timestamp}
	log.Fields = []Field{{Key: protoLog.Key, Value: protoLog.Value}}
	return log
}

func fromProtoSpanLeagcy(protoSpan *protogen.Span) (*Span, error) {
	span := &Span{
		ServiceName:   protoSpan.ServiceName,
		OperationName: protoSpan.OperationName,
		TraceID:       protoSpan.TraceId,
		SpanID:        protoSpan.SpanId,
		ParentID:      protoSpan.ParentId,
		Env:           protoSpan.Env,
	}
	span.StartTime = time.Unix(protoSpan.StartAt/int64(time.Second), protoSpan.StartAt%int64(time.Second))
	span.Duration = time.Duration(protoSpan.FinishAt - protoSpan.StartAt)
	span.References = []SpanRef{{
		RefType: RefTypeChildOf,
		TraceID: protoSpan.TraceId,
		SpanID:  protoSpan.ParentId,
	}}
	span.Tags = make(map[string]interface{})
	for _, tag := range protoSpan.Tags {
		newTag := convertLeagcyTag(tag)
		span.Tags[newTag.Key] = newTag.Value
	}
	span.Logs = make([]Log, 0, len(protoSpan.Logs))
	for _, log := range protoSpan.Logs {
		span.Logs = append(span.Logs, convertLeagcyLog(log))
	}
	return span, nil
}

func timeFromTimestamp(t *timestamp.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}

func durationFromDuration(d *duration.Duration) time.Duration {
	return time.Duration(d.Seconds*int64(time.Second) + int64(d.Nanos))
}

func convertSpanRef(protoRef *protogen.SpanRef) SpanRef {
	ref := SpanRef{
		TraceID: protoRef.TraceId,
		SpanID:  protoRef.SpanId,
	}
	switch protoRef.RefType {
	case protogen.SpanRef_CHILD_OF:
		ref.RefType = RefTypeChildOf
	case protogen.SpanRef_FOLLOWS_FROM:
		ref.RefType = RefTypeFollowsFrom
	}
	return ref
}

func unSerializeInt64(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
}

func unSerializeBool(data []byte) bool {
	return data[0] == byte(1)
}

func unSerializeFloat64(data []byte) float64 {
	value := binary.BigEndian.Uint64(data)
	return math.Float64frombits(value)
}

func convertTag(protoTag *protogen.Tag) Tag {
	tag := Tag{Key: protoTag.Key}
	switch protoTag.Kind {
	case protogen.Tag_STRING:
		tag.Kind = TagString
		tag.Value = string(protoTag.Value)
	case protogen.Tag_INT:
		tag.Kind = TagInt
		tag.Value = unSerializeInt64(protoTag.Value)
	case protogen.Tag_BOOL:
		tag.Kind = TagBool
		tag.Value = unSerializeBool(protoTag.Value)
	case protogen.Tag_FLOAT:
		tag.Kind = TagFloat
		tag.Value = unSerializeFloat64(protoTag.Value)
	}
	return tag
}

func convertLog(protoLog *protogen.Log) Log {
	log := Log{Timestamp: protoLog.Timestamp}
	log.Fields = make([]Field, 0, len(protoLog.Fields))
	for _, protoFiled := range protoLog.Fields {
		log.Fields = append(log.Fields, Field{Key: protoFiled.Key, Value: protoFiled.Value})
	}
	return log
}

func fromProtoSpanInternal(protoSpan *protogen.Span) (*Span, error) {
	span := &Span{
		ServiceName:   protoSpan.ServiceName,
		OperationName: protoSpan.OperationName,
		TraceID:       protoSpan.TraceId,
		SpanID:        protoSpan.SpanId,
		Env:           protoSpan.Env,
		ParentID:      protoSpan.ParentId,
		StartTime:     timeFromTimestamp(protoSpan.StartTime),
		Duration:      durationFromDuration(protoSpan.Duration),
	}
	span.References = make([]SpanRef, 0, len(protoSpan.References))
	for _, ref := range protoSpan.References {
		span.References = append(span.References, convertSpanRef(ref))
	}
	span.Tags = make(map[string]interface{})
	for _, tag := range protoSpan.Tags {
		newTag := convertTag(tag)
		span.Tags[newTag.Key] = newTag.Value
	}
	span.Logs = make([]Log, 0, len(protoSpan.Logs))
	for _, log := range protoSpan.Logs {
		span.Logs = append(span.Logs, convertLog(log))
	}
	return span, nil
}
