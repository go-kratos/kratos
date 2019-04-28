// Package server generate by warden_gen
package server

import (
	"context"

	v1 "go-common/app/service/main/coupon/api"
	"go-common/app/service/main/coupon/model"
	"go-common/app/service/main/coupon/service"
	"go-common/library/net/rpc/warden"
)

// New VipInfo warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	v1.RegisterCouponServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.CouponServer = &server{}

func (s *server) CaptchaToken(c context.Context, req *v1.CaptchaTokenReq) (res *v1.CaptchaTokenReply, err error) {
	var token *model.Token
	if token, err = s.svr.CaptchaToken(c, req.Ip); err != nil || token == nil {
		return
	}
	return &v1.CaptchaTokenReply{
		Token: token.Token,
		Url:   token.URL,
	}, err
}

func (s *server) UseCouponCode(c context.Context, req *v1.UseCouponCodeReq) (res *v1.UseCouponCodeResp, err error) {
	var data *model.UseCouponCodeResp
	if data, err = s.svr.UseCouponCode(c, &model.ArgUseCouponCode{
		Token:  req.Token,
		Code:   req.Code,
		Verify: req.Verify,
		IP:     req.Ip,
		Mid:    req.Mid,
	}); err != nil || data == nil {
		return
	}
	return &v1.UseCouponCodeResp{
		CouponToken:          data.CouponToken,
		CouponAmount:         data.CouponAmount,
		FullAmount:           data.FullAmount,
		PlatfromLimitExplain: data.PlatfromLimitExplain,
		ProductLimitMonth:    data.ProductLimitMonth,
		ProductLimitRenewal:  data.ProductLimitRenewal,
	}, err
}

func (s *server) UsableAllowanceCouponV2(c context.Context, req *v1.UsableAllowanceCouponV2Req) (res *v1.UsableAllowanceCouponV2Reply, err error) {
	var (
		data *model.CouponTipInfo
		ci   *v1.ModelCouponAllowancePanelInfo
	)
	if data, err = s.svr.UsableAllowanceCouponV2(c, req); err != nil {
		return
	}
	if data.CouponInfo != nil {
		ci = &v1.ModelCouponAllowancePanelInfo{
			CouponToken:         data.CouponInfo.CouponToken,
			CouponAmount:        data.CouponInfo.Amount,
			State:               data.CouponInfo.State,
			FullAmount:          data.CouponInfo.FullAmount,
			FullLimitExplain:    data.CouponInfo.FullLimitExplain,
			ScopeExplain:        data.CouponInfo.ScopeExplain,
			CouponDiscountPrice: data.CouponInfo.CouponDiscountPrice,
			StartTime:           data.CouponInfo.StartTime,
			ExpireTime:          data.CouponInfo.ExpireTime,
			Selected:            int32(data.CouponInfo.Selected),
			DisablesExplains:    data.CouponInfo.DisablesExplains,
			OrderNo:             data.CouponInfo.OrderNO,
			Name:                data.CouponInfo.Name,
			Usable:              int32(data.CouponInfo.Usable),
		}
	}
	return &v1.UsableAllowanceCouponV2Reply{
		CouponTip:  data.CouponTip,
		CouponInfo: ci,
	}, err
}
