package service

import (
	"context"
	"time"

	"go-common/app/service/main/push/dao"
	"go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) loadBusinessproc() {
	for {
		if s.loadBusiness() != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		time.Sleep(time.Duration(s.c.Push.LoadBusinessInteval))
	}
}

func (s *Service) loadBusiness() (err error) {
	res, err := s.dao.Businesses(context.Background())
	if err != nil {
		log.Error("s.dao.Business() error(%v)", err)
		return
	}
	if len(res) > 0 {
		s.businesses = res
	}
	return
}

func (s *Service) checkBusiness(id int64, token string) error {
	b, ok := s.businesses[id]
	if !ok {
		log.Error("business is not exist. business(%d) token(%s)", id, token)
		dao.PromError("service:业务方不存在")
		return ecode.PushBizAuthErr
	}
	if token != b.Token {
		log.Error("wrong token business(%d) token(%s) need(%s)", id, token, b.Token)
		dao.PromError("service:业务方token错误")
		return ecode.PushBizAuthErr
	}
	if b.PushSwitch == model.SwitchOff {
		log.Error("business was forbidden. business(%d) token(%s)", id, token)
		dao.PromError("service:业务方被禁止推送")
		return ecode.PushBizForbiddenErr
	}
	// 校验免打扰时间
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	if (hour >= b.SilentTime.BeginHour && minute >= b.SilentTime.BeginMinute) &&
		(hour <= b.SilentTime.EndHour && minute <= b.SilentTime.EndMinute) {
		log.Warn("in silent time, forbidden. business(%d) now(%v)", id, time.Now())
		return ecode.PushSilenceErr
	}
	return nil
}
