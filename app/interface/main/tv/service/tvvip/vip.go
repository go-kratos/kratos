package tvvip

import (
	"context"
	tvmdl "go-common/app/interface/main/tv/model/tvvip"
	pb "go-common/app/service/main/tv/api"
)

const (
	ystSystemError = "999"
)

// VipInfo implementation
func (s *Service) VipInfo(ctx context.Context, mid int64) (resp *pb.UserInfoReply, err error) {
	return s.tvVipClient.UserInfo(ctx, &pb.UserInfoReq{Mid: mid})
}

func (s *Service) YstVipInfo(ctx context.Context, mid int64, sign string) (resp *pb.YstUserInfoReply, err error) {
	return s.tvVipClient.YstUserInfo(ctx, &pb.YstUserInfoReq{Mid: mid, Sign: sign})
}

// ChangeHistory implementation
func (s *Service) ChangeHistory(ctx context.Context, id int32) (resp *pb.ChangeHistoryReply, err error) {
	return s.tvVipClient.ChangeHistory(ctx, &pb.ChangeHistoryReq{Id: id})
}

// ChangeHistorys implementation
func (s *Service) ChangeHistorys(ctx context.Context, mid int64, from, to, pn, ps int32) (resp *pb.ChangeHistorysReply, err error) {
	return s.tvVipClient.ChangeHistorys(ctx, &pb.ChangeHistorysReq{Mid: mid, From: from, To: to, Pn: pn, Ps: ps})
}

// PanelInfo implemention
func (s *Service) PanelInfo(ctx context.Context, mid int64) (resp *pb.PanelInfoReply, err error) {
	resp, err = s.tvVipClient.PanelInfo(ctx, &pb.PanelInfoReq{Mid: mid})
	return
}

// GuestPanelInfo implemention
func (s *Service) GuestPanelInfo(ctx context.Context) (resp *pb.GuestPanelInfoReply, err error) {
	return s.tvVipClient.GuestPanelInfo(ctx, &pb.GuestPanelInfoReq{})
}

// CreateQr implemention
func (s *Service) CreateQr(ctx context.Context, req *tvmdl.CreateQrReq) (resp *pb.CreateQrReply, err error) {
	pr := new(pb.CreateQrReq)
	req.CopyIntoPbCreateOrReq(pr)
	return s.tvVipClient.CreateQr(ctx, pr)
}

// CreateGuestQr implemention
func (s *Service) CreateGuestQr(ctx context.Context, req *tvmdl.CreateGuestQrReq) (resp *pb.CreateGuestQrReply, err error) {
	pr := new(pb.CreateGuestQrReq)
	req.CopyIntoPbCreateGuestQrReq(pr)
	return s.tvVipClient.CreateGuestQr(ctx, pr)
}

// TokenInfo implemention
func (s *Service) TokenInfo(ctx context.Context, tokens []string) (resp *pb.TokenInfoReply, err error) {
	req := &pb.TokenInfoReq{
		Token: tokens,
	}
	return s.tvVipClient.TokenInfo(ctx, req)
}

// CreateOrder implementation
func (s *Service) CreateOrder(ctx context.Context, clientIp string, req *tvmdl.CreateOrderReq) (resp *pb.CreateOrderReply, err error) {
	pr := new(pb.CreateOrderReq)
	req.CopyIntoPbCreateOrderReq(pr)
	return s.tvVipClient.CreateOrder(ctx, pr)
}

// CreateGuestOrder implementation
func (s *Service) CreateGuestOrder(ctx context.Context, mid int64, clientIp string, req *tvmdl.CreateGuestOrderReq) (resp *pb.CreateGuestOrderReply, err error) {
	pr := new(pb.CreateGuestOrderReq)
	pr.Mid = mid
	req.CopyIntoPbCreateGuestOrderReq(pr)
	return s.tvVipClient.CreateGuestOrder(ctx, pr)
}

// PayCallback implementation
func (s *Service) PayCallback(ctx context.Context, req *tvmdl.YstPayCallbackReq) (resp *pb.PayCallbackReply) {
	var err error
	pr := new(pb.PayCallbackReq)
	req.CopyIntoPbPayCallbackReq(pr)
	resp, err = s.tvVipClient.PayCallback(ctx, pr)
	if err != nil {
		resp = new(pb.PayCallbackReply)
		resp.Result = ystSystemError
		resp.Msg = err.Error()
	}
	return
}

// ContractCallback implementation
func (s *Service) WxContractCallback(ctx context.Context, req *tvmdl.WxContractCallbackReq) (resp *pb.WxContractCallbackReply) {
	var err error
	wc := new(pb.WxContractCallbackReq)
	req.CopyIntoPbWxContractCallbackReq(wc)
	resp, err = s.tvVipClient.WxContractCallback(ctx, wc)
	if err != nil {
		resp = new(pb.WxContractCallbackReply)
		resp.Result = ystSystemError
		resp.Msg = err.Error()
	}
	return
}
