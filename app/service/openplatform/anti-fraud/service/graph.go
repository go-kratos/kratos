package service

import (
	"context"
	gtm "go-common/app/common/openplatform/geetest/model"
	gtsvc "go-common/app/common/openplatform/geetest/service"
	"go-common/app/service/openplatform/anti-fraud/api/grpc/v1"
	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"
)

//GraphPrepare 拉起图片验证
func (s *Service) GraphPrepare(c context.Context, req *v1.GraphPrepareRequest) (res *v1.GraphPrepareResponse, err error) {
	go s.d.IncrGeetestCount(context.Background())
	var data *gtm.ProcessRes
	if data, err = gtsvc.PreProcess(c, req.MID, 1, req.IP, req.ClientType, model.GeetestCaptchaID, model.GeetestPrivateKey); err != nil {
		return
	}
	res = new(v1.GraphPrepareResponse)
	res.Success = int64(data.Success)
	res.CaptchaID = data.CaptchaID
	res.Challenge = data.Challenge
	res.NewCaptcha = int64(data.NewCaptcha)
	res.Voucher = s.d.Voucher(c, req.MID, req.IP, req.ItemID, req.Customer, model.VoucherTypePull)
	return
}

//GraphCheck 图形验证
func (s *Service) GraphCheck(c context.Context, req *v1.GraphCheckRequest) (res *v1.GraphCheckResponse, err error) {
	if err = s.d.CheckVoucher(c, req.MID, req.Voucher, model.VoucherTypePull); err != nil {
		return
	}
	if status := gtsvc.Validate(c, req.Challenge, req.Validate, req.Seccode, req.ClientType, req.IP, model.GeetestCaptchaID, model.GeetestPrivateKey, int(req.Success), req.MID); !status {
		err = ecode.AntiValidateFailed
		return
	}
	res = new(v1.GraphCheckResponse)
	res.NewVoucher = s.d.Voucher(c, req.MID, req.IP, req.ItemID, req.Customer, model.VoucherTypeCheck)
	return res, nil
}
