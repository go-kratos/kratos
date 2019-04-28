package grpc

import (
	"context"
	pb "go-common/app/service/main/tv/api"
	"go-common/app/service/main/tv/internal/model"
	"go-common/app/service/main/tv/internal/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

// New new warden rpc server
func New(c *warden.ServerConfig, svc *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	pb.RegisterTVServiceServer(ws.Server(), &server{svc})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	log.Info("start grpc server")
	return ws
}

type server struct {
	svr *service.Service
}

var _ pb.TVServiceServer = &server{}

// UserInfo implementation
func (s *server) UserInfo(ctx context.Context, req *pb.UserInfoReq) (resp *pb.UserInfoReply, err error) {
	ui, err := s.svr.UserInfo(ctx, req.Mid)
	if err != nil {
		return
	}
	if ui == nil {
		return nil, ecode.NothingFound
	}
	resp = &pb.UserInfoReply{}
	resp.DeepCopyFromUserInfo(ui)
	return
}

// ChangeHistory implementation
func (s *server) ChangeHistory(ctx context.Context, req *pb.ChangeHistoryReq) (resp *pb.ChangeHistoryReply, err error) {
	ch, err := s.svr.ChangeHistory(ctx, req.Id)
	if err != nil {
		return
	}
	if ch == nil {
		return nil, ecode.NothingFound
	}
	resp = &pb.ChangeHistoryReply{}
	resp.DeepCopyFromUserChangeHistory(ch)
	return
}

// ChangeHistorys implementation
func (s *server) ChangeHistorys(ctx context.Context, req *pb.ChangeHistorysReq) (resp *pb.ChangeHistorysReply, err error) {
	chs, total, err := s.svr.ChangeHistorys(ctx, req.Mid, req.From, req.To, req.Pn, req.Ps)
	if err != nil {
		return
	}
	resp = &pb.ChangeHistorysReply{Total: int32(total)}
	resp.Historys = make([]*pb.ChangeHistoryReply, 0, len(chs))
	for _, ch := range chs {
		chr := &pb.ChangeHistoryReply{}
		chr.DeepCopyFromUserChangeHistory(ch)
		resp.Historys = append(resp.Historys, chr)
	}
	return
}

func suitType2String(st int8) string {
	switch st {
	case model.SuitTypeAll:
		return "ALL"
	case model.SuitTypeMvip:
		return "MVIP"
	default:
		return "ALL"
	}
}

// PanelInfo implemention
func (s *server) PanelInfo(ctx context.Context, req *pb.PanelInfoReq) (resp *pb.PanelInfoReply, err error) {
	pi, err := s.svr.PanelInfo(ctx, req.Mid)
	if err != nil {
		return
	}
	resp = &pb.PanelInfoReply{}
	resp.PriceConfigs = make(map[string]*pb.PanelPriceConfigs)
	for st, ps := range pi {
		ppcs := &pb.PanelPriceConfigs{}
		ppcs.PriceConfigs = make([]*pb.PanelPriceConfig, 0)
		for _, p := range ps {
			item := &pb.PanelPriceConfig{}
			item.DeepCopyFromPanelPriceConfig(p)
			ppcs.PriceConfigs = append(ppcs.PriceConfigs, item)
		}
		resp.PriceConfigs[suitType2String(st)] = ppcs
	}
	return
}

// GuestPanelInfo implemention
func (s *server) GuestPanelInfo(ctx context.Context, req *pb.GuestPanelInfoReq) (resp *pb.GuestPanelInfoReply, err error) {
	pi, err := s.svr.GuestPanelInfo(ctx)
	if err != nil {
		return
	}
	resp = &pb.GuestPanelInfoReply{}
	resp.PriceConfigs = make(map[string]*pb.PanelPriceConfigs)
	for st, ps := range pi {
		ppcs := &pb.PanelPriceConfigs{}
		ppcs.PriceConfigs = make([]*pb.PanelPriceConfig, 0, len(ps))
		for _, p := range ps {
			item := &pb.PanelPriceConfig{}
			item.DeepCopyFromPanelPriceConfig(p)
			ppcs.PriceConfigs = append(ppcs.PriceConfigs, item)
		}
		resp.PriceConfigs[suitType2String(st)] = ppcs
	}
	return
}

// PayOrder implementation
func (s *server) PayOrder(ctx context.Context, req *pb.PayOrderReq) (resp *pb.PayOrderReply, err error) {
	resp = &pb.PayOrderReply{}
	return
}

// CreateQr implementation
func (s *server) CreateQr(ctx context.Context, req *pb.CreateQrReq) (resp *pb.CreateQrReply, err error) {
	qr, err := s.svr.CreateQr(ctx, req.Mid, req.Pid, req.BuyNum, req.Guid, req.AppChannel)
	if err != nil {
		return
	}
	resp = &pb.CreateQrReply{}
	resp.DeepCopyFromQR(qr)
	return
}

// CreateGuestQr implementation
func (s *server) CreateGuestQr(ctx context.Context, req *pb.CreateGuestQrReq) (resp *pb.CreateGuestQrReply, err error) {
	qr, err := s.svr.CreateGuestQr(ctx, req.Pid, req.BuyNum, req.Guid, req.AppChannel)
	if err != nil {
		return
	}
	resp = &pb.CreateGuestQrReply{}
	resp.DeepCopyFromQR(qr)
	return
}

// CreateOrder implementation
func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (resp *pb.CreateOrderReply, err error) {
	pi, err := s.svr.CreateOrder(ctx, req.Token, req.Platform, req.PaymentType, req.ClientIp)
	if err != nil {
		return
	}
	resp = &pb.CreateOrderReply{}
	resp.DeepCopyFromPayInfo(pi)
	return
}

// CreateGuestOrder implementation
func (s *server) CreateGuestOrder(ctx context.Context, req *pb.CreateGuestOrderReq) (resp *pb.CreateGuestOrderReply, err error) {
	pi, err := s.svr.CreateGuestOrder(ctx, req.Mid, req.Token, req.Platform, req.PaymentType, req.ClientIp)
	if err != nil {
		return
	}
	resp = &pb.CreateGuestOrderReply{}
	resp.DeepCopyFromPayInfo(pi)
	return
}

// RenewVip implementation
func (s *server) TokenInfo(ctx context.Context, req *pb.TokenInfoReq) (resp *pb.TokenInfoReply, err error) {
	ti, err := s.svr.TokenInfos(ctx, req.Token)
	if err != nil {
		return
	}
	resp = &pb.TokenInfoReply{}
	resp.Tokens = make([]*pb.TokenInfo, 0)
	for _, v := range ti {
		t := &pb.TokenInfo{}
		t.DeepCopyFromTokenInfo(v)
		resp.Tokens = append(resp.Tokens, t)
	}
	return
}

// RenewVip implementation
func (s *server) RenewVip(ctx context.Context, req *pb.RenewVipReq) (resp *pb.RenewVipReply, err error) {
	err = s.svr.RenewVip(ctx, req.Mid)
	if err != nil {
		return
	}
	resp = &pb.RenewVipReply{}
	return
}

// YstUserInfo implementation
func (s *server) YstUserInfo(ctx context.Context, req *pb.YstUserInfoReq) (resp *pb.YstUserInfoReply, err error) {
	resp = &pb.YstUserInfoReply{}
	ui, err := s.svr.YstUserInfo(ctx, req.DeepCopyAsYstUserInfoReq())
	if err != nil {
		resp.Result = "998"
		resp.Msg = err.Error()
		return
	}
	resp.DeepCopyFromUserInfo(ui)
	resp.Result = "0"
	resp.Msg = "ok"
	return
}

// PayCallback implementation
func (s *server) PayCallback(ctx context.Context, req *pb.PayCallbackReq) (resp *pb.PayCallbackReply, err error) {
	ystReq := &model.YstPayCallbackReq{}
	req.DeepCopyAsIntoYstPayCallbackReq(ystReq)
	ystReply := s.svr.PayCallback(ctx, ystReq)
	resp = &pb.PayCallbackReply{}
	resp.DeepCopyFromYstPayCallbackReply(ystReply)
	return
}

// WxContractCallback implementation.
func (s *server) WxContractCallback(ctx context.Context, req *pb.WxContractCallbackReq) (resp *pb.WxContractCallbackReply, err error) {
	resp = &pb.WxContractCallbackReply{}
	return
}
