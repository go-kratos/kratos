package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/library/ecode"
)

// SyncOrder 同步订单方法
func (s *Service) SyncOrder(c context.Context, oi *model.DistOrderArg) (affect int64, err error) {
	if oi.Stat == model.DistOrderNormal {
		has, _ := s.dao.HasOrder(c, oi)
		if has {
			err = ecode.TicketRecordDupli
			errors.Wrap(err, fmt.Sprintf("重复的正向订单:%v", oi.Oid))
			return
		}
	}
	if oi.Stat == model.DistOrderRefunded {
		roi := &model.DistOrderArg{}
		roi.Stat = model.DistOrderNormal
		roi.Oid = oi.Oid
		has, _ := s.dao.HasOrder(c, roi)
		if !has {
			err = ecode.TicketRecordLost
			errors.Wrap(err, fmt.Sprintf("缺少正向订单:%v", oi.Oid))
			return
		}
	}
	affect, err = s.dao.InsertOrder(c, oi)
	return
}

//GetOrder 获取分销订单方法
func (s *Service) GetOrder(c context.Context, oid uint64) (oi []*model.OrderInfo, err error) {
	oi, err = s.dao.GetOrder(c, oid)
	if err != nil {
		errors.Wrap(err, fmt.Sprintf("查询订单失败:%v", oid))
		return
	}
	return
}
