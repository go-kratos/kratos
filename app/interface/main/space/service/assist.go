package service

import (
	"context"

	"go-common/app/service/main/assist/model/assist"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var _emptyAssists = make([]*assist.AssistUp, 0)

// RiderList get rider list by mid
func (s *Service) RiderList(c context.Context, mid int64, pn, ps int) (res *assist.AssistUpsPager, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &assist.ArgAssistUps{AssistMid: mid, Pn: int64(pn), Ps: int64(ps), RealIP: ip}
	if res, err = s.ass.AssistUps(c, arg); err != nil {
		log.Error("s.ass.AssistUps(%d,%d,%d) error(%v)", mid, pn, ps, err)
	}
	if len(res.Data) == 0 {
		res.Data = _emptyAssists
	}
	return
}

// ExitRider del rider with mid and upMid
func (s *Service) ExitRider(c context.Context, mid, upMid int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if err = s.ass.AssistExit(c, &assist.ArgAssist{Mid: upMid, AssistMid: mid, RealIP: ip}); err != nil {
		log.Error("s.add.DelAssist(%d,%d) error(%v)", mid, upMid, err)
	}
	return
}
