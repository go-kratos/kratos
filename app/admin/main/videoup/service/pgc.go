package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// PassByPGC update pgc archive state to StateOpen.
func (s *Service) PassByPGC(c context.Context, aid int64, gid int64, attrs map[uint]int32, redirectURL string, now time.Time) (err error) {
	// archive
	var a *archive.Archive
	if a, err = s.arc.Archive(c, aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", aid, err)
		return
	}
	log.Info("aid(%d) begin tran pass pgc", aid)
	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if a.State != archive.StateOpen {
		var firstPass bool
		if firstPass, err = s.txUpArcState(c, tx, a.Aid, archive.StateOpen); err != nil {
			tx.Rollback()
			log.Error("PassByPGC s.txUpArcState(aid(%d),state(%d)) error(%v)", aid, archive.StateOpen, err)
			return
		}
		a.State = archive.StateOpen
		log.Info("archive(%d) update archive state(%d)", a.Aid, a.State)

		// archive ptime
		if firstPass {
			pTime := xtime.Time(now.Unix())
			if _, err = s.arc.TxUpArcPTime(tx, a.Aid, pTime); err != nil {
				tx.Rollback()
				log.Error("s.arc.TxUpArcPTime(%d, %d) error(%v)", a.Aid, pTime, err)
				return
			}
			a.PTime = pTime
			log.Info("archive(%d) second_round upPTime(%d)", a.Aid, a.PTime)
		}
		var round = s.archiveRound(c, a, a.Aid, a.Mid, a.TypeID, a.Round, a.State, false)
		if _, err = s.arc.TxUpArcRound(tx, a.Aid, round); err != nil {
			tx.Rollback()
			log.Error("s.arc.TxUpArcRound(%d, %d) error(%v)", a.Aid, round, err)
			return
		}
		a.Round = round
		log.Info("archive(%d) second_round upRound(%d)", a.Aid, a.Round)
	}
	var conts []string
	if conts, err = s.txUpArcAttrs(tx, a, attrs, redirectURL); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran pass pgc", aid)
	if _, err := s.oversea.UpPolicyRelation(c, aid, gid); err != nil {
		conts = append(conts, fmt.Sprintf("[地区展示]应用策略组ID[%d]", gid))
	}
	if len(conts) > 0 {
		s.arc.AddArcOper(c, a.Aid, 221, a.Attribute, a.TypeID, int16(a.State), a.Round, 1, strings.Join(conts, "，"), "")
	}
	// NOTE: send second_round for sync dede.
	s.busSecondRound(aid, 0, false, false, false, false, false, false, "", nil)
	return
}

// ModifyByPGC update pgc archive attributes.
func (s *Service) ModifyByPGC(c context.Context, aid int64, gid int64, attrs map[uint]int32, redirectURL string) (err error) {
	// archive
	var a *archive.Archive
	if a, err = s.arc.Archive(c, aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", aid, err)
		return
	}
	log.Info("aid(%d) begin tran modify pgc", aid)
	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	var conts []string
	if conts, err = s.txUpArcAttrs(tx, a, attrs, redirectURL); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran modify pgc", aid)
	if _, err := s.oversea.UpPolicyRelation(c, aid, gid); err != nil {
		conts = append(conts, fmt.Sprintf("[地区展示]应用策略组ID[%d]", gid))
	}
	if len(conts) > 0 {
		s.arc.AddArcOper(c, a.Aid, 221, a.Attribute, a.TypeID, int16(a.State), a.Round, 1, strings.Join(conts, "，"), "")
	}
	// NOTE: send second_round for sync dede.
	s.busSecondRound(aid, 0, false, false, false, false, false, false, "", nil)
	return
}

// LockByPGC  update pgc archive state to StateForbidLock.
func (s *Service) LockByPGC(c context.Context, aid int64) (err error) {
	// archive
	var a *archive.Archive
	if a, err = s.arc.Archive(c, aid); err != nil || a == nil {
		log.Error("s.arc.Archive(%d) error(%v) or a==nil", aid, err)
		return
	}
	if a.State == archive.StateForbidLock {
		return
	}
	log.Info("aid(%d) begin tran lock pgc", aid)
	// begin tran
	var tx *sql.Tx
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran() error(%v)", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("wocao jingran recover le error(%v)", r)
		}
	}()
	if _, err = s.txUpArcState(c, tx, a.Aid, archive.StateForbidLock); err != nil {
		tx.Rollback()
		log.Error("s.txUpArcState(aid(%d),state(%d)) error(%v)", aid, archive.StateForbidLock, err)
		return
	}
	a.State = archive.StateForbidLock
	log.Info("archive(%d) update archive state(%d)", a.Aid, a.State)
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	log.Info("aid(%d) end tran lock pgc", aid)
	s.arc.AddArcOper(c, a.Aid, 221, a.Attribute, a.TypeID, int16(a.State), a.Round, 1, "", "")
	// NOTE: send second_round for sync dede.
	s.busSecondRound(aid, 0, false, false, false, false, false, false, "", nil)
	return
}
