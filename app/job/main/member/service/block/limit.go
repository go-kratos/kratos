package block

import (
	"context"
	"time"

	model "go-common/app/job/main/member/model/block"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func (s *Service) limitExpireHandler(c context.Context) {
	if s.conf.BlockProperty.LimitExpireCheckLimit <= 0 {
		log.Error("conf.Conf.Property.LimitExpireCheckLimit [%d] <= 0", s.conf.BlockProperty.LimitExpireCheckLimit)
		return
	}
	var (
		mids    = make([]int64, s.conf.BlockProperty.LimitExpireCheckLimit)
		startID int64
		err     error
	)
	for len(mids) >= s.conf.BlockProperty.LimitExpireCheckLimit {
		log.Info("limit expire handle startID (%d)", startID)
		if startID, mids, err = s.dao.UserStatusList(c, model.BlockStatusLimit, startID, s.conf.BlockProperty.LimitExpireCheckLimit); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, mid := range mids {
			var ok bool
			if ok, err = s.limitExpireCheck(c, mid); err != nil {
				log.Error("%+v", err)
				continue
			}
			if ok {
				if err = s.limitExpireRemove(c, mid); err != nil {
					log.Error("error: %+v, mid: %d", err, mid)
					continue
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (s *Service) limitExpireCheck(c context.Context, mid int64) (ok bool, err error) {
	var (
		his *model.DBHistory
	)
	if his, err = s.dao.UserLastHistory(c, mid); err != nil {
		return
	}
	if his == nil {
		return
	}
	if his.Action != model.BlockActionLimit {
		return
	}
	if his.StartTime.Add(time.Duration(his.Duration) * time.Second).After(time.Now()) {
		return
	}
	ok = true
	return
}

func (s *Service) limitExpireRemove(c context.Context, mid int64) (err error) {
	var (
		_reason, _comment = "系统自动解封", "系统自动解封"
		stime             = time.Now()
		db                = &model.DBHistory{
			MID:       mid,
			AdminID:   model.BlockJOBManagerID,
			AdminName: model.BlockJOBManagerName,
			Source:    model.BlockSourceRemove,
			Area:      model.BlockAreaNone,
			Reason:    _reason,
			Comment:   _comment,
			Action:    model.BlockActionSelfRemove,
			StartTime: stime,
			Duration:  0,
			Notify:    false,
		}
		tx *xsql.Tx
	)
	if tx, err = s.dao.BeginTX(c); err != nil {
		return
	}
	if err = s.dao.TxInsertHistory(c, tx, db); err != nil {
		tx.Rollback()
		return
	}
	count, err := s.dao.TxUpsertUser(c, tx, mid, model.BlockStatusFalse)
	if err != nil || count == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		err = errors.WithStack(err)
	}
	s.mission(func() {
		if err := s.notifyRemoveMSG(context.TODO(), []int64{mid}); err != nil {
			log.Error("%+v", err)
		}
	})
	s.mission(func() {
		s.AddAuditLog(context.TODO(), model.BlockActionSelfRemove, model.BlockJOBManagerID, model.BlockJOBManagerName, []int64{mid}, 0, model.BlockSourceRemove, model.BlockAreaNone, _reason, _comment, false, stime)
	})
	s.cache.Save(func() {
		if err := s.dao.DeleteUserCache(context.TODO(), mid); err != nil {
			log.Error("%+v", err)
		}
		if databusErr := s.dao.AccountNotify(context.TODO(), mid); databusErr != nil {
			log.Error("%+v", databusErr)
		}
	})
	return
}
