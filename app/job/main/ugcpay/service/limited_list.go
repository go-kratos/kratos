package service

import (
	"context"
	"time"

	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/model"
	"go-common/library/log"
)

type limitedList interface {
	LimitSize() int
	BeginID(ctx context.Context) (id int64, err error)
	List(ctx context.Context, beginID int64) (maxID int64, list []interface{}, err error)
}

func runLimitedList(ctx context.Context, ll limitedList, sleep time.Duration, handler func(ctx context.Context, ele interface{}) error) (err error) {
	if ll.LimitSize() <= 0 {
		return
	}
	var (
		id   int64
		list = make([]interface{}, ll.LimitSize())
	)
	if id, err = ll.BeginID(ctx); err != nil {
		return
	}
	for len(list) >= ll.LimitSize() {
		if id, list, err = ll.List(ctx, id); err != nil {
			return
		}
		for _, ele := range list {
			if err = handler(ctx, ele); err != nil {
				log.Error("handle failed, ele: %+v, err: %+v", ele, err)
				err = nil
				continue
			}
			if sleep > 0 {
				time.Sleep(sleep)
			}
		}
	}
	return
}

type orderPaidLL struct {
	beginTime time.Time
	endTime   time.Time
	limit     int
	dao       *dao.Dao
}

func (o *orderPaidLL) LimitSize() int {
	return o.limit
}

func (o *orderPaidLL) BeginID(ctx context.Context) (id int64, err error) {
	return o.dao.MinIDOrderPaid(ctx, o.beginTime)
}

func (o *orderPaidLL) List(ctx context.Context, beginID int64) (maxID int64, list []interface{}, err error) {
	var rawList []*model.Order
	if maxID, rawList, err = o.dao.OrderPaidList(ctx, o.beginTime, o.endTime, beginID, o.limit); err != nil {
		return
	}
	log.Info("orderPaidLL beginID: %d, beginTime: %+v, endTime: %+v, limit: %d, size: %d", beginID, o.beginTime, o.endTime, o.limit, len(rawList))
	for _, r := range rawList {
		list = append(list, r)
	}
	return
}

type orderRefundedLL struct {
	beginTime time.Time
	endTime   time.Time
	limit     int
	dao       *dao.Dao
}

func (o *orderRefundedLL) LimitSize() int {
	return o.limit
}

func (o *orderRefundedLL) BeginID(ctx context.Context) (id int64, err error) {
	return o.dao.MinIDOrderRefunded(ctx, o.beginTime)
}

func (o *orderRefundedLL) List(ctx context.Context, beginID int64) (maxID int64, list []interface{}, err error) {
	var rawList []*model.Order
	if maxID, rawList, err = o.dao.OrderRefundedList(ctx, o.beginTime, o.endTime, beginID, o.limit); err != nil {
		return
	}
	log.Info("orderRefundedLL beginID: %d, beginTime: %+v, endTime: %+v, limit: %d, size: %d", beginID, o.beginTime, o.endTime, o.limit, len(rawList))
	for _, r := range rawList {
		list = append(list, r)
	}
	return
}

type dailyBillLLByVer struct {
	ver   int64
	limit int
	dao   *dao.Dao
}

func (d *dailyBillLLByVer) LimitSize() int {
	return d.limit
}

func (d *dailyBillLLByVer) BeginID(ctx context.Context) (id int64, err error) {
	return d.dao.MinIDDailyBillByVer(ctx, d.ver)
}

func (d *dailyBillLLByVer) List(ctx context.Context, beginID int64) (maxID int64, list []interface{}, err error) {
	var rawList []*model.DailyBill
	if maxID, rawList, err = d.dao.DailyBillListByVer(ctx, d.ver, beginID, d.limit); err != nil {
		return
	}
	log.Info("dailyBillLLByVer beginID: %d, ver: %d, limit: %d, size: %d", beginID, d.ver, d.limit, len(rawList))
	for _, r := range rawList {
		list = append(list, r)
	}
	return
}

type dailyBillLLByMonthVer struct {
	monthVer int64
	limit    int
	dao      *dao.Dao
}

func (d *dailyBillLLByMonthVer) LimitSize() int {
	return d.limit
}

func (d *dailyBillLLByMonthVer) BeginID(ctx context.Context) (id int64, err error) {
	return d.dao.MinIDDailyBillByMonthVer(ctx, d.monthVer)
}

func (d *dailyBillLLByMonthVer) List(ctx context.Context, beginID int64) (maxID int64, list []interface{}, err error) {
	var rawList []*model.DailyBill
	if maxID, rawList, err = d.dao.DailyBillListByMonthVer(ctx, d.monthVer, beginID, d.limit); err != nil {
		return
	}
	log.Info("dailyBillLLByMonthVer beginID: %d, monthVer: %d, limit: %d, size: %d", beginID, d.monthVer, d.limit, len(rawList))
	for _, r := range rawList {
		list = append(list, r)
	}
	return
}

type monthlyBillLL struct {
	ver   int64
	limit int
	dao   *dao.Dao
}

func (m *monthlyBillLL) LimitSize() int {
	return m.limit
}

func (m *monthlyBillLL) BeginID(ctx context.Context) (id int64, err error) {
	return m.dao.MinIDMonthlyBill(ctx, m.ver)
}

func (m *monthlyBillLL) List(ctx context.Context, beginID int64) (maxID int64, list []interface{}, err error) {
	var rawList []*model.Bill
	if maxID, rawList, err = m.dao.MonthlyBillList(ctx, m.ver, beginID, m.limit); err != nil {
		return
	}
	log.Info("monthlyBillLL beginID: %d, ver: %d, limit: %d, size: %d", beginID, m.ver, m.limit, len(rawList))
	for _, r := range rawList {
		list = append(list, r)
	}
	return
}
