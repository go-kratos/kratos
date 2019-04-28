package service

import (
	"context"

	"go-common/app/service/main/passport-game/model"
)

// RegV3 RegV3
func (s *Service) RegV3(c context.Context, tdoRegV3 model.TdoRegV3) (regV3 *model.ResRegV3, err error) {
	regV3, err = s.d.RegV3(c, tdoRegV3)
	return
}

// RegV2 RegV2
func (s *Service) RegV2(c context.Context, tdoRegV2 model.TdoRegV2) (regV2 *model.ResRegV2, err error) {
	regV2, err = s.d.RegV2(c, tdoRegV2)
	return
}

// Reg Reg
func (s *Service) Reg(c context.Context, tdoReg model.TdoReg) (reg *model.ResReg, err error) {
	reg, err = s.d.Reg(c, tdoReg)
	return
}

// ByTel ByTel
func (s *Service) ByTel(c context.Context, tdoByTel model.TdoByTel) (byTel *model.ResByTel, err error) {
	byTel, err = s.d.ByTel(c, tdoByTel)
	return
}

// Captcha Captcha
func (s *Service) Captcha(c context.Context, ip string) (captchaData *model.CaptchaData, err error) {
	captchaData, err = s.d.Captcha(c, ip)
	return
}

// SendSms SendSms
func (s *Service) SendSms(c context.Context, tdoSendSms model.TdoSendSms) (err error) {
	err = s.d.SendSms(c, tdoSendSms)
	return
}
