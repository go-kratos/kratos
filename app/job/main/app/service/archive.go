package service

import (
	"context"
	"encoding/json"
	"go-common/library/ecode"

	"go-common/app/job/main/app/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_updateAct         = "update"
	_insertAct         = "insert"
	_tableArchive      = "archive"
	_upContributeRetry = 5
)

// arcConsumeproc consumer archive
func (s *Service) arcConsumeproc() {
	var (
		msg *databus.Message
		ok  bool
		err error
	)
	msgs := s.archiveNotifySub.Messages()
	for {
		if msg, ok = <-msgs; !ok {
			log.Info("arc databus Consumer exit")
			break
		}
		var ms = &model.ArcMsg{}
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			msg.Commit()
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		switch ms.Table {
		case _tableArchive:
			s.archiveUpdate(ms.Action, ms.New)
		}
		msg.Commit()
	}
	s.waiter.Done()
}

func (s *Service) archiveUpdate(action string, nwMsg []byte) {
	nw := &model.Archive{}
	if err := json.Unmarshal(nwMsg, nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	switch action {
	case _updateAct, _insertAct:
		s.upViewCache([]int64{nw.Aid})
		s.upViewContribute(nw.Mid)
	}
}

func (s *Service) rebuildViewCache(c context.Context, v *archive.View3) (err error) {
	if !v.IsNormal() {
		if err = s.vdao.DelArcCache(c, v.Aid); err != nil {
			log.Error("%+v", err)
			return
		}
		if err = s.vdao.DelViewCache(c, v.Aid); err != nil {
			log.Error("%+v", err)
			return
		}
		if err = s.spdao.DelContributeIDCache(c, v.Author.Mid, v.Aid, model.GotoAv); err != nil {
			log.Error("%+v", err)
			return
		}
	}
	if err = s.vdao.UpArcCache(c, v.Archive3); err != nil {
		log.Error("%+v", err)
		return
	}
	if err = s.vdao.UpViewCache(c, v); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

func (s *Service) upViewCache(aids []int64) {
	if len(aids) == 0 {
		return
	}
	var (
		vs  map[int64]*archive.View3
		err error
	)
	c := context.Background()
	if vs, err = s.arcRPC.Views3(c, &archive.ArgAids2{Aids: aids}); err != nil {
		log.Error("s.arcRPC.Views3(%+v) error(%v)", aids, err)
		return
	}
	if len(vs) == 0 {
		log.Warn("archive(%+v) not exist", aids)
		return
	}
	for aid, v := range vs {
		var cacheErr error
		if cacheErr = s.rebuildViewCache(c, v); cacheErr != nil {
			retry := &model.Retry{Action: model.ActionUpView}
			retry.Data.Aid = aid
			s.vdao.PushFail(c, retry)
			continue
		}
		log.Info("update view cahce aid(%d) success", aid)
	}
}

func (s *Service) upViewContribute(mid int64) {
	if mid == 0 {
		return
	}
	const (
		_minCnt = 20
		_count  = 6
	)
	var err error
	c := context.Background()
	defer func() {
		if err != nil {
			retry := &model.Retry{Action: model.ActionUpViewContribute}
			retry.Data.Mid = mid
			s.vdao.PushFail(c, retry)
			return
		}
		log.Info("update view contribute cache mid(%d) success", mid)
	}()
	arg := &archive.ArgUpCount2{Mid: mid}
	cnt, err := s.arcRPC.UpCount2(c, arg)
	if err != nil {
		log.Error("s.arcRPC.UpCount2(%v) error(%v)", arg, err)
		return
	}
	if cnt < _minCnt {
		return
	}
	arg2 := &archive.ArgUpArcs2{Mid: mid, Pn: 1, Ps: _count}
	uas, err := s.arcRPC.UpArcs3(c, arg2)
	if err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			log.Warn("s.arcRPC.UpArcs3(%v) error(%v)", arg2, err)
			err = nil
			return
		}
		log.Error("s.arcRPC.UpArcs3(%v) error(%v)", arg2, err)
		return
	}
	if len(uas) < _count {
		return
	}
	aids := make([]int64, 0, len(uas))
	for _, a := range uas {
		aids = append(aids, a.Aid)
	}
	if err = s.vdao.UpViewContributeCache(c, mid, aids); err != nil {
		log.Error("%+v", err)
	}
}
