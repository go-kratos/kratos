package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/app/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func (s *Service) statConsumeproc(bus string) {
	defer s.waiter.Done()
	statSub, ok := s.statSub[bus]
	if !ok {
		return
	}
	for {
		msg, ok := <-statSub.Messages()
		if !ok {
			log.Info("stat %s databus Consumer exit", bus)
			break
		}
		msg.Commit()
		ms := &model.StatMsg{}
		if err := json.Unmarshal(msg.Value, ms); err != nil {
			log.Error("stat %s json.Unmarshal(%s) error(%v)", bus, msg.Value, err)
			continue
		}
		if (ms.Type != model.TypeArchive && ms.Type != model.TypeArchiveHis) || ms.ID < 1 {
			log.Warn("stat %s error", msg.Value)
			continue
		}
		ms.BusType = bus
		s.statChan[ms.ID%s.sliceCnt] <- ms
		log.Info("chan(%d) got message(%+v)", ms.ID%s.sliceCnt, ms)
	}
}

func (s *Service) statproc(i int64) {
	defer s.waiter.Done()
	var statChan = s.statChan[i]
	for {
		ms, ok := <-statChan
		if !ok {
			log.Error("statChan[%d] closed", i)
			break
		}
		s.upStatCache(ms.BusType, ms)
	}
}

func (s *Service) upStatCache(bus string, ms *model.StatMsg) {
	var (
		st  *api.Stat
		err error
	)
	c := context.Background()
	defer func() {
		if err != nil {
			log.Error("%+v", err)
			retry := &model.Retry{Action: model.ActionUpStat}
			retry.Data.Aid = ms.ID
			if err = s.vdao.PushFail(c, retry); err != nil {
				log.Error("%+v", err)
			}
			return
		}
		log.Info("update stat cache aid(%d) st(%+v) success", ms.ID, st)
	}()
	if st, err = s.vdao.StatCache(c, ms.ID); err != nil {
		return
	}
	if st == nil {
		arg := &archive.ArgAid2{Aid: ms.ID}
		if st, err = s.arcRPC.Stat3(c, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
			return
		}
	}
	if st == nil {
		return
	}
	switch bus {
	case model.TypeForView:
		st.View = ms.Count
	case model.TypeForDm:
		st.Danmaku = ms.Count
	case model.TypeForReply:
		st.Reply = ms.Count
	case model.TypeForFav:
		st.Fav = ms.Count
	case model.TypeForCoin:
		st.Coin = ms.Count
	case model.TypeForShare:
		st.Share = ms.Count
	case model.TypeForLike:
		st.Like = ms.Count
		st.DisLike = ms.DislikeCount
	case model.TypeForRank:
		st.HisRank = ms.Count
	}
	err = s.vdao.UpStatCache(c, st)
}
