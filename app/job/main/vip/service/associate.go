package service

import (
	"context"
	"time"

	"go-common/app/job/main/vip/model"
	v1 "go-common/app/service/main/vip/api"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/log"
)

func (s *Service) eleEompensateJob() {
	log.Info("ele grant eompensate job start..................")
	if succeed := s.dao.AddTransferLock(context.TODO(), "lock:elegrant"); succeed {
		if err := s.EleGrantCompensate(context.TODO()); err != nil {
			log.Error("error(%+v)", err)
		}
	}
	log.Info("ele grant eompensate job end..................")
}

// EleGrantCompensate ele frant compensate.
func (s *Service) EleGrantCompensate(c context.Context) (err error) {
	var res []*model.VipOrderActivityRecord
	if res, err = s.dao.NotGrantActOrders(c, vipmol.PanelTypeEle, s.c.Property.NotGrantLimit); err != nil {
		return
	}
	for _, v := range res {
		if _, err = s.vipgRPC.EleVipGrant(c, &v1.EleVipGrantReq{OrderNo: v.OrderNo}); err != nil {
			log.Error("EleVipGrant a(%s) err(%+v)", v.OrderNo, err)
			continue
		}
		time.Sleep(time.Second)
	}
	return
}
