package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/service/main/dapper-query/dao"
	"go-common/app/service/main/dapper-query/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// timeOffset secode
var _timeOffset int64 = 1800

// DependItem service depend
type DependItem struct {
	ServiceName   string
	Component     string
	OperationName string
}

func parseDependItem(spans []*model.Span, serviceName, operationName string) (depends []DependItem) {
	spanMap := make(map[uint64][]*model.Span)
	var root *model.Span
	for _, span := range spans {
		if span.ServiceName == serviceName && span.IsServer() {
			root = span
			continue
		}
		spanMap[span.ParentID] = append(spanMap[span.ParentID], span)
	}
	if root == nil {
		return
	}
	for _, span := range spanMap[root.SpanID] {
		if span.IsServer() {
			continue
		}
		if peerSpans, ok := spanMap[span.SpanID]; ok {
			span = peerSpans[0]
			depends = append(depends, DependItem{
				Component:     span.StringTag("component"),
				ServiceName:   span.ServiceName,
				OperationName: span.OperationName,
			})
		} else {
			peerService := span.StringTag("peer.service")
			// in old dapper sdk service such as redis, memcache, mysql be save as service_name not in peer.service
			if peerService == "" {
				peerService = span.ServiceName
			}
			depends = append(depends, DependItem{
				Component:     span.StringTag("component"),
				ServiceName:   peerService,
				OperationName: span.OperationName,
			})
		}
	}
	return depends
}

func (s *service) fetchServiceDepend(ctx context.Context, serviceName string, operationNames []string) []DependItem {
	var mx sync.Mutex
	var depends []DependItem
	appendDepend := func(items []DependItem) {
		mx.Lock()
		depends = append(depends, items...)
		mx.Unlock()
	}
	end := time.Now().Unix() - _timeOffset
	start := end - 3600
	sel := &dao.Selector{Limit: 3, Start: start, End: end}

	group := &errgroup.Group{}
	for i := range operationNames {
		operationName := operationNames[i]
		group.Go(func() error {
			refs, err := s.daoImpl.QuerySpanList(ctx, serviceName, operationName, sel, dao.TimeDesc)
			if err != nil && err != dao.ErrNotFound {
				if err != dao.ErrNotFound {
					// only log error don't return error
					log.Warn("query span list error: %s", err)
				}
				return nil
			}
			for _, ref := range refs {
				spans, err := s.daoImpl.Trace(ctx, ref.TraceID)
				if err != nil {
					if err != dao.ErrNotFound {
						// only log error don't return error
						log.Warn("fetch trace span error: %s", err)
					}
				}
				appendDepend(parseDependItem(spans, serviceName, operationName))
			}
			return nil
		})
	}
	group.Wait()
	return depends
}
