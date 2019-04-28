package service

import (
	"context"

	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) checkFollow(c context.Context, mid, assistMid int64) (err error) {
	follow, err := s.acc.IsFollow(c, mid, assistMid)
	if err != nil {
		log.Error("s.ass.IsFollow(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if !follow {
		log.Error("s.ass.IsFollow AssistNotFollowUp(%d,%d) error(%v)", mid, assistMid, err)
		err = ecode.AssistNotFollowUp
		return
	}
	return
}

// checkIdentify func
func (s *Service) checkIdentify(c context.Context, assistMid int64) (err error) {
	if err = s.acc.IdentifyInfo(c, assistMid, ""); err != nil {
		log.Error("s.acc.IdentifyInfo IdentifyInfoFailed assistMid(%d)", assistMid)
		return
	}
	return
}

func (s *Service) checkBanned(c context.Context, assistMid int64) (err error) {
	if err = s.acc.UserBanned(c, assistMid); err != nil {
		log.Error("s.UserBanned err: (%d) error(%v)", assistMid, err)
		return
	}
	return
}

func (s *Service) checkIsAssist(c context.Context, mid, assistMid int64) (err error) {
	assist, err := s.ass.Assist(c, mid, assistMid)
	if err != nil {
		log.Error("s.ass.Assist(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if assist != nil {
		log.Error("s.ass.Assist(%d,%d) assist is not nil error(%v)", mid, assistMid, err)
		err = ecode.AssistAlreadyExist
		return
	}
	return
}

func (s *Service) checkIsNotAssist(c context.Context, mid, assistMid int64) (err error) {
	assist, err := s.ass.Assist(c, mid, assistMid)
	if err != nil {
		log.Error("s.ass.Assist(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if assist == nil {
		err = ecode.AssistNotExist
		log.Error("s.ass.Assist(%d,%d) assist is nil error(%v)", mid, assistMid, err)
		return
	}
	return
}
