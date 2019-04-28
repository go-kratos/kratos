package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/library/ecode"

	"go-common/app/service/main/dapper-query/api/v1"
	"go-common/app/service/main/dapper-query/conf"
	"go-common/app/service/main/dapper-query/dao"
	"go-common/app/service/main/dapper-query/model"
	"go-common/app/service/main/dapper-query/pkg/cltclient"
	"go-common/app/service/main/dapper-query/service/opslog"
	"go-common/library/log"
)

const (
	defaultTimeRange = 3600
)

type service struct {
	opsLog  *opslog.Service
	daoImpl dao.Dao
	clt     *cltclient.Client
}

var _ v1.BMDapperQueryServer = &service{}

// New DapperQueryService implement DapperQueryServer
func New(cfg *conf.Config) (v1.DapperQueryServer, error) {
	daoImpl, err := dao.New(cfg)
	if err != nil {
		return nil, err
	}
	clt, err := cltclient.New(cfg.Collectors.Nodes, nil)
	if err != nil {
		return nil, err
	}
	opsLog := opslog.NewOpsLogService(cfg, daoImpl)
	return &service{daoImpl: daoImpl, clt: clt, opsLog: opsLog}, nil
}

// OpsLog 获取 OpsLog 数据
func (s *service) OpsLog(ctx context.Context, req *v1.OpsLogReq) (*v1.OpsLogReply, error) {
	return s.opsLog.OpsLog(ctx, req)
}

func (s *service) ListServiceName(ctx context.Context, req *v1.ListServiceNameReq) (*v1.ListServiceNameReply, error) {
	serviceNames, err := s.daoImpl.ServiceNames(ctx)
	if err != nil {
		return nil, err
	}
	reply := &v1.ListServiceNameReply{
		ServiceNames: serviceNames,
	}
	return reply, nil
}

func (s *service) ListOperationName(ctx context.Context, req *v1.ListOperationNameReq) (*v1.ListOperationNameReply, error) {
	serviceName := req.ServiceName
	operationNames, err := s.daoImpl.OperationNames(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	reply := &v1.ListOperationNameReply{OperationNames: operationNames}
	return reply, nil
}

func (s *service) ListSpan(ctx context.Context, req *v1.ListSpanReq) (*v1.ListSpanReply, error) {
	serviceName := req.ServiceName
	operationName := req.OperationName

	start := req.Start
	end := req.End
	// if start or end not set, set time range to last hour
	if start == 0 || end == 0 {
		end = time.Now().Unix()
		start = end - defaultTimeRange
	}

	order := req.Order
	if order == "" {
		order = dao.TimeDesc
	}
	if !s.daoImpl.SupportOrder(order) {
		return nil, ecode.Errorf(ecode.RequestErr, "request order %s unsupport yet ㄟ( ▔, ▔  )ㄏ", order)
	}

	onlyError := req.OnlyError
	offset := int(req.Offset)

	limit := int(req.Limit)
	if limit == 0 {
		limit = 50
	}
	sel := &dao.Selector{Start: start, End: end, Offset: offset, Limit: limit, OnlyError: onlyError}
	spanListRefs, err := s.daoImpl.QuerySpanList(ctx, serviceName, operationName, sel, order)
	if err != nil {
		log.Error("query span list error: %s", err.Error())
	}
	reply := &v1.ListSpanReply{Items: make([]*v1.SpanListItem, 0, len(spanListRefs))}
	for _, spanListRef := range spanListRefs {
		item, err := s.getSpanListItem(ctx, serviceName, spanListRef)
		if err != nil {
			log.Error("query span list item error: %s", err)
			continue
		}
		reply.Items = append(reply.Items, item)
	}
	return reply, nil
}

func (s *service) getSpanListItem(ctx context.Context, serviceName string, spanListRef model.SpanListRef) (*v1.SpanListItem, error) {
	spans, err := s.daoImpl.Trace(ctx, spanListRef.TraceID, spanListRef.SpanID)
	if err != nil {
		return nil, err
	}
	var currentSpan *model.Span
	for _, span := range spans {
		if span.SpanID == spanListRef.SpanID {
			currentSpan = span
		}
	}
	if currentSpan == nil {
		return nil, fmt.Errorf("can't find span: %x in traceid: %x", spanListRef.SpanID, spanListRef.TraceID)
	}

	var mark string
	if currentSpan.IsError() {
	rangelog:
		for _, log := range currentSpan.Logs {
			for _, field := range log.Fields {
				if field.Key == "message" {
					mark = string(field.Value)
					break rangelog
				}
			}
		}
	}

	item := &v1.SpanListItem{
		TraceId:       spanListRef.TraceIDStr(),
		SpanId:        spanListRef.SpanIDStr(),
		ParentId:      strconv.FormatUint(currentSpan.ParentID, 16),
		ServiceName:   currentSpan.ServiceName,
		OperationName: currentSpan.OperationName,
		StartTime:     currentSpan.StartTime.Format("2006-01-02T15:04:05.000"),
		Duration:      currentSpan.Duration.String(),
		IsError:       currentSpan.IsError(),
		RegionZone:    currentSpan.StringTag("region") + ":" + currentSpan.StringTag("zone"),
		ContainerIp:   currentSpan.StringTag("ip"),
		Mark:          mark,
		Tags:          toAPITags(currentSpan.Tags),
	}
	return item, nil
}

func (s *service) Trace(ctx context.Context, req *v1.TraceReq) (*v1.TraceReply, error) {
	traceIDStr := strings.Split(req.TraceId, ":")[0]
	traceID, err := strconv.ParseUint(traceIDStr, 16, 64)
	if err != nil {
		return nil, ecode.Errorf(ecode.RequestErr, "invalid traceID %s", req.TraceId)
	}
	spanID, _ := strconv.ParseUint(req.SpanId, 16, 64)
	spans, err := s.daoImpl.Trace(ctx, traceID)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	var apiSpans []*v1.Span
	var root *v1.Span
	compatibleLegacySpan(spans)
	serviceMap := make(map[string]bool)
	for _, span := range spans {
		apiSpan := s.toAPISpan(span)
		if root == nil && ((spanID == 0 && span.IsServer()) || (span.SpanID == spanID && span.IsServer())) {
			root = apiSpan
		}
		apiSpans = append(apiSpans, apiSpan)
		serviceMap[span.ServiceName] = true
		if peerService := span.StringTag("peer.service"); peerService != "" {
			serviceMap[peerService] = true
		}
	}
	if root == nil {
		// NOTE: return root=nil trace if not found, it help web judge
		return &v1.TraceReply{}, nil
	}
	root.Level = 1
	parentMap := make(map[string][]*v1.Span)
	for _, apiSpan := range apiSpans {
		// skip root span
		if apiSpan.SpanId == root.SpanId {
			continue
		}
		parentMap[apiSpan.ParentId] = append(parentMap[apiSpan.ParentId], apiSpan)
	}
	maxLevel := setChilds(root, parentMap, root.Level)
	reply := &v1.TraceReply{
		Root:         root,
		SpanCount:    int32(len(spans) - len(parentMap)),
		MaxLevel:     maxLevel,
		ServiceCount: int32(len(serviceMap)),
	}
	return reply, nil
}

func (s *service) RawTrace(ctx context.Context, req *v1.RawTraceReq) (*v1.RawTraceReply, error) {
	traceID, err := strconv.ParseUint(req.TraceId, 16, 64)
	if err != nil {
		return nil, ecode.Errorf(ecode.RequestErr, "invalid traceID %s", req.TraceId)
	}
	spans, err := s.daoImpl.Trace(ctx, traceID)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	reply := &v1.RawTraceReply{
		Items: make([]*v1.Span, 0, len(spans)),
	}
	for _, span := range spans {
		reply.Items = append(reply.Items, s.toAPISpan(span))
	}
	sort.Slice(reply.Items, func(i, j int) bool {
		return reply.Items[i].StartTime > reply.Items[j].StartTime
	})
	return reply, nil
}

func toAPITags(tags map[string]interface{}) map[string]*v1.TagValue {
	apiTags := make(map[string]*v1.TagValue)
	for key, value := range tags {
		apitag := &v1.TagValue{}
		switch val := value.(type) {
		case string:
			apitag.Value = &v1.TagValue_StringValue{StringValue: val}
		case int64:
			apitag.Value = &v1.TagValue_Int64Value{Int64Value: val}
		case bool:
			apitag.Value = &v1.TagValue_BoolValue{BoolValue: val}
		case float64:
			apitag.Value = &v1.TagValue_FloatValue{FloatValue: float32(val)}
		}
		apiTags[key] = apitag
	}
	return apiTags
}

func (s *service) toAPISpan(span *model.Span) *v1.Span {
	apiSpan := &v1.Span{
		ServiceName:   span.ServiceName,
		OperationName: span.OperationName,
		TraceId:       span.TraceIDStr(),
		SpanId:        span.SpanIDStr(),
		ParentId:      span.ParentIDStr(),
		StartTime:     span.StartTime.UnixNano(),
		Duration:      int64(span.Duration),
	}
	for _, log := range span.Logs {
		apilog := &v1.Log{Timestamp: log.Timestamp}
		for _, field := range log.Fields {
			apilog.Fields = append(apilog.Fields, &v1.Field{Key: field.Key, Value: string(field.Value)})
		}
		apiSpan.Logs = append(apiSpan.Logs, apilog)
	}
	apiSpan.Tags = toAPITags(span.Tags)
	return apiSpan
}

// OperationNameRank 查询 OperationName 排名列表
func (s *service) OperationNameRank(ctx context.Context, req *v1.OperationNameRankReq) (*v1.OperationNameRankReply, error) {
	serviceName := req.ServiceName
	start, end := req.Start, req.End
	// if start or end not set, set time range to last hour
	if start == 0 || end == 0 {
		end = time.Now().Unix()
		start = end - defaultTimeRange
	}
	rankType := req.RankType
	if rankType != "" {
		if !model.VerifyRankType(rankType) {
			return nil, ecode.Errorf(ecode.RequestErr, "request rankType %s unsupport yet", rankType)
		}
	} else {
		rankType = model.MaxDurationRank
	}
	values, err := s.daoImpl.MeanOperationNameField(ctx, map[string]string{
		dao.ServiceNameTag: serviceName,
		dao.SpanKindTag:    "server",
	}, rankType, start, end, []string{dao.OperationNameTag})
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	model.SortRank(values)
	reply := &v1.OperationNameRankReply{RankType: rankType, Items: make([]*v1.RankItem, 0, len(values))}
	for _, value := range values {
		operationName := value.Tag[dao.OperationNameTag]
		reply.Items = append(reply.Items, &v1.RankItem{
			ServiceName:   serviceName,
			OperationName: operationName,
			Value:         value.Value,
		})
	}
	return reply, nil
}

// DependsRank DependsRank
func (s *service) DependsRank(ctx context.Context, req *v1.DependsRankReq) (*v1.DependsRankReply, error) {
	serviceName := req.ServiceName
	start, end := req.Start, req.End
	// if start or end not set, set time range to last hour
	if start == 0 || end == 0 {
		end = time.Now().Unix()
		start = end - defaultTimeRange
	}
	rankType := req.RankType
	if rankType != "" {
		if !model.VerifyRankType(rankType) {
			return nil, ecode.Errorf(ecode.RequestErr, "request rankType %s unsupport yet", rankType)
		}
	} else {
		rankType = model.MaxDurationRank
	}
	values, err := s.daoImpl.MeanOperationNameField(ctx, map[string]string{
		dao.ServiceNameTag: serviceName,
		dao.SpanKindTag:    "client",
	}, rankType, start, end, []string{dao.PeerServiceTag, dao.OperationNameTag})
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	reply := &v1.DependsRankReply{RankType: rankType, Items: make([]*v1.RankItem, 0, len(values))}
	for _, value := range values {
		operationName := value.Tag[dao.OperationNameTag]
		serviceName := value.Tag[dao.PeerServiceTag]
		reply.Items = append(reply.Items, &v1.RankItem{
			ServiceName:   serviceName,
			OperationName: operationName,
			Value:         value.Value,
		})
	}
	return reply, nil
}

// SpanSeries 获取 span 的时间序列数据
func (s *service) SpanSeries(ctx context.Context, req *v1.SpanSeriesReq) (*v1.SpanSeriesReply, error) {
	start, end := req.Start, req.End
	// if start or end not set, set time range to last hour
	if start == 0 || end == 0 {
		end = time.Now().Unix()
		start = end - defaultTimeRange
	}
	interval := (end - start) / 120
	if interval < 5 {
		interval = 5
	}
	serviceName := req.ServiceName
	operationName := req.OperationName
	fields := strings.Split(req.Fields, ",")
	seriesFn := s.daoImpl.SpanSeriesMean
	// 为错误序列特别处理
	if len(fields) == 1 && fields[0] == "errors" {
		seriesFn = s.daoImpl.SpanSeriesCount
	}
	series, err := seriesFn(ctx, serviceName, operationName, fields, start, end, interval)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	reply := &v1.SpanSeriesReply{Interval: int64(interval)}
	reply.Times = make([]string, len(series.Timestamps))
	for i, timestamp := range series.Timestamps {
		formatedTime := time.Unix(timestamp/int64(time.Second), timestamp%int64(time.Second)).Format("2006-01-02T15:04:05")
		reply.Times[i] = formatedTime
	}
	reply.Items = make([]*v1.SeriesItem, len(series.Items))
	for i, item := range series.Items {
		apiItem := &v1.SeriesItem{Field: item.Field, Values: make([]*int64, len(item.Rows))}
		for j, val := range item.Rows {
			if val == nil {
				apiItem.Values[j] = nil
			} else {
				valInt64 := int64(*val)
				apiItem.Values[j] = &valInt64
			}
		}
		reply.Items[i] = apiItem
	}
	return reply, nil
}

func (s *service) SamplePoint(ctx context.Context, req *v1.SamplePointReq) (*v1.SamplePointReply, error) {
	serviceName := req.ServiceName
	operationName := req.OperationName
	startTime, err := time.ParseInLocation("2006-01-02T15:04:05", req.Time, time.Local)
	if err != nil {
		return nil, ecode.Error(ecode.RequestErr, err.Error())
	}
	start, end := startTime.Unix()-req.Interval, startTime.Unix()+req.Interval
	spanListRef, err := s.daoImpl.QuerySpanList(ctx, serviceName, operationName, &dao.Selector{Start: start, End: end, Limit: 50, OnlyError: req.OnlyError}, dao.TimeDesc)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	reply := &v1.SamplePointReply{}
	for _, ref := range spanListRef {
		reply.Items = append(reply.Items, &v1.SamplePointItem{
			TraceId:  ref.TraceIDStr(),
			SpanId:   ref.SpanIDStr(),
			Duration: ref.Duration,
			IsError:  ref.IsError,
		})
	}
	return reply, nil
}

// CltStatus CltStatus
func (s *service) CltStatus(ctx context.Context, req *v1.CltStatusReq) (*v1.CltStatusReply, error) {
	nodes, err := s.clt.Status(ctx)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	reply := new(v1.CltStatusReply)
	for _, node := range nodes {
		apiNode := &v1.CltNode{
			Node:     node.Node,
			QueueLen: int64(node.QueueLen),
		}
		for _, client := range node.Clients {
			apiNode.Clients = append(apiNode.Clients, &v1.Client{
				Addr:     client.Addr,
				ErrCount: client.ErrCount,
				Rate:     client.Rate,
				UpTime:   client.UpTime,
			})
		}
		reply.Nodes = append(reply.Nodes, apiNode)
	}
	return reply, nil
}

// DependsTopology 依赖拓扑
func (s *service) DependsTopology(ctx context.Context, req *v1.DependsTopologyReq) (*v1.DependsTopologyReply, error) {
	serviceNames, err := s.daoImpl.ServiceNames(ctx)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	reply := &v1.DependsTopologyReply{}
	for _, serviceName := range serviceNames {
		peerServices, err := s.daoImpl.PeerService(ctx, serviceName)
		if err != nil {
			return nil, ecode.Error(ecode.ServerErr, err.Error())
		}
		log.V(10).Info("peerService: %+v", peerServices)
		for _, peerService := range peerServices {
			reply.Items = append(reply.Items, &v1.DependsTopologyItem{ServiceName: serviceName, DependOn: peerService})
		}
	}
	return reply, nil
}

// ServiceDepend 查询服务依赖
func (s *service) ServiceDepend(ctx context.Context, req *v1.ServiceDependReq) (*v1.ServiceDependReply, error) {
	operationNames, err := s.daoImpl.OperationNames(ctx, req.ServiceName)
	if err != nil {
		if err == dao.ErrNotFound {
			return nil, ecode.Errorf(ecode.NothingFound, "service %s not exists", req.ServiceName)
		}
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	if req.OperationName != "" {
		for _, operationName := range operationNames {
			if operationName == req.OperationName {
				operationNames = []string{operationName}
				break
			}
		}
		if len(operationNames) != 1 {
			return nil, ecode.Errorf(ecode.RequestErr, "operationName %s not exists for service %s", req.OperationName, req.ServiceName)
		}
	}
	items := s.fetchServiceDepend(ctx, req.ServiceName, operationNames)
	return &v1.ServiceDependReply{Items: mergeDepends(items)}, nil
}

func mergeDepends(items []DependItem) []*v1.ServiceDependItem {
	var apiItems []*v1.ServiceDependItem
	apiItemsExists := func(serviceName, component string) (*v1.ServiceDependItem, bool) {
		for _, apiItem := range apiItems {
			if apiItem.ServiceName == serviceName && apiItem.Component == component {
				return apiItem, true
			}
		}
		return nil, false
	}
	addOperationName := func(apiItem *v1.ServiceDependItem, operationName string) {
		for _, name := range apiItem.OperationNames {
			if name == operationName {
				return
			}
		}
		apiItem.OperationNames = append(apiItem.OperationNames, operationName)
	}
	for _, item := range items {
		if apiItem, ok := apiItemsExists(item.ServiceName, item.Component); ok {
			addOperationName(apiItem, item.OperationName)
		} else {
			apiItems = append(apiItems, &v1.ServiceDependItem{
				ServiceName:    item.ServiceName,
				Component:      item.Component,
				OperationNames: []string{item.OperationName},
			})
		}
	}
	return apiItems
}
