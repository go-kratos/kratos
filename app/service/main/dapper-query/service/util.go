package service

import (
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	apiv1 "go-common/app/service/main/dapper-query/api/v1"
	"go-common/app/service/main/dapper-query/model"
)

func compatibleLegacySpan(spans []*model.Span) bool {
	var fixed bool
	set := make(map[uint64]*model.Span)
	for _, sp1 := range spans {
		if sp2, ok := set[sp1.SpanID]; ok {
			fixed = fixed || fixParentID(sp1, sp2)
			delete(set, sp1.SpanID)
		} else {
			set[sp1.SpanID] = sp1
		}
	}
	return fixed
}

func fixParentID(sp1, sp2 *model.Span) bool {
	var client, server *model.Span
	for _, sp := range []*model.Span{sp1, sp2} {
		if sp.IsServer() {
			server = sp
		} else {
			client = sp
		}
	}
	if client == nil || server == nil {
		return false
	}
	server.ParentID = client.SpanID
	return true
}

func setChilds(node *apiv1.Span, parentMap map[string][]*apiv1.Span, level int32) int32 {
	spans, ok := parentMap[node.SpanId]
	if !ok {
		return level
	}
	level++
	delete(parentMap, node.SpanId)
	// compatible old span pair, client server has same span_id, parent_span_id
	for _, span := range spans {
		span.Level = int32(level)
		node.Childs = append(node.Childs, span)
	}
	if node.Childs != nil {
		sort.Slice(node.Childs, func(i, j int) bool {
			iStartTime := node.Childs[i].StartTime
			jStartTime := node.Childs[j].StartTime
			return iStartTime < jStartTime
		})
	}
	var newLevel int32
	for _, cnode := range node.Childs {
		if ret := setChilds(cnode, parentMap, level); ret > newLevel {
			newLevel = ret
		}
	}
	return newLevel
}

func protoTimestamp(t time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

func protoDuration(d time.Duration) *duration.Duration {
	return &duration.Duration{
		Seconds: d.Nanoseconds() / int64(time.Second),
		Nanos:   int32(d.Nanoseconds() % int64(time.Second)),
	}
}
