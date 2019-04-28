package service

import (
	"context"
	"time"

	"go-common/app/job/main/point/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	_timeFormat = "2006-01-02 15:04:05"
)

// AddPoint add point.
func (s *Service) AddPoint(c context.Context, p *model.VipPoint) (err error) {
	if _, err = s.dao.AddPoint(c, p); err != nil {
		log.Error("%+v", err)
		return
	}
	s.dao.DelPointInfoCache(c, p.Mid)
	return
}

// UpdatePoint update point.
func (s *Service) UpdatePoint(c context.Context, p *model.VipPoint, oldpoint *model.VipPoint) (err error) {
	if _, err = s.dao.UpdatePoint(c, p, oldpoint.Ver); err != nil {
		log.Error("%+v", err)
		return
	}
	s.dao.DelPointInfoCache(c, p.Mid)
	return
}

// AddPointHistory add point history.
func (s *Service) AddPointHistory(c context.Context, h *model.VipPointChangeHistoryMsg) (err error) {
	var (
		history    = new(model.VipPointChangeHistory)
		changeTime time.Time
	)
	history.ChangeType = h.ChangeType
	history.Mid = h.Mid
	history.Operator = h.Operator
	history.OrderID = h.OrderID
	history.Point = h.Point
	history.PointBalance = h.PointBalance
	history.RelationID = h.RelationID
	history.Remark = h.Remark
	if changeTime, err = time.ParseInLocation(_timeFormat, h.ChangeTime, time.Local); err != nil {
		log.Error("time.ParseInLocation error %+v", err)
		return
	}
	history.ChangeTime = xtime.Time(changeTime.Unix())
	if _, err = s.dao.AddPointHistory(c, history); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

// Notify notify.
func (s *Service) Notify(c context.Context, msg *model.VipPointChangeHistoryMsg) (err error) {
	for _, url := range s.c.Properties.NotifyCacheDelURL {
		if err = s.dao.NotifyCacheDel(c, url, msg.Mid, "127.0.0.1"); err != nil {
			log.Error("NotifyCacheDel fail(%d) %+v", msg.Mid, err)
		}
	}
	return
}
