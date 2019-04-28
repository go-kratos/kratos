package grpc

import (
	"context"
	"go-common/app/service/main/member/api"
	"go-common/app/service/main/member/service"
	"go-common/app/service/main/member/service/block"
	"go-common/library/net/rpc/warden"
)

// New Member warden rpc server
func New(cfg *warden.ServerConfig, s *service.Service) *warden.Server {
	w := warden.NewServer(cfg)
	api.RegisterMemberServer(w.Server(), &MemberServer{svr: s, blockSvr: s.BlockImpl()})

	ws, err := w.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

// MemberServer define member Service
type MemberServer struct {
	svr      *service.Service
	blockSvr *block.Service
}

var _ api.MemberServer = &MemberServer{}

// Base get member base info
func (s *MemberServer) Base(ctx context.Context, req *api.MemberMidReq) (*api.BaseInfoReply, error) {
	res, err := s.svr.BaseInfo(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	return api.FromBaseInfo(res), nil
}

// Bases batch get members base info
func (s *MemberServer) Bases(ctx context.Context, req *api.MemberMidsReq) (*api.BaseInfosReply, error) {
	res, err := s.svr.BatchBaseInfo(ctx, req.Mids)
	if err != nil {
		return nil, err
	}

	baseInfos := make(map[int64]*api.BaseInfoReply, len(res))
	baseInfosReply := &api.BaseInfosReply{
		BaseInfos: baseInfos,
	}

	for k, v := range res {
		baseInfos[k] = api.FromBaseInfo(v)
	}

	return baseInfosReply, nil
}

// Member get member full information
func (s *MemberServer) Member(ctx context.Context, req *api.MemberMidReq) (*api.MemberInfoReply, error) {
	res, err := s.svr.Member(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	memberInfoReply := api.FromMember(res)
	return memberInfoReply, nil
}

// Members Batch get members info
func (s *MemberServer) Members(ctx context.Context, req *api.MemberMidsReq) (*api.MemberInfosReply, error) {
	res, err := s.svr.Members(ctx, req.Mids)
	if err != nil {
		return nil, err
	}

	memberInfos := make(map[int64]*api.MemberInfoReply, len(res))
	for k, v := range res {
		memberInfoReply := api.FromMember(v)
		memberInfos[k] = memberInfoReply
	}
	memberInfosReply := &api.MemberInfosReply{
		MemberInfos: memberInfos,
	}

	return memberInfosReply, nil
}

// NickUpdated Whether the member's nickname has been updated
func (s *MemberServer) NickUpdated(ctx context.Context, req *api.MemberMidReq) (*api.NickUpdatedReply, error) {
	res, err := s.svr.NickUpdated(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	nickUpdatedReply := &api.NickUpdatedReply{
		NickUpdated: res,
	}
	return nickUpdatedReply, nil
}

// SetNickUpdated Mark nickname as updated
func (s *MemberServer) SetNickUpdated(ctx context.Context, req *api.MemberMidReq) (*api.EmptyStruct, error) {
	err := s.svr.SetNickUpdated(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetOfficialDoc Set offical document
func (s *MemberServer) SetOfficialDoc(ctx context.Context, req *api.OfficialDocReq) (*api.EmptyStruct, error) {
	err := s.svr.SetOfficialDoc(ctx, api.ToArgOfficialDoc(req))
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetSex Set member's sex
func (s *MemberServer) SetSex(ctx context.Context, req *api.UpdateSexReq) (*api.EmptyStruct, error) {
	err := s.svr.SetSex(ctx, req.Mid, req.Sex)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetName Set member's name
func (s *MemberServer) SetName(ctx context.Context, req *api.UpdateUnameReq) (*api.EmptyStruct, error) {
	err := s.svr.SetName(ctx, req.Mid, req.Name)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetFace Set member's face
func (s *MemberServer) SetFace(ctx context.Context, req *api.UpdateFaceReq) (*api.EmptyStruct, error) {
	err := s.svr.SetFace(ctx, req.Mid, req.Face)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetRank Set member's rank
func (s *MemberServer) SetRank(ctx context.Context, req *api.UpdateRankReq) (*api.EmptyStruct, error) {
	err := s.svr.SetRank(ctx, req.Mid, req.Rank)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetBirthday Set member's birthday
func (s *MemberServer) SetBirthday(ctx context.Context, req *api.UpdateBirthdayReq) (*api.EmptyStruct, error) {
	err := s.svr.SetBirthday(ctx, req.Mid, req.Birthday)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// SetSign Set member's sign
func (s *MemberServer) SetSign(ctx context.Context, req *api.UpdateSignReq) (*api.EmptyStruct, error) {
	err := s.svr.SetSign(ctx, req.Mid, req.Sign)
	if err != nil {
		return nil, err
	}
	emptyStruct := &api.EmptyStruct{}
	return emptyStruct, nil
}

// OfficialDoc Get member's offical doc
func (s *MemberServer) OfficialDoc(ctx context.Context, req *api.MidReq) (*api.OfficialDocInfoReply, error) {
	res, err := s.svr.OfficialDoc(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	officialDocInfoReply := api.FromOfficialDoc(res)
	return officialDocInfoReply, nil
}
