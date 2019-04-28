package service

import (
	"context"
	"time"

	"go-common/app/interface/main/growup/model"

	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"golang.org/x/sync/errgroup"
)

func (s *Service) getUpFans(c context.Context, mid int64) (fans int64, err error) {
	pfl, err := s.dao.ProfileWithStat(c, mid)
	if err != nil {
		return
	}
	fans = int64(pfl.Follower)
	return
}

// GetUpStatus status of user in growup plan by mid
func (s *Service) GetUpStatus(c context.Context, mid int64, ip string) (status *model.UpStatus, err error) {
	id, err := s.dao.Blocked(c, mid)
	if err != nil {
		log.Error("s.dao.Blocked mid(%d) error(%v)", mid, err)
		return
	}
	status = &model.UpStatus{}
	status.Blocked = id != 0
	if status.Blocked {
		return
	}
	status.Status = make([]*model.BusinessStatus, 3)

	var g errgroup.Group
	g.Go(func() (err error) {
		stat, err := s.dao.AvUpStatus(c, mid)
		if err != nil {
			log.Error("s.dao.AvUpStatus mid(%d) error(%v)", mid, err)
			return
		}
		if stat.AccountState == 0 {
			if stat.QuitAt.After(stat.CTime) {
				stat.AccountState = 8
			}
		}
		// check av threshold
		if stat.AccountState == 3 {
			stat.ShowPanel = true
		} else {
			stat.ShowPanel, err = s.checkAvStat(c, mid, ip)
			if err != nil {
				log.Error("s.checkAvStat mid(%d) error(%v)", mid, err)
				return
			}
		}
		stat.IsWhite = true
		stat.Type = 0
		status.Status[0] = stat
		return
	})

	g.Go(func() (err error) {
		stat, err := s.dao.ColumnUpStatus(c, mid)
		if err != nil {
			log.Error("s.dao.ColumnUpStatus mid(%d) error(%v)", mid, err)
			return
		}
		stat.IsWhite = true
		if stat.AccountState == 3 {
			stat.ShowPanel = true
		} else {
			stat.ShowPanel, err = s.checkArticleStat(c, mid, ip)
			if err != nil {
				log.Error("s.checkArticleStat mid(%d) error(%v)", mid, err)
				return
			}
		}
		stat.Type = 1
		status.Status[1] = stat
		return
	})

	g.Go(func() (err error) {
		stat, err := s.dao.BgmUpStatus(c, mid)
		if err != nil {
			return
		}
		if stat.AccountState == 3 {
			stat.ShowPanel = true
		} else {
			stat.ShowPanel, err = s.checkBgmStat(c, mid)
			if err != nil {
				return
			}
		}
		stat.Type = 2
		status.Status[2] = stat
		return
	})
	err = g.Wait()
	return
}

func (s *Service) checkAvStat(c context.Context, mid int64, ip string) (ok bool, err error) {
	identify, err := s.dao.UpBusinessInfos(c, mid)
	if err != nil {
		log.Error("s.dao.UpBusinessInfos mid(%d) error(%v)", mid, err)
		return
	}
	if identify.Archive != 1 {
		ok = false
		return
	}

	stat, err := s.avStat(c, mid, ip)
	if err != nil {
		log.Error("s.dao.AvStat mid(%d) error(%v)", mid, err)
		return
	}

	if stat.Fans >= s.conf.Threshold.LimitFanCnt || stat.View >= s.conf.Threshold.LimitTotalClick {
		ok = true
	} else {
		ok = false
	}
	return
}

func (s *Service) checkArticleStat(c context.Context, mid int64, ip string) (ok bool, err error) {
	stat, err := s.dao.ArticleStat(c, mid, ip)
	if err != nil {
		log.Error("s.dao.ArticleStat mid(%d) error(%v)", mid, err)
		return
	}
	if stat.View >= s.conf.Threshold.LimitArticleView {
		ok = true
	} else {
		ok = false
	}
	return
}

func (s *Service) checkBgmStat(c context.Context, mid int64) (ok bool, err error) {
	count, err := s.dao.BgmUpCount(c, mid)
	if err != nil {
		return
	}
	if count > 0 {
		ok = true
		return
	}
	count, err = s.dao.BgmWhiteList(c, mid)
	if err != nil {
		return
	}
	if count > 0 {
		ok = true
	}
	return
}

// JoinAv add user to growup plan (video)
func (s *Service) JoinAv(c context.Context, accountType int, mid int64, signType int) (err error) {
	id, err := s.dao.Blocked(c, mid)
	if err != nil {
		log.Error("s.dao.GetBlocked mid(%d) error(%v)", mid, err)
		return
	}
	if id != 0 {
		log.Info("mid(%d) is blocked", mid)
		return ecode.GrowupDisabled
	}

	ip := metadata.String(c, metadata.RemoteIP)
	ok, err := s.checkAvStat(c, mid, ip)
	if err != nil {
		log.Error("s.checkAvStat mid(%d) error(%v)", mid, err)
		return
	}
	if !ok {
		log.Info("mid(%d) video not reach standard", mid)
		return ecode.GrowupDisabled
	}

	nickname, categoryID, err := s.dao.CategoryInfo(c, mid)
	if err != nil {
		return
	}
	fans, err := s.dao.Fans(c, mid)
	if err != nil {
		return
	}
	state, err := s.dao.GetAccountState(c, "up_info_video", mid)
	if err != nil {
		return
	}

	// if account state is 2 3 4 5 6 7 return
	if state >= 2 && state < 8 {
		return
	}

	now := xtime.Time(time.Now().Unix())
	// sign_type: 1.basic; 2.first publish; 0:default.
	v := &model.UpInfo{
		MID:          mid,
		Nickname:     nickname,
		AccountType:  accountType,
		MainCategory: categoryID,
		Fans:         fans,
		AccountState: 2,
		SignType:     signType,
		ApplyAt:      now,
	}

	_, err = s.dao.InsertUpInfo(c, "up_info_video", "total_play_count", v)
	return
}

// Quit user quit growup plan
func (s *Service) Quit(c context.Context, mid int64, reason string) (err error) {
	var (
		tx        *sql.Tx
		now       = time.Now().Unix()
		quitAt    = xtime.Time(now)
		expiredIn = xtime.Time(now + 86400*3)
	)

	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}

	nickname, err := s.dao.Nickname(c, mid)
	if err != nil {
		return
	}

	current, err := s.dao.CreditScore(c, mid)
	if err != nil {
		return
	}

	_, err = s.dao.TxQuit(tx, "up_info_video", mid, quitAt, expiredIn, reason)
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = s.dao.TxQuit(tx, "up_info_column", mid, quitAt, expiredIn, reason)
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = s.dao.TxQuit(tx, "up_info_bgm", mid, quitAt, expiredIn, reason)
	if err != nil {
		tx.Rollback()
		return
	}

	cr := &model.CreditRecord{
		MID:       mid,
		OperateAt: xtime.Time(now),
		Operator:  nickname,
		Reason:    5,
		Deducted:  1,
		Remaining: current - 1,
	}

	_, err = s.dao.TxInsertCreditRecord(tx, cr)
	if err != nil {
		tx.Rollback()
		return
	}
	// quit deduct 1
	_, err = s.dao.TxDeductCreditScore(tx, 1, mid)
	if err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
	}
	return
}
