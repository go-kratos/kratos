package opslog

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/app/service/main/dapper-query/api/v1"
	"go-common/app/service/main/dapper-query/conf"
	"go-common/app/service/main/dapper-query/dao"
	"go-common/app/service/main/dapper-query/model"
	"go-common/app/service/main/dapper-query/pkg/opslog"
	"go-common/app/service/main/dapper-query/util"
)

// Service OpslogService
type Service struct {
	daoImpl dao.Dao
	client  opslog.Client
}

// NewOpsLogService opslog service provide opslog query service
func NewOpsLogService(cfg *conf.Config, d dao.Dao) *Service {
	client := opslog.New(cfg.OpsLog.API, nil)
	return &Service{daoImpl: d, client: client}
}

// OpsLog 获取 OpsLog 数据
func (s *Service) OpsLog(ctx context.Context, req *v1.OpsLogReq) (*v1.OpsLogReply, error) {
	traceIDStr := strings.Split(req.TraceId, ":")[0]
	traceID, err := strconv.ParseUint(traceIDStr, 16, 64)
	if err != nil {
		return nil, ecode.Error(ecode.RequestErr, err.Error())
	}
	var spanIDs []uint64
	spanID, err := strconv.ParseUint(req.SpanId, 16, 64)
	if err == nil {
		spanIDs = append(spanIDs, spanID)
	}
	spans, err := s.daoImpl.Trace(ctx, traceID, spanIDs...)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	if len(spans) != 0 {
		return s.opsLogTraceExists(ctx, traceID, spanID, req.TraceField, spans)
	}
	return s.opslogFromGuess(ctx, traceID, req)
}

func (s *Service) opslogFromGuess(ctx context.Context, traceID uint64, req *v1.OpsLogReq) (*v1.OpsLogReply, error) {
	if req.ServiceName == "" || req.OperationName == "" {
		return nil, ecode.Error(ecode.RequestErr, "service_name and operation_name required")
	}
	if req.End == 0 {
		req.End = time.Now().Unix()
	}
	if req.Start == 0 {
		req.Start = req.End - 3600
	}
	// TODO: get better way to query service name depends
	refs, err := s.daoImpl.QuerySpanList(ctx, req.ServiceName, req.OperationName, &dao.Selector{Start: req.Start, End: req.End, Limit: 3}, dao.TimeDesc)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	serviceMap := make(map[string]struct{})
	var spans []*model.Span
	for _, ref := range refs {
		if ref.IsError {
			continue
		}
		spans, err = s.daoImpl.Trace(ctx, ref.TraceID)
		if err != nil {
			log.Warn("get trace %x error: %s", ref.TraceID, err)
		}
		for _, span := range spans {
			serviceMap[span.ServiceName] = struct{}{}
		}
	}
	serviceNames := make([]string, 0, len(serviceMap))
	for key := range serviceMap {
		serviceNames = append(serviceNames, key)
	}
	records, err := s.client.Query(ctx, serviceNames, traceID, util.SessionIDFromContext(ctx), req.Start, req.End)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	return &v1.OpsLogReply{Records: toAPIRecords(records)}, nil
}

func (s *Service) opsLogTraceExists(ctx context.Context, traceID, spanID uint64, traceField string, spans []*model.Span) (*v1.OpsLogReply, error) {
	var start, end int64
	serviceMap := make(map[string]bool)
	for _, span := range spans {
		if start == 0 {
			// start set 300 second before span start_time
			start = span.StartTime.Unix() - 300
			// end set 300 second after span start_time
			end = start + 600
		}
		if !span.IsServer() {
			continue
		}
		if spanID != 0 && span.SpanID != spanID {
			continue
		}
		serviceMap[span.ServiceName] = true
	}
	serviceNames := make([]string, 0, len(serviceMap))
	for key := range serviceMap {
		serviceNames = append(serviceNames, key)
	}
	var options []opslog.Option
	if traceField != "" {
		options = append(options, opslog.SetTraceField(traceField))
	}
	records, err := s.client.Query(ctx, serviceNames, traceID, util.SessionIDFromContext(ctx), start, end, options...)
	if err != nil {
		return nil, ecode.Error(ecode.ServerErr, err.Error())
	}
	return &v1.OpsLogReply{Records: toAPIRecords(records)}, nil
}

func toAPIRecords(records []*opslog.Record) []*v1.OpsLogRecord {
	apiRecords := make([]*v1.OpsLogRecord, 0, len(records))
	for _, record := range records {
		apiRecord := &v1.OpsLogRecord{
			Time:    record.Time.Format("2006-01-02T15:04:05.000"),
			Level:   record.Level,
			Message: record.Message,
			Fields:  make(map[string]*v1.TagValue),
		}
		for k, v := range record.Fields {
			switch val := v.(type) {
			case string:
				apiRecord.Fields[k] = &v1.TagValue{Value: &v1.TagValue_StringValue{StringValue: val}}
			case bool:
				apiRecord.Fields[k] = &v1.TagValue{Value: &v1.TagValue_BoolValue{BoolValue: val}}
			case float64:
				apiRecord.Fields[k] = &v1.TagValue{Value: &v1.TagValue_FloatValue{FloatValue: float32(val)}}
			}
		}
		apiRecords = append(apiRecords, apiRecord)
	}
	return apiRecords
}
