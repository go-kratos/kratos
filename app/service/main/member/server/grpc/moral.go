package grpc

import (
	"context"
	"go-common/app/service/main/member/api"
)

// Moral Get member moral info
func (s *MemberServer) Moral(ctx context.Context, req *api.MemberMidReq) (*api.MoralReply, error) {
	res, err := s.svr.Moral(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	moralReply := &api.MoralReply{
		Mid:             res.Mid,
		Moral:           res.Moral,
		Added:           res.Added,
		Deducted:        res.Deducted,
		LastRecoverDate: res.LastRecoverDate,
	}

	return moralReply, nil
}

// MoralLog Get member moral logs
func (s *MemberServer) MoralLog(ctx context.Context, req *api.MemberMidReq) (*api.UserLogsReply, error) {
	res, err := s.svr.MoralLog(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	userLogs := make([]*api.UserLogReply, 0, len(res))
	for _, v := range res {
		userLog := &api.UserLogReply{
			Mid:     v.Mid,
			Ip:      v.IP,
			Ts:      v.TS,
			LogId:   v.LogID,
			Content: v.Content,
		}
		userLogs = append(userLogs, userLog)
	}
	userLogsReply := &api.UserLogsReply{
		UserLogs: userLogs,
	}

	return userLogsReply, nil
}

// AddMoral Add member's moral value
func (s *MemberServer) AddMoral(ctx context.Context, req *api.UpdateMoralReq) (*api.EmptyStruct, error) {
	err := s.svr.UpdateMoral(ctx, api.ToArgUpdateMoral(req))
	if err != nil {
		return nil, err
	}
	return &api.EmptyStruct{}, nil
}

// BatchAddMoral Batch add member's moral value
func (s *MemberServer) BatchAddMoral(ctx context.Context, req *api.UpdateMoralsReq) (*api.UpdateMoralsReply, error) {
	res, err := s.svr.UpdateMorals(ctx, api.ToArgUpdateMorals(req))
	if err != nil {
		return nil, err
	}
	updateMoralsReply := &api.UpdateMoralsReply{
		AfterMorals: res,
	}
	return updateMoralsReply, nil
}
