package dao

import (
	"context"

	accmdl "go-common/app/service/main/account/model"
	spymdl "go-common/app/service/main/spy/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// UserScore get userscore by id , will init score if score not exist.
func (d *Dao) UserScore(c context.Context, mid int64) (userScore *spymdl.UserScore, err error) {
	var argUserScore = &spymdl.ArgUserScore{
		Mid: mid,
		IP:  metadata.String(c, metadata.RemoteIP),
	}
	if userScore, err = d.spyRPC.UserScore(c, argUserScore); err != nil {
		log.Error("dao.spyRPC.UserScore(%v) error(%v)", argUserScore, err)
		return
	}
	return
}

// AccInfo get account info by mid
func (d *Dao) AccInfo(c context.Context, mid int64) (ai *accmdl.Info, err error) {
	arg := &accmdl.ArgMid{Mid: mid}
	if ai, err = d.accRPC.Info3(c, arg); err != nil || ai == nil {
		log.Error("s.accRPC.Info(%d) error(%v)", mid, err)
	}
	return
}

// ResetBase reset user base score.
func (d *Dao) ResetBase(c context.Context, mid int64, operator string) (err error) {
	arg := &spymdl.ArgReset{
		Mid:       mid,
		BaseScore: true,
		Operator:  operator,
	}
	if err = d.spyRPC.UpdateBaseScore(c, arg); err != nil {
		log.Error("s.spyRPC.UpdateBaseScore(%v) error(%v)", arg, err)
		return
	}
	return
}

// RefreshBase reset user base score.
func (d *Dao) RefreshBase(c context.Context, mid int64, operator string) (err error) {
	arg := &spymdl.ArgReset{
		Mid:       mid,
		BaseScore: true,
		Operator:  operator,
	}
	if err = d.spyRPC.RefreshBaseScore(c, arg); err != nil {
		log.Error("s.spyRPC.RefreshBaseScore(%v) error(%v)", arg, err)
		return
	}
	return
}

// ResetEvent reset user event score.
func (d *Dao) ResetEvent(c context.Context, mid int64, operator string) (err error) {
	arg := &spymdl.ArgReset{
		Mid:        mid,
		EventScore: true,
		Operator:   operator,
	}
	if err = d.spyRPC.UpdateEventScore(c, arg); err != nil {
		log.Error("s.spyRPC.UpdateEventScore(%v) error(%v)", arg, err)
		return
	}
	return
}

// ClearCount clear count.
func (d *Dao) ClearCount(c context.Context, mid int64, operator string) (err error) {
	arg := &spymdl.ArgReset{
		Mid:        mid,
		ReLiveTime: true,
		Operator:   operator,
	}
	if err = d.spyRPC.ClearReliveTimes(c, arg); err != nil {
		log.Error("d.spyRPC.ClearReliveTimes(%v) error(%v)", arg, err)
		return
	}
	return
}
