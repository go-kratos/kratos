package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/archive-shjd/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// consumerproc consumer all topic
func (s *Service) consumerproc(k string, d *databus.Databus) {
	defer s.waiter.Done()
	var msgs = d.Messages()
	for {
		var (
			err error
			ok  bool
			msg *databus.Message
			now = time.Now().Unix()
		)
		msg, ok = <-msgs
		if !ok || s.close {
			log.Info("databus(%s) consumer exit", k)
			return
		}
		msg.Commit()
		var ms = &model.StatCount{}
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
			continue
		}
		if ms.Aid <= 0 || (ms.Type != "archive" && ms.Type != "archive_his") {
			log.Warn("message(%s) error", msg.Value)
			continue
		}
		if now-ms.TimeStamp > 1800 {
			log.Warn("topic(%s) message(%s) too early", msg.Topic, msg.Value)
			continue
		}
		stat := &model.StatMsg{Aid: ms.Aid, Type: k, Ts: ms.TimeStamp}
		switch k {
		case model.TypeForView:
			stat.Click = ms.Count
		case model.TypeForDm:
			stat.DM = ms.Count
		case model.TypeForReply:
			stat.Reply = ms.Count
		case model.TypeForFav:
			stat.Fav = ms.Count
		case model.TypeForCoin:
			stat.Coin = ms.Count
		case model.TypeForShare:
			stat.Share = ms.Count
		case model.TypeForRank:
			stat.HisRank = ms.Count
		case model.TypeForLike:
			stat.Like = ms.Count
		default:
			log.Error("unknow type(%s) message(%s)", k, msg.Value)
			continue
		}
		s.subStatCh[stat.Aid%_sharding] <- stat
		log.Info("got message(%+v)", stat)
	}
}

func (s *Service) statDealproc(i int) {
	defer s.waiter.Done()
	var (
		ch  = s.subStatCh[i]
		sm  = s.statSM[i]
		c   = context.TODO()
		ls  *lastTmStat
		err error
	)
	for {
		ms, ok := <-ch
		if !ok {
			log.Info("statDealproc(%d) quit", i)
			return
		}
		// get stat
		if ls, ok = sm[ms.Aid]; !ok {
			var stat *api.Stat
			for _, arc := range s.arcRPCs {
				if stat, err = arc.Stat3(c, &archive.ArgAid2{Aid: ms.Aid}); err == nil {
					break
				}
			}
			if stat == nil {
				log.Error("stat(%d) is nill", ms.Aid)
				continue
			}
			ls = &lastTmStat{}
			ls.stat = stat
			sm[ms.Aid] = ls
		}
		model.Merge(ms, ls.stat)
		// update cache
		st := &api.Stat{
			Aid:     ls.stat.Aid,
			View:    ls.stat.View,
			Danmaku: ls.stat.Danmaku,
			Reply:   ls.stat.Reply,
			Fav:     ls.stat.Fav,
			Coin:    ls.stat.Coin,
			Share:   ls.stat.Share,
			NowRank: ls.stat.NowRank,
			HisRank: ls.stat.HisRank,
			Like:    ls.stat.Like,
		}
		for cluster, arc := range s.arcRPCs {
			if err = arc.SetStat2(c, st); err != nil {
				log.Error("s.arcRPC.SetStat2(%s) (%+v) error(%v)", cluster, st, err)
			}
		}
	}
}
