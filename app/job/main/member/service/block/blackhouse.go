package block

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "go-common/app/job/main/member/model/block"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

func (s *Service) creditExpireHandler(c context.Context) {
	if s.conf.BlockProperty.CreditExpireCheckLimit <= 0 {
		log.Error("conf.Conf.Property.creditExpireCheckLimit [%d] <= 0", s.conf.BlockProperty.CreditExpireCheckLimit)
		return
	}
	var (
		mids    = make([]int64, s.conf.BlockProperty.CreditExpireCheckLimit)
		startID int64
		err     error
	)
	for len(mids) >= s.conf.BlockProperty.CreditExpireCheckLimit {
		log.Info("black house expire handle startID (%d)", startID)
		if startID, mids, err = s.dao.UserStatusList(c, model.BlockStatusCredit, startID, s.conf.BlockProperty.CreditExpireCheckLimit); err != nil {
			log.Error("%+v", err)
			return
		}
		for _, mid := range mids {
			log.Info("Start handle black house mid (%d)", mid)
			var ok bool
			if ok, err = s.creditExpireCheck(c, mid); err != nil {
				log.Error("%+v", err)
				continue
			}
			if ok {
				log.Info("Start remove black house mid (%d)", mid)
				if err = s.creditExpireRemove(c, mid); err != nil {
					log.Error("error: %+v, mid: %d", err, mid)
				}
			}
		}
	}
}

func (s *Service) creditExpireCheck(c context.Context, mid int64) (ok bool, err error) {
	var (
		his *model.DBHistory
		ex  *model.DBExtra
	)
	if his, err = s.dao.UserLastHistory(c, mid); err != nil {
		return
	}
	if his == nil {
		return
	}
	log.Info("Credit check his (%+v)", his)
	if his.Action != model.BlockActionLimit {
		return
	}
	if ex, err = s.dao.UserExtra(c, mid); err != nil {
		return
	}
	if ex == nil {
		return
	}
	log.Info("Credit check extra (%+v)", his)
	if ex.ActionTime.Before(his.StartTime) {
		return
	}
	if his.StartTime.Add(time.Duration(his.Duration) * time.Second).After(time.Now()) {
		return
	}
	ok = true
	return
}

func (s *Service) creditExpireRemove(c context.Context, mid int64) (err error) {
	var (
		db = &model.DBHistory{
			MID:       mid,
			AdminID:   -1,
			AdminName: "sys",
			Source:    model.BlockSourceRemove,
			Area:      model.BlockAreaNone,
			Reason:    "小黑屋自动解封",
			Comment:   "小黑屋自动解封",
			Action:    model.BlockActionSelfRemove,
			StartTime: time.Now(),
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

func (s *Service) notifyRemoveMSG(c context.Context, mids []int64) (err error) {
	code, title, content := s.MSGRemoveInfo()
	if err = s.dao.SendSysMsg(c, code, mids, title, content, ""); err != nil {
		return
	}
	return
}

// databus
func (s *Service) creditsubproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("%+v", errors.WithStack(fmt.Errorf("s.creditsubproc panic(%v)", x)))
			go s.creditsubproc()
			log.Info("s.creditsubproc recover")
		}
	}()
	var (
		msg      *databus.Message
		eventMSG *model.CreditAnswerMSG
		err      error
		msgChan  = s.creditSub.Messages()
		c        = context.TODO()
	)
	for msg = range msgChan {
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit error(%v)", err)
		}
		eventMSG = &model.CreditAnswerMSG{}
		if err = json.Unmarshal([]byte(msg.Value), eventMSG); err != nil {
			log.Error("%+v", errors.WithStack(err))
			continue
		}
		if err = s.handleCreditAnswerMSG(c, eventMSG); err != nil {
			log.Error("%+v", err)
			continue
		}
		log.Info("s.handleCreditAnswerMSG(%v) msg", eventMSG)
	}

	log.Info("creditsubproc end")
}

func (s *Service) handleCreditAnswerMSG(c context.Context, msg *model.CreditAnswerMSG) (err error) {
	if msg.MID <= 0 {
		return
	}
	var (
		extra = &model.DBExtra{
			MID:              msg.MID,
			CreditAnswerFlag: true,
			ActionTime:       msg.MTime.Time(),
		}
		checkFlag bool
	)
	if err = s.dao.InsertExtra(c, extra); err != nil {
		return
	}
	// 及时检查解封
	log.Info("Start check black house mid (%d) from answer", extra.MID)
	if checkFlag, err = s.creditExpireCheck(c, extra.MID); err != nil {
		return
	}
	if checkFlag {
		log.Info("Start remove black house mid (%d)", extra.MID)
		if err = s.creditExpireRemove(c, extra.MID); err != nil {
			log.Error("error: %+v, mid: %d", err, extra.MID)
			return
		}
	}
	return
}
