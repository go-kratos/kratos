package service

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/model"

	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// JoinBgm join bgm
func (s *Service) JoinBgm(c context.Context, mid int64, accountType, signType int) (err error) {
	id, err := s.dao.Blocked(c, mid)
	if err != nil {
		log.Error("s.dao.GetBlocked mid(%d) error(%v)", mid, err)
		return
	}
	if id != 0 {
		log.Info("mid(%d) is blocked", mid)
		return ecode.GrowupDisabled
	}

	ok, err := s.checkBgmStat(c, mid)
	if err != nil {
		log.Info("s.checkBgmStat error(%v)", err)
		return
	}
	if !ok {
		return ecode.GrowupDisabled
	}

	count, err := s.dao.BGMCount(c, mid)
	if err != nil {
		log.Info("s.dao.BGMCount error(%v)", err)
		return
	}

	avStat, err := s.dao.GetAccountState(c, "up_info_video", mid)
	if err != nil {
		return
	}

	if avStat >= 5 && avStat <= 7 {
		return ecode.GrowupDisabled
	}

	columnStat, err := s.dao.GetAccountState(c, "up_info_column", mid)
	if err != nil {
		return
	}

	if columnStat >= 5 && columnStat <= 7 {
		return ecode.GrowupDisabled
	}

	card, err := s.dao.Card(c, mid)
	if err != nil {
		log.Error("s.dao.Card(%d) error(%v)", mid, err)
		return
	}
	fans, err := s.dao.Fans(c, mid)
	if err != nil {
		return
	}

	now := xtime.Time(time.Now().Unix())
	// sign_type: 1.basic; 2.first publish; 0:default.
	v := &model.UpInfo{
		MID:         mid,
		Nickname:    card.Name,
		AccountType: accountType,
		Fans:        fans,
		SignType:    signType,
		SignedAt:    now,
		Bgms:        count,
	}

	var tx *sql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}

	if _, err = s.dao.TxInsertBgmUpInfo(tx, v); err != nil {
		tx.Rollback()
		return
	}

	if _, err = s.dao.TxInsertCreditScore(tx, mid); err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error("tx.Commit error(%v)", err)
	}
	return
}
