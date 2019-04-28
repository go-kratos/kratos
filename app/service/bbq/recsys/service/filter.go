package service

import (
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	"go-common/app/service/bbq/recsys/service/retrieve"
	"go-common/library/log"
	"strconv"
)

//FilterManager ...
type FilterManager struct {
	filterNodes        []FilterNode
	relatedFilterNodes []FilterNode
	upsFilterNodes     []FilterNode
}

//FilterNode ...
type FilterNode interface {
	doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile)
}

// NewFilterManager new a filter manager
func NewFilterManager() (m *FilterManager) {
	m = &FilterManager{
		filterNodes:        make([]FilterNode, 0),
		relatedFilterNodes: make([]FilterNode, 0),
		upsFilterNodes:     make([]FilterNode, 0),
	}
	defaultFilterNode := &DefaultFilterNode{}
	bloomFilterNode := &BloomFilterNode{}
	durationFilterNode := &DurationFilterNode{}
	followsFilterNode := &FollowsFilterNode{}
	blackFilterNode := &BlackFilterNode{}
	m.filterNodes = append(m.filterNodes, defaultFilterNode, bloomFilterNode, blackFilterNode, durationFilterNode)

	relatedFilterNode := &RelatedFilterNode{}
	m.relatedFilterNodes = append(m.relatedFilterNodes, defaultFilterNode, relatedFilterNode)

	m.upsFilterNodes = append(m.upsFilterNodes, defaultFilterNode, blackFilterNode, followsFilterNode)

	return
}

func (m *FilterManager) filter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {
	for _, filterNode := range m.filterNodes {
		filterNode.doFilter(req, response, profile)
	}
}

func (m *FilterManager) relatedFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {
	for _, filterNode := range m.relatedFilterNodes {
		filterNode.doFilter(req, response, profile)
	}
}

func (m *FilterManager) upsFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {
	for _, filterNode := range m.upsFilterNodes {
		filterNode.doFilter(req, response, profile)
	}
}

//DefaultFilterNode ...
type DefaultFilterNode struct {
	FilterNode
}

func (f *DefaultFilterNode) doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {

	if req.DebugFlag {
		log.Info("Default Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}

	records := make([]*rpc.RecsysRecord, 0)
	viewedVideoSet := make(map[int64]int64)
	for _, record := range response.List {
		if _, ok := viewedVideoSet[record.Svid]; !ok {
			records = append(records, record)
		}
		viewedVideoSet[record.Svid] = 1
	}
	response.List = records

	if req.DebugFlag {
		log.Info("Default Filter Node after records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}

}

//BloomFilterNode ...
type BloomFilterNode struct {
	FilterNode
}

func (f *BloomFilterNode) doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {

	if req.DebugFlag {
		log.Info("Bloom Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}

	if response.Message[model.ResponseDownGrade] == "1" {
		if req.DebugFlag {
			log.Info("Do not do Bloom Filter in down grade state, traceID is (%v)", req.TraceID)
		}
		return
	}

	if profile.BloomFilter == nil {
		return
	}
	records := make([]*rpc.RecsysRecord, 0)

	for _, record := range response.List {
		svid := uint64(record.Svid)
		if !profile.BloomFilter.MightContainUint64(svid) {
			records = append(records, record)
		}
	}

	response.List = records
	if req.DebugFlag {
		log.Info("Bloom Filter Node after records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}
}

// RelatedFilterNode ...
type RelatedFilterNode struct {
	FilterNode
}

func (f *RelatedFilterNode) doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {

	records := make([]*rpc.RecsysRecord, 0)
	upMID := response.Message[retrieve.SourceUpMID]

	for _, record := range response.List {
		if record.Svid != req.SVID && record.Map[model.UperMid] != upMID {
			records = append(records, record)
		}
	}

	response.List = records
}

// FollowsFilterNode
type FollowsFilterNode struct {
	FilterNode
}

//去掉关注过的up主
func (f *FollowsFilterNode) doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {

	if req.DebugFlag {
		log.Info("Follow Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}

	records := make([]*rpc.RecsysRecord, 0)
	for _, record := range response.List {
		upMID, _ := strconv.ParseInt(record.Map[model.UperMid], 10, 64)
		if _, ok := profile.BBQFollow[upMID]; ok {
			continue
		}
		records = append(records, record)
	}
	response.List = records

	if req.DebugFlag {
		log.Info("Follow Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}
}

//BlackFilterNode ...
type BlackFilterNode struct {
	FilterNode
}

func (f *BlackFilterNode) doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {

	if req.DebugFlag {
		log.Info("Black Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}

	records := make([]*rpc.RecsysRecord, 0)
	for _, record := range response.List {
		upMID, _ := strconv.ParseInt(record.Map[model.UperMid], 10, 64)
		if _, ok := profile.BBQBlack[upMID]; ok {
			continue
		}
		records = append(records, record)
	}
	response.List = records

	if req.DebugFlag {
		log.Info("Black Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}
}

//DurationFilterNode ...
type DurationFilterNode struct {
	FilterNode
}

func (f *DurationFilterNode) doFilter(req *rpc.RecsysRequest, response *rpc.RecsysResponse, profile *model.UserProfile) {

	if req.DebugFlag {
		log.Info("Duration Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}

	records := make([]*rpc.RecsysRecord, 0)
	for _, record := range response.List {
		duration, _ := strconv.ParseInt(record.Map[model.Duration], 10, 64)
		if duration > 60 || duration < 15 {
			continue
		}
		records = append(records, record)
	}
	response.List = records

	if req.DebugFlag {
		log.Info("Duration Filter Node before records size: (%v), traceID is (%v)", len(response.List), req.TraceID)
	}
}
