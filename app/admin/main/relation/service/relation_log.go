package service

import (
	"context"
	"time"

	"go-common/app/admin/main/relation/model"
)

// RelationLog is.
func (s *Service) RelationLog(ctx context.Context, mid, fid int64, from time.Time, to time.Time) (model.RelationLogList, error) {
	logs, err := s.dao.RelationLogs(ctx, mid, fid, from, to)
	if err != nil {
		return nil, err
	}
	// order log by mtime with desc
	logs.OrderByMTime(true)

	uids := make([]int64, 0, len(logs)*2)
	for _, l := range logs {
		uids = append(uids, l.Mid)
		uids = append(uids, l.Fid)
	}
	uinfos, err := s.dao.RPCInfos(ctx, uids)
	if err != nil {
		return nil, err
	}
	for _, l := range logs {
		if mi, ok := uinfos[l.Mid]; ok {
			l.MemberName = mi.Name
		}
	}
	for _, l := range logs {
		if fi, ok := uinfos[l.Fid]; ok {
			l.FollowingName = fi.Name
		}
	}

	return logs, nil
}
