package grpc

import (
	"context"
	"go-common/app/service/main/member/api"
	"go-common/app/service/main/member/model"
)

// AddUserMonitor add user monitor
func (s *MemberServer) AddUserMonitor(ctx context.Context, req *api.AddUserMonitorReq) (*api.EmptyStruct, error) {
	argAddUserMonitor := &model.ArgAddUserMonitor{
		Mid:      req.Mid,
		Operator: req.Operator,
		Remark:   req.Remark,
	}
	err := s.svr.AddUserMonitor(ctx, argAddUserMonitor)
	if err != nil {
		return nil, err
	}

	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// IsInMonitor check whether the member is in monitored status
func (s *MemberServer) IsInMonitor(ctx context.Context, req *api.MidReq) (*api.IsInMonitorReply, error) {
	res, err := s.svr.IsInMonitor(ctx, &model.ArgMid{
		Mid:    req.Mid,
		RealIP: req.RealIP,
	})
	if err != nil {
		return nil, err
	}

	isInMonitorReply := &api.IsInMonitorReply{
		IsInMonitor: res,
	}
	return isInMonitorReply, nil
}
