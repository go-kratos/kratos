package grpc

import (
	"context"
	"go-common/app/service/main/member/api"
)

// RealnameStatus get the member realname status
func (s *MemberServer) RealnameStatus(ctx context.Context, req *api.MemberMidReq) (*api.RealnameStatusReply, error) {
	res, err := s.svr.RealnameStatus(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	var realnameStatusReply = &api.RealnameStatusReply{
		RealnameStatus: int8(res),
	}

	return realnameStatusReply, nil
}

// RealnameApplyStatus get member realname apply status
func (s *MemberServer) RealnameApplyStatus(ctx context.Context, req *api.MemberMidReq) (*api.RealnameApplyInfoReply, error) {
	res, err := s.svr.RealnameApplyStatus(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	var realnameStatusReply = &api.RealnameApplyInfoReply{
		Status: int8(res.Status),
		Remark: res.Remark,
	}

	return realnameStatusReply, nil
}

// RealnameTelCapture mobilePhone realname certification
func (s *MemberServer) RealnameTelCapture(ctx context.Context, req *api.MemberMidReq) (*api.EmptyStruct, error) {
	_, err := s.svr.RealnameTelCapture(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	return &api.EmptyStruct{}, nil
}

// RealnameApply apply for realname certification
func (s *MemberServer) RealnameApply(ctx context.Context, req *api.ArgRealnameApplyReq) (*api.EmptyStruct, error) {
	err := s.svr.RealnameApply(ctx, req.Mid, int(req.CaptureCode), req.Realname, req.CardType, req.CardCode, req.Country, req.HandIMGToken, req.FrontIMGToken, req.BackIMGToken)
	if err != nil {
		return nil, err
	}

	return &api.EmptyStruct{}, nil
}

// RealnameDetail detail about realname by mid
func (s *MemberServer) RealnameDetail(ctx context.Context, req *api.MemberMidReq) (*api.RealnameDetailReply, error) {
	res, err := s.svr.RealnameDetail(ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	var realnameDetail = &api.RealnameDetailReply{
		Realname: res.Realname,
		Card:     res.Card,
		CardType: int8(res.CardType),
		Status:   int8(res.Status),
		Gender:   res.Gender,
		HandImg:  res.HandIMG,
	}

	return realnameDetail, nil
}

// RealnameStrippedInfo is
func (s *MemberServer) RealnameStrippedInfo(ctx context.Context, req *api.MemberMidReq) (*api.RealnameStrippedInfoReply, error) {
	return s.svr.RealnameStrippedInfo(ctx, req.Mid)
}

// MidByRealnameCard is
func (s *MemberServer) MidByRealnameCard(ctx context.Context, req *api.MidByRealnameCardsReq) (*api.MidByRealnameCardReply, error) {
	return s.svr.MidByRealnameCard(ctx, req)
}
