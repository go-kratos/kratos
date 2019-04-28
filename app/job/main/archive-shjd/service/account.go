package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/archive-shjd/model"
	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_actForUname = "updateUname"
	_actForFace  = "updateFace"
	_actForAdmin = "updateByAdmin"
)

func (s *Service) accountNotifyproc() {
	defer s.waiter.Done()
	var msgs = s.accountNotify.Messages()
	for {
		var (
			msg *databus.Message
			ok  bool
			err error
			c   = context.TODO()
		)
		if msg, ok = <-msgs; !ok {
			log.Error("s.cachesub.messages closed")
			return
		}
		msg.Commit()
		m := &model.AccountNotify{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		log.Info("accountNotify got key(%s) value(%s)", msg.Key, msg.Value)
		if m.Action != _actForAdmin && m.Action != _actForFace && m.Action != _actForUname {
			log.Warn("accountNotify skip action(%s) values(%s)", m.Action, msg.Value)
			continue
		}
		var count int
		if count, err = s.arcRPCs["group1"].UpCount2(c, &arcmdl.ArgUpCount2{Mid: m.Mid}); err != nil {
			log.Error("s.arcRPC.UpCount2(%d) error(%v)", m.Mid, err)
			continue
		}
		if count == 0 {
			log.Info("accountNotify mid(%d) passed(%d)", m.Mid, count)
			continue
		}
		if m.Action == _actForAdmin {
			// check uname or face is updated
			var am []*api.Arc
			if am, err = s.arcRPCs["group1"].UpArcs3(c, &arcmdl.ArgUpArcs2{Mid: m.Mid, Ps: 2, Pn: 1}); err != nil {
				if ecode.Cause(err).Equal(ecode.NothingFound) {
					err = nil
					log.Info("accountNotify mid(%d) no passed archive", m.Mid)
					continue
				}
				log.Error("accountNotify mid(%d) error(%v)", m.Mid, err)
				continue
			}
			if len(am) == 0 {
				log.Info("accountNotify mid(%d) no passed archive", m.Mid)
				continue
			}
			var info *accmdl.Info
			if info, err = s.accRPC.Info3(c, &accmdl.ArgMid{Mid: m.Mid}); err != nil {
				log.Error("accountNotify accRPC.info3(%d) error(%v)", m.Mid, err)
				continue
			}
			if info.Name == am[0].Author.Name && info.Face == am[0].Author.Face {
				log.Info("accountNotify face(%s) name(%s) not change", info.Face, info.Name)
				continue
			}
		}
		s.notifyMu.Lock()
		s.notifyMid[m.Mid] = struct{}{}
		s.notifyMu.Unlock()
	}
}

func (s *Service) clearMidCache() {
	defer s.waiter.Done()
	for {
		time.Sleep(5 * time.Second)
		s.notifyMu.Lock()
		mids := s.notifyMid
		s.notifyMid = make(map[int64]struct{})
		s.notifyMu.Unlock()
		log.Info("start clearMidCache mids(%d)", len(mids))
		for mid := range mids {
			s.updateUpperCache(context.TODO(), mid)
		}
		log.Info("finish clearMidCache mids(%d)", len(mids))
		if s.close && len(s.notifyMid) == 0 {
			return
		}
	}
}

func (s *Service) updateUpperCache(c context.Context, mid int64) (err error) {
	failedCnt := 0
	for k, rpc := range s.arcRPCs {
		pn := 1
		for {
			var arcs []*api.Arc
			if arcs, err = rpc.UpArcs3(c, &arcmdl.ArgUpArcs2{Mid: mid, Pn: pn}); err != nil {
				log.Error("rpc(%s) UpArcs3(%d) error(%v)", k, mid, err)
				break
			}
			pn++
			if len(arcs) == 0 {
				break
			}
			for _, arc := range arcs {
				if err = rpc.ArcCache2(c, &arcmdl.ArgCache2{Aid: arc.Aid, Tp: arcmdl.CacheUpdate}); err != nil {
					log.Error("s.arcRPC(%d).ArcCache2(%d, %s) mid(%d) error(%v)", k, arc.Aid, arcmdl.CacheUpdate, mid, err)
					failedCnt++
					continue
				}
			}

		}
	}
	if failedCnt > 0 {
		log.Error("accountNotify updateUpperCache mid(%d) failed(%d)", mid, failedCnt)
		return
	}
	log.Info("accountNofity updateUpperCache mid(%d)", mid)
	return
}
