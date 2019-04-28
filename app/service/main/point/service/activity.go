package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/point/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) activityGiveTimes(c context.Context, mid int64, changeType int, point int64) (sendTime int64, err error) {
	var (
		startDate   time.Time
		endDate     time.Time
		obtain      int64
		buyPoint    int64
		amount      int64
		alreadySend int64
		phs         []*model.PointHistory
		now         = time.Now()
		ok          bool
	)
	if changeType != model.Contract {
		return
	}
	if startDate, err = time.ParseInLocation(_timeFormat, s.c.Property.PointActiveStartDate, time.Local); err != nil {
		log.Error("time.ParseInLocation %+v", err)
		return
	}
	if endDate, err = time.ParseInLocation(_timeFormat, s.c.Property.PointActiveEndDate, time.Local); err != nil {
		log.Error("time.ParseInLocation %+v", err)
		return
	}
	if now.Before(startDate) || now.After(endDate) {
		return
	}
	if obtain, ok = s.c.Property.PointGetRule[strconv.Itoa(changeType)]; !ok || obtain == 0 {
		log.Info("not found rete by changeType %+v", changeType)
		return
	}
	if phs, err = s.dao.SelPointHistory(c, mid, xtime.Time(startDate.Unix()), xtime.Time(endDate.Unix())); err != nil {
		log.Error("%+v", err)
		return
	}
	for _, v := range phs {
		if v.ChangeType == model.Contract {
			buyPoint += v.Point
			// old system pgc remark TNT.
			if model.ActivityGiveRemark == v.Remark {
				alreadySend += v.Point
			}
		}
	}
	buyPoint += point
	amount = buyPoint / obtain
	if alreadySend == 0 {
		if amount >= model.ActivityMixBuyBp && amount < model.ActivityOutOfBuyBp {
			sendTime = model.ActivitySendTimes1
		} else if amount >= model.ActivityOutOfBuyBp {
			sendTime = model.ActivitySendTimes2
		}
	} else if alreadySend == model.ActivityGivePoint && amount >= model.ActivityOutOfBuyBp {
		sendTime = model.ActivitySendTimes1
	}
	return
}

func (s *Service) activeSendPoint(c context.Context, tx *xsql.Tx, ph *model.PointHistory) (activePoint int64, err error) {
	var (
		sendTime     int64
		pointBalance int64
	)
	sendTime, err = s.activityGiveTimes(c, ph.Mid, ph.ChangeType, ph.Point)
	for i := 0; i < int(sendTime); i++ {
		activePoint += model.ActivityGivePoint
		phAdd := new(model.PointHistory)
		if pointBalance, err = s.updatePoint(c, tx, ph.Mid, model.ActivityGivePoint); err != nil {
			log.Error("%+v", err)
			activePoint = 0
			return
		}
		phAdd.Point = model.ActivityGivePoint
		phAdd.PointBalance = pointBalance
		phAdd.ChangeType = model.PointSystem
		phAdd.RelationID = ph.OrderID
		phAdd.Remark = model.ActivityGiveRemark
		phAdd.Mid = ph.Mid
		phAdd.ChangeTime = xtime.Time(time.Now().Unix())
		s.dao.InsertPointHistory(c, tx, phAdd)
	}
	log.Info("send point total->%v  mid:%v orderID:%v", activePoint, ph.Mid, ph.OrderID)
	return
}
