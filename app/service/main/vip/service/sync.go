package service

import (
	"context"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// SyncUser .
// FIXME after remove vip java.
func (s *Service) SyncUser(c context.Context, user *model.VipUserInfo) (err error) {
	var (
		vipdb *model.VipInfoDB
		tx    *sql.Tx
		eff   int64
	)
	s.fixRecentTime(user)
	log.Info("SyncUser start %+v, tonew: %+v", user, user.ToNew())
	if vipdb, err = s.dao.VipInfo(c, user.Mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Info("sync user info begin (%+v)", vipdb)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			if x := tx.Rollback(); x != nil {
				err = errors.WithStack(x)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = errors.Wrapf(err, "user(%d) commit failed.", user.Mid)
			return
		}
		err = s.dao.DelVipInfoCache(context.Background(), user.Mid)
	}()
	if vipdb != nil {
		if vipdb.Ver >= user.Ver {
			log.Info("sync user(%d) info ver(new:%d, db:%d) too low ", user.Mid, user.Ver, vipdb.Ver)
			return
		}
		if eff, err = s.dao.SyncUpdateUser(tx, user.ToNew(), vipdb.Ver); eff != 1 {
			log.Info("sync user info update eff (%d), err(%+v)", eff, err)
			err = errors.Errorf("SyncUpdateUser mid:%d, eff:%d", user.Mid, eff)
			return
		}
	} else {
		err = s.dao.SyncAddUser(tx, user.ToNew())
	}
	if err != nil {
		log.Info("sync user info err (%+v)", err)
		return
	}
	if user.AutoRenewed == 1 {
		udh := &model.VipUserDiscountHistory{
			DiscountID: model.VipUserFirstDiscount,
			Status:     model.DiscountUsed,
			Mid:        user.Mid,
		}
		if _, err = s.dao.TxAddUserDiscount(tx, udh); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	log.Info("sync user info success (%d)", user.Mid)
	return
}

func (s *Service) fixRecentTime(user *model.VipUserInfo) {
	if user.VipRecentTime == 0 {
		user.VipRecentTime = user.Mtime
	}
	if user.AnnualVipOverdueTime == 0 {
		if empty, err := time.ParseInLocation("2006-01-02 15:04:05", "2016-01-01 00:00:00", time.Local); err == nil {
			user.AnnualVipOverdueTime = xtime.Time(empty.Unix())
		}
	}
}
