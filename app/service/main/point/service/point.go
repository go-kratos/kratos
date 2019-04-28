package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"go-common/app/service/main/point/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

//PointInfo .
func (s *Service) PointInfo(c context.Context, mid int64) (pi *model.PointInfo, err error) {
	var (
		cache = true
	)
	if pi, err = s.dao.PointInfoCache(c, mid); err != nil {
		log.Error("%+v", err)
		cache = false
	}
	if pi != nil {
		return
	}
	if pi, err = s.dao.PointInfo(c, mid); err != nil {
		log.Error("%+v", err)
		return
	}
	if pi == nil {
		pi = new(model.PointInfo)
		pi.Mid = mid
	}
	if cache {
		s.dao.SetPointInfoCache(c, pi)
	}
	return
}

//PointHistory .
func (s *Service) PointHistory(c context.Context, mid int64, cursor, ps int) (phs []*model.PointHistory, total int, ncursor int, err error) {
	if total, err = s.dao.PointHistoryCount(c, mid); err != nil {
		log.Error("%+v", err)
		return
	}
	if total <= 0 {
		return
	}
	if cursor <= 0 {
		cursor = _defcursor
	}
	if ps <= 0 {
		ps = _defps
	}
	if phs, err = s.dao.PointHistory(c, mid, cursor, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	if len(phs) != 0 {
		ncursor = int(phs[len(phs)-1].ID)
	}
	return
}

//OldPointHistory old point history .
func (s *Service) OldPointHistory(c context.Context, mid int64, pn, ps int) (phs []*model.OldPointHistory, total int, err error) {
	if total, err = s.dao.PointHistoryCount(c, mid); err != nil {
		log.Error("%+v", err)
		return
	}
	if total <= 0 {
		return
	}
	if pn <= 0 {
		pn = _defpn
	}
	if ps <= 0 {
		ps = _defps
	}
	if phs, err = s.dao.OldPointHistory(c, mid, (pn-1)*ps, ps); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

//PointAddByBp by bp.
func (s *Service) PointAddByBp(c context.Context, pa *model.ArgPointAdd) (p int64, err error) {
	var (
		hid         int
		activePoint int64
		ph          = new(model.PointHistory)
		rate        int64
		ok          bool
	)
	if hid, err = s.dao.ExistPointOrder(c, pa.OrderID); err != nil {
		log.Error("point add %+v", err)
		return
	}
	if hid > 0 {
		log.Error("point add repeated consumption %+v", pa)
		return
	}
	if rate, ok = s.c.Property.PointGetRule[strconv.Itoa(pa.ChangeType)]; !ok || rate == 0 {
		log.Info("point add not found rete %+v", pa.ChangeType)
		return
	}
	p = int64(math.Ceil(pa.Bcoin * float64(rate)))
	if p == 0 {
		return
	}
	ph.ChangeType = int(pa.ChangeType)
	ph.OrderID = pa.OrderID
	ph.Point = p
	ph.Remark = pa.Remark
	ph.Mid = pa.Mid
	ph.RelationID = pa.RelationID
	ph.ChangeTime = xtime.Time(time.Now().Unix())
	if _, activePoint, err = s.updatePointWithHistory(c, ph); err != nil {
		log.Error("add point mid(%d) ph(%v) %+v", ph.Mid, ph, err)
		return
	}
	p += activePoint
	s.dao.DelPointInfoCache(c, pa.Mid)
	return
}

func (s *Service) updatePointWithHistory(c context.Context, ph *model.PointHistory) (pointBalance int64, activePoint int64, err error) {
	var (
		tx *xsql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err == nil {
			if err = tx.Commit(); err != nil {
				err = errors.WithStack(err)
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
	}()
	if pointBalance, err = s.updatePoint(c, tx, ph.Mid, ph.Point); err != nil {
		err = errors.WithStack(err)
		return
	}
	ph.PointBalance = pointBalance
	if _, err = s.dao.InsertPointHistory(c, tx, ph); err != nil {
		err = errors.WithStack(err)
		return
	}
	if ph.ChangeType == model.Contract && ph.Point >= 0 {
		if activePoint, err = s.activeSendPoint(c, tx, ph); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

func (s *Service) updatePoint(c context.Context, tx *xsql.Tx, mid, point int64) (pb int64, err error) {
	var (
		pi  *model.PointInfo
		ver int64
		a   int64
	)
	if pi, err = s.dao.TxPointInfo(c, tx, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if point == 0 {
		return
	}
	if pi == nil {
		if point < 0 {
			err = fmt.Errorf("point not enough")
			return
		}
		pb = point
		pi = new(model.PointInfo)
		pi.Ver = 1
		pi.PointBalance = point
		pi.Mid = mid
		if a, err = s.dao.InsertPoint(c, tx, pi); err != nil {
			err = errors.WithStack(err)
			return
		}
		if a != 1 {
			err = fmt.Errorf("operation failed")
			return
		}

	} else {
		pb = pi.PointBalance + point
		if pb < 0 {
			err = fmt.Errorf("point not enough")
			return
		}
		pi.PointBalance = pb
		ver = pi.Ver
		pi.Ver++
		if a, err = s.dao.UpdatePointInfo(c, tx, pi, ver); err != nil {
			err = errors.WithStack(err)
			return
		}
		if a != 1 {
			err = fmt.Errorf("operation failed")
			return
		}
	}
	return
}

//ConsumePoint .
func (s *Service) ConsumePoint(c context.Context, pc *model.ArgPointConsume) (status int8, err error) {
	var (
		ph *model.PointHistory
	)
	pc.Point = ^pc.Point + 1

	ph = new(model.PointHistory)
	ph.ChangeType = int(pc.ChangeType)
	ph.ChangeTime = xtime.Time(time.Now().Unix())
	ph.Mid = pc.Mid
	ph.RelationID = pc.RelationID
	ph.Remark = pc.Remark
	ph.Point = pc.Point
	if _, _, err = s.updatePointWithHistory(c, ph); err != nil {
		log.Error("consume point mid(%d) ph(%v) %+v", ph.Mid, ph, err)
		return
	}
	s.dao.DelPointInfoCache(c, pc.Mid)
	status = model.SUCCESS
	return
}

// AddPoint .
func (s *Service) AddPoint(c context.Context, pc *model.ArgPoint) (status int8, err error) {
	ph := new(model.PointHistory)
	ph.ChangeType = int(pc.ChangeType)
	ph.ChangeTime = xtime.Time(time.Now().Unix())
	ph.Mid = pc.Mid
	ph.Remark = pc.Remark
	ph.Point = pc.Point
	ph.Operator = pc.Operator
	if _, _, err = s.updatePointWithHistory(c, ph); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.dao.DelPointInfoCache(c, pc.Mid)
	status = model.SUCCESS
	return
}
