package service

import (
	"context"
	"encoding/json"

	jobmdl "go-common/app/job/main/archive/model/databus"
	"go-common/app/job/main/archive/model/retry"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) videoupConsumer() {
	defer s.waiter.Done()
	var msgs = s.videoupSub.Messages()
	for {
		var (
			msg *databus.Message
			ok  bool
			err error
		)
		if msg, ok = <-msgs; !ok {
			log.Error("s.videoupSub.messages closed")
			return
		}
		if s.closeSub {
			return
		}
		msg.Commit()
		m := &jobmdl.Videoup{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
			continue
		}
		log.Info("videoupMessage key(%s) value(%s) start", msg.Key, msg.Value)
		if m.Aid <= 0 {
			log.Warn("aid(%d) <= 0 WTF(%s)", m.Aid, msg.Value)
			continue
		}
		switch m.Route {
		case jobmdl.RouteAutoOpen, jobmdl.RouteDelayOpen, jobmdl.RouteDeleteArchive, jobmdl.RouteSecondRound, jobmdl.RouteFirstRoundForbid, jobmdl.RouteForceSync:
			select {
			case s.videoupAids[m.Aid%int64(s.c.ChanSize)] <- m.Aid:
			default:
				rt := &retry.Info{Action: retry.FailResultAdd}
				rt.Data.Aid = m.Aid
				s.PushFail(context.TODO(), rt)
				log.Warn("s.videoupAids is full!!! async databus archive(%d)", m.Aid)
			}
		}
		log.Info("videoupMessage key(%s) value(%s) finish", msg.Key, msg.Value)
	}
}

func (s *Service) delVideoCache(aid int64, cids []int64) (err error) {
	for _, cid := range cids {
		for k, rpc := range s.arcServices {
			if err = rpc.DelVideo2(context.TODO(), &archive.ArgVideo2{Aid: aid, Cid: cid}); err != nil {
				log.Error("s.arcRpc(%d).DelVideo2(%d, %d) error(%v)", k, aid, cid, err)
				if ecode.Cause(err) != ecode.NothingFound {
					rt := &retry.Info{Action: retry.FailDelVideoCache}
					rt.Data.Aid = aid
					rt.Data.Cids = []int64{cid}
					s.PushFail(context.TODO(), rt)
					log.Error("delVideoCache error(%v)", err)
				}
			}
		}
	}
	return
}

func (s *Service) upVideoCache(aid int64, cids []int64) (err error) {
	for _, cid := range cids {
		for k, rpc := range s.arcServices {
			if err = rpc.UpVideo2(context.TODO(), &archive.ArgVideo2{Aid: aid, Cid: cid}); err != nil {
				log.Error("s.arcRpc(%d).UpVideo2(%d, %d) error(%v)", k, aid, cid, err)
				if ecode.Cause(err) != ecode.NothingFound {
					rt := &retry.Info{Action: retry.FailUpVideoCache}
					rt.Data.Aid = aid
					rt.Data.Cids = []int64{cid}
					s.PushFail(context.TODO(), rt)
					log.Error("upVideoCache error(%v)", err)
				}
			}
		}
	}
	return
}
