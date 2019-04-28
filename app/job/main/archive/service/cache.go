package service

import (
	"context"
	"encoding/json"
	"time"

	jobmdl "go-common/app/job/main/archive/model/databus"
	"go-common/app/job/main/archive/model/result"
	accgrpc "go-common/app/service/main/account/api"
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

func (s *Service) cachesubproc() {
	defer s.waiter.Done()
	var msgs = s.cacheSub.Messages()
	for {
		var (
			msg *databus.Message
			ok  bool
			err error
		)
		if msg, ok = <-msgs; !ok {
			log.Error("s.cachesub.messages closed")
			return
		}
		if s.closeSub {
			return
		}
		m := &jobmdl.Rebuild{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
			continue
		}
		log.Info("cacheSub key(%s) value(%s) start", msg.Key, msg.Value)
		var retryError error
		for {
			var (
				a         *result.Archive
				infoReply *accgrpc.InfoReply
				c         = context.TODO()
			)
			if retryError != nil {
				time.Sleep(10 * time.Millisecond)
			}
			a, retryError = s.resultDao.Archive(c, m.Aid)
			if retryError != nil {
				log.Error("s.resultDao.Archive(%d) error(%v)", m.Aid, retryError)
				continue
			}
			if a == nil || a.Mid == 0 {
				log.Info("cache break archive(%d) not exist or mid==0", m.Aid)
				break
			}
			infoReply, retryError = s.accGRPC.Info3(c, &accgrpc.MidReq{Mid: a.Mid})
			if retryError != nil {
				if ecode.Cause(retryError).Equal(ecode.MemberNotExist) {
					log.Info("archive(%d) mid(%d) not exist", m.Aid, a.Mid)
					break
				}
				log.Error("s.acc.RPC.Info3(%d) error(%v)", m.Aid, retryError)
				continue
			}
			if infoReply == nil {
				log.Error("infoReply mid(%d) err is nil,but info is nil too", a.Mid)
				break
			}
			if infoReply.Info.Name == "" || infoReply.Info.Face == "" {
				log.Error("empty info mid(%d) info(%+v)", infoReply.Info.Mid, infoReply.Info)
				break
			}
			for k, arcRPC := range s.arcServices {
				if retryError = arcRPC.ArcCache2(c, &arcmdl.ArgCache2{Aid: m.Aid, Tp: arcmdl.CacheUpdate}); retryError != nil {
					log.Error("s.arcRPC(%d).ArcCache2(%d) error(%v)", k, m.Aid, retryError)
					continue
				}
			}
			log.Info("archive(%d) mid(%d) uname(%s) update success", m.Aid, infoReply.Info.Mid, infoReply.Info.Name)
			break
		}
		msg.Commit()
	}
}

func (s *Service) accountNotifyproc() {
	defer s.waiter.Done()
	var msgs = s.accountNotifySub.Messages()
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
		if s.closeSub {
			return
		}
		msg.Commit()
		m := &jobmdl.AccountNotify{}
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
		if count, err = s.arcServices[0].UpCount2(c, &arcmdl.ArgUpCount2{Mid: m.Mid}); err != nil {
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
			if am, err = s.arcServices[0].UpArcs3(c, &arcmdl.ArgUpArcs2{Mid: m.Mid, Ps: 2, Pn: 1}); err != nil {
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
			var reply *accgrpc.InfoReply
			if reply, err = s.accGRPC.Info3(c, &accgrpc.MidReq{Mid: m.Mid}); err != nil || reply == nil {
				log.Error("accountNotify accRPC.info3(%d) error(%v)", m.Mid, err)
				continue
			}
			if reply.Info.Name == am[0].Author.Name && reply.Info.Face == am[0].Author.Face {
				log.Info("accountNotify face(%s) name(%s) not change", reply.Info.Face, reply.Info.Name)
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
		for mid := range mids {
			s.updateUpperCache(context.TODO(), mid)
		}
		if s.closeSub && len(s.notifyMid) == 0 {
			return
		}
	}
}

func (s *Service) updateUpperCache(c context.Context, mid int64) (err error) {
	// update archive cache
	var aids []int64
	if aids, err = s.resultDao.UpPassed(c, mid); err != nil {
		log.Error("s.resultDao.UpPassed(%d) error(%v)", mid, err)
		return
	}
	failedCnt := 0
	for _, aid := range aids {
		for k, rpc := range s.arcServices {
			if err = rpc.ArcCache2(c, &arcmdl.ArgCache2{Aid: aid}); err != nil {
				log.Error("s.arcRPC(%d).ArcCache2(%d) mid(%d) error(%v)", k, aid, mid, err)
				failedCnt++
			}
		}
	}
	if failedCnt > 0 {
		log.Error("accountNotify updateUpperCache mid(%d) failed(%d)", mid, failedCnt)
		return
	}
	log.Info("accountNofity updateUpperCache mid(%d) successed(%d)", mid, len(aids))
	return
}
