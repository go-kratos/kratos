package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/job/main/feed/dao"
	"go-common/app/job/main/feed/model"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_insertAct = "insert"
	_updateAct = "update"

	_retryTimes = 10
	_retrySleep = time.Second
)

func retry(_retryTimes int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; i < _retryTimes; i++ {
		err = callback()
		if err == nil {
			return
		}
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
			return
		}
		time.Sleep(sleep)
		log.Error("retrying after error: %v", err)
	}
	return fmt.Errorf("after %d _retryTimes, last error: %s", _retryTimes, err)
}

func (s *Service) archiveUpdate(action string, nwMsg []byte, oldMsg []byte) {
	var err error
	nw := &model.Archive{}
	if err = json.Unmarshal(nwMsg, nw); err != nil {
		dao.PromError("archive:解析新稿件")
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	switch action {
	case _insertAct:
		s.insert(nw)
	case _updateAct:
		old := &model.Archive{}
		if err = json.Unmarshal(oldMsg, old); err != nil {
			dao.PromError("archive:解析旧稿件")
			log.Error("json.Unmarshal(%s) error(%v)", string(oldMsg), err)
			return
		}
		s.update(nw, old)
	}
}

func (s *Service) insert(nw *model.Archive) {
	var (
		err   error
		sleep = time.Second
		ts    int64
	)
	if !nw.IsNormal() {
		return
	}
	if ts, err = transTime(nw.PubTime); err != nil {
		dao.PromError("archive:插入时解析时间")
		log.Error("s.insert.transTime(%v) err: %v", nw.PubTime, err)
		err = nil
	}
	if err = retry(_retryTimes, sleep, func() error {
		return s.feedRPC.AddArc(context.TODO(), &feedmdl.ArgArc{Aid: nw.ID, Mid: nw.Mid, PubDate: ts})
	}); err != nil {
		dao.PromError("archive:增加新稿件")
		log.Error("s.feedRPC.AddArc(%v)", nw.ID)
	} else {
		dao.PromInfo("archive:增加新稿件")
		log.Info("feed: Archive(%v) passed, new_arc: %+v.", nw.ID, *nw)
	}
}

func (s *Service) update(nw *model.Archive, old *model.Archive) {
	var (
		err     error
		pubDate int64
	)
	if old.Mid != nw.Mid {
		arg := &feedmdl.ArgChangeUpper{Aid: nw.ID, OldMid: old.Mid, NewMid: nw.Mid}
		if err = retry(_retryTimes, _retrySleep, func() error { return s.feedRPC.ChangeArcUpper(context.TODO(), arg) }); err != nil {
			dao.PromError("archive:修改up主")
			log.Error("s.feedRPC.ChangeArcUpper(%v) error(%v)", arg.Aid, err)
		} else {
			dao.PromInfo("archive:修改up主")
			log.Info("feed: Archive(%v) change upper, old: %+v new: %+v", old.ID, *old, *nw)
		}
	}
	if nw.IsNormal() && ((old.PubTime != nw.PubTime) || (old.State != nw.State)) {
		if pubDate, err = transTime(nw.PubTime); err != nil {
			dao.PromError("archive:解析发布时间")
			log.Error("s.Feed ParsePubTime(%v) error(%v)", nw.PubTime, err)
			err = nil
		}
		if err = retry(_retryTimes, _retrySleep, func() error {
			return s.feedRPC.AddArc(context.TODO(), &feedmdl.ArgArc{Aid: nw.ID, PubDate: pubDate, Mid: nw.Mid})
		}); err != nil {
			dao.PromError("archive:增加新稿件")
			log.Error("s.feedRPC.AddArc error(%v) old_arc: %+v, new_arc: %+v.", err, nw.ID, *old, *nw)
		} else {
			dao.PromInfo("archive:增加新稿件")
			log.Info("feed: Archive(%v) passed, old_arc: %+v, new_arc: %+v.", nw.ID, *old, *nw)
		}
		return
	}
	if (old.State != nw.State) && !nw.IsNormal() {
		if err = retry(_retryTimes, _retrySleep, func() error { return s.feedRPC.DelArc(context.TODO(), &feedmdl.ArgAidMid{Aid: nw.ID, Mid: nw.Mid}) }); err != nil {
			dao.PromError("archive:删除稿件")
			log.Error("s.feedRPC.DelArc(%v) error(%v)", nw.ID, err)
		} else {
			dao.PromInfo("archive:删除稿件")
			log.Info("feed: Archive(%v) not passed, old_arc: %+v, new_arc:%+v.", nw.ID, *old, *nw)
		}
	}
}

func transTime(t string) (res int64, err error) {
	var ts time.Time
	if ts, err = time.Parse("2006-01-02 15:04:05 MST", t+" CST"); err != nil {
		return
	}
	res = ts.Unix()
	return
}
