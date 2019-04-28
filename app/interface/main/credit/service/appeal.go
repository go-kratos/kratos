package service

import (
	"context"
	"time"

	"go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// AddAppeal  add a new appeal .
func (s *Service) AddAppeal(c context.Context, btid, bid, mid int64, reason string) (err error) {
	var (
		isID   bool
		origin string
		ctime  xtime.Time
		caseID int64
	)
	infos, err := s.BlockedUserList(c, mid)
	if err != nil {
		err = errors.Wrap(err, "s.dao.BlockedUserList error")
		return
	}
	for _, v := range infos {
		if bid == v.ID {
			isID = true
			origin = v.OriginContent
			ctime = v.CTime
			caseID = v.CaseID
		}
	}
	if !isID {
		err = ecode.CreditBlockNotExist
		return
	}
	if xtime.Time(time.Now().AddDate(0, 0, -7).Unix()) > ctime {
		err = ecode.CreditBlockExpired
		return
	}
	if err = s.dao.AddAppeal(c, s.tagMap[int8(btid)], btid, caseID, mid, model.Business, origin, reason); err != nil {
		err = errors.Wrap(err, "s.AddAppeal error")
	}
	return
}

// AppealState appeal status .
func (s *Service) AppealState(c context.Context, mid, bid int64) (state bool, err error) {
	block, err := s.BlockedInfoAppeal(c, bid, mid)
	if err != nil {
		err = errors.Wrap(err, "BlockedInfo error")
		return
	}
	if block == nil || block.UID != mid {
		err = ecode.CreditBlockNotExist
		return
	}
	if xtime.Time(time.Now().AddDate(0, 0, -7).Unix()) > block.CTime {
		err = ecode.CreditBlockExpired
		return
	}
	aps, err := s.dao.AppealList(c, mid, model.Business)
	if err != nil {
		err = errors.Wrap(err, "s.dao.AppealList error")
		return
	}
	for _, v := range aps {
		if block.CaseID == v.Oid {
			return
		}
	}
	state = true
	return
}
