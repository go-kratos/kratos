package service

import (
	"context"
	"runtime/debug"
	"time"

	"go-common/app/job/main/ugcpay/model"
	"go-common/library/log"
)

// SyncElecOrderList 同步老充电订单
func (s *Service) SyncElecOrderList(c context.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("syncElecOrderSync panic(%+v) :\n %s", x, debug.Stack())
		}
	}()
	var (
		limit   = 100
		list    = make([]*model.DBOldElecPayOrder, limit)
		startID = int64(0)
		err     error
	)
	log.Info("Start syncElecOrderSync from elec_pay_order")

	for len(list) >= limit {
		log.Info("sync progress elec_pay_order fromID (%d)", startID)
		// 1. load old data
		if startID, list, err = s.dao.OldElecOrderList(_ctx, startID, limit); err != nil {
			log.Error("%+v", err)
			return
		}
		// 2. save new data
		for _, ele := range list {
			if err = s.handleElecOrder(_ctx, ele); err != nil {
				log.Error("s.handleElecOrder: %+v, err: %+v", ele, err)
				return
			}
		}
		// 3. give db a break time
		time.Sleep(time.Millisecond * 20)
	}

	log.Info("End syncElecOrderSync from elec_pay_order")
}

// SyncElecMessageList 同步老充电留言
func (s *Service) SyncElecMessageList(c context.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("syncElecMessageList panic(%+v) :\n %s", x, debug.Stack())
		}
	}()
	var (
		limit   = 100
		list    = make([]*model.DBOldElecMessage, limit)
		startID = int64(0)
		err     error
	)
	log.Info("Start syncElecMessageList from elec_message")

	for len(list) >= limit {
		log.Info("sync progress elec_message fromID (%d)", startID)
		// 1. load old data
		if startID, list, err = s.dao.OldElecMessageList(_ctx, startID, limit); err != nil {
			log.Error("%+v", err)
			return
		}
		// 2. save new data
		for _, ele := range list {
			if err = s.handleOldElecMessage(_ctx, ele); err != nil {
				log.Error("s.handleOldElecMessage: %+v, err: %+v", ele, err)
				return
			}
		}
		// 3. give db a break time
		time.Sleep(time.Millisecond * 20)
	}

	log.Info("End syncElecMessageList from elec_message")
}
