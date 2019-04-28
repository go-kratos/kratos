package coupon

import (
	"context"

	"go-common/app/interface/main/account/conf"
	v1 "go-common/app/service/main/coupon/api"
	"go-common/app/service/main/coupon/model"
	courpc "go-common/app/service/main/coupon/rpc/client"
)

// Service .
type Service struct {
	// conf
	c *conf.Config
	// rpc
	couRPC *courpc.Service
	// coupon grpc service
	coupongRPC v1.CouponClient
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		couRPC: courpc.New(c.RPCClient2.Coupon),
	}
	coupongRPC, err := v1.NewClient(c.CouponClient)
	if err != nil {
		panic(err)
	}
	s.coupongRPC = coupongRPC
	return
}

// AllowanceList allowance list.
func (s *Service) AllowanceList(c context.Context, mid int64, state int8) (res []*model.CouponAllowancePanelInfo, err error) {
	res, err = s.couRPC.AllowanceList(c, &model.ArgAllowanceList{Mid: mid, State: state})
	return
}

// CouponPage coupon list.
func (s *Service) CouponPage(c context.Context, a *model.ArgRPCPage) (res *model.CouponPageRPCResp, err error) {
	res, err = s.couRPC.CouponPage(c, a)
	return
}

// CouponCartoonPage coupon cartoon list.
// func (s *Service) CouponCartoonPage(c context.Context, a *model.ArgRPCPage) (res *model.CouponCartoonPageResp, err error) {
// 	res, err = s.couRPC.CouponCartoonPage(c, a)
// 	return
// }

// PrizeCards .
func (s *Service) PrizeCards(c context.Context, a *model.ArgCount) (res []*model.PrizeCardRep, err error) {
	res, err = s.couRPC.PrizeCards(c, a)
	return
}

// PrizeDraw .
func (s *Service) PrizeDraw(c context.Context, a *model.ArgPrizeDraw) (res *model.PrizeCardRep, err error) {
	res, err = s.couRPC.PrizeDraw(c, a)
	return
}

// CaptchaToken captcha token.
func (s *Service) CaptchaToken(c context.Context, a *v1.CaptchaTokenReq) (res *v1.CaptchaTokenReply, err error) {
	return s.coupongRPC.CaptchaToken(c, a)
}

// UseCouponCode use coupon code.
func (s *Service) UseCouponCode(c context.Context, a *model.ArgUseCouponCode) (res *v1.UseCouponCodeResp, err error) {
	return s.coupongRPC.UseCouponCode(c, &v1.UseCouponCodeReq{
		Token:  a.Token,
		Code:   a.Code,
		Verify: a.Verify,
		Ip:     a.IP,
		Mid:    a.Mid,
	})
}
