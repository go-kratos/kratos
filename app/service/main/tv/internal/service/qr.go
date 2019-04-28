package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/tv/internal/model"
	xtime "go-common/library/time"
)

func (s *Service) MakePayParam(c context.Context, mid int64, pid int32, buyNum int32, guid string, appChannel string) (p *model.PayParam) {
	return &model.PayParam{
		Mid:        mid,
		Pid:        pid,
		OrderNo:    s.makeOrderNo(),
		BuyNum:     buyNum,
		Guid:       guid,
		AppChannel: appChannel,
		Status:     model.PayOrderStatusPending,
		ExpireAt:   xtime.Time(time.Now().Unix() + int64(s.c.CacheTTL.PayParamTTL)),
	}
}

func (s *Service) CreateQr(c context.Context, mid int64, pid int32, buyNum int32, guid string, appChannel string) (qr *model.QR, err error) {
	payParam := s.MakePayParam(c, mid, pid, buyNum, guid, appChannel)
	token := payParam.MD5()
	qr = &model.QR{
		ExpireAt: payParam.ExpireAt,
		Token:    token,
		URL:      fmt.Sprintf("%s?token=%s", s.c.PAY.QrURL, token),
	}
	s.dao.AddCachePayParam(c, token, payParam)
	return qr, nil
}

func (s *Service) CreateGuestQr(c context.Context, pid int32, buyNum int32, guid string, appChannel string) (qr *model.QR, err error) {
	payParam := s.MakePayParam(c, -1, pid, buyNum, guid, appChannel)
	token := payParam.MD5()
	qr = &model.QR{
		ExpireAt: payParam.ExpireAt,
		Token:    token,
		URL:      fmt.Sprintf("%s?token=%s", s.c.PAY.GuestQrURL, token),
	}
	s.dao.AddCachePayParam(c, token, payParam)
	return qr, nil
}
