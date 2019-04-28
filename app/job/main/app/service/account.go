package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/app/model"
	accapi "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

const (
	_actForUname = "updateUname"
	_actForFace  = "updateFace"
	_actForAdmin = "updateByAdmin"
)

func (s *Service) accConsumeproc() {
	var (
		msg *databus.Message
		ok  bool
		err error
	)
	msgs := s.accountNotifySub.Messages()
	for {
		if msg, ok = <-msgs; !ok {
			log.Info("acc databus Consumer exit")
			break
		}
		var ms = &model.AccMsg{}
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			msg.Commit()
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		switch ms.Action {
		case _actForFace, _actForUname, _actForAdmin:
			s.notifyMidMap.Lock()
			s.notifyMidMap.Map[ms.Mid] = ms.Action
			s.notifyMidMap.Unlock()
		}
		msg.Commit()
	}
	s.waiter.Done()
}

func (s *Service) notifyConsumeproc() {
	for {
		time.Sleep(5 * time.Second)
		if s.closed && len(s.notifyMidMap.Map) == 0 {
			break
		}
		s.notifyMidMap.Lock()
		midMap := s.notifyMidMap.Map
		s.notifyMidMap.Map = map[int64]string{}
		s.notifyMidMap.Unlock()
		log.Info("notifyConsumeproc mid map len(%d)", len(midMap))
		for mid, action := range midMap {
			if err := s.upNotifyArc(mid, action); err != nil {
				log.Error("%+v", err)
			}
		}
	}
	s.waiter.Done()
}

func (s *Service) upNotifyArc(mid int64, action string) (err error) {
	var (
		cnt          int
		accInfoReply *accapi.InfoReply
		res          map[int64][]*archive.AidPubTime
	)
	c := context.Background()
	defer func() {
		if err != nil {
			log.Error("%+v", err)
			// TODO 等待error确定后开启
			// retry := &model.Retry{Action: model.ActionUpAccount}
			// retry.Data.Mid = mid
			// retry.Data.Action = action
			// s.vdao.PushFail(c, retry)
			return
		}
		log.Info("update notify archive cache mid(%d) action(%s) success", mid, action)
	}()
	arg := &archive.ArgUpCount2{Mid: mid}
	if cnt, err = s.arcRPC.UpCount2(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if cnt < 1 {
		return
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		arg := &archive.ArgUpsArcs2{Mids: []int64{mid}, Pn: 1, Ps: cnt}
		if res, err = s.arcRPC.UpsPassed2(ctx, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
		}
		return
	})
	g.Go(func() (err error) {
		arg := &accapi.MidReq{Mid: mid}
		if accInfoReply, err = s.accAPI.Info3(ctx, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
			return
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	if accInfoReply.Info == nil {
		return
	}
	for _, a := range res[mid] {
		if a == nil {
			continue
		}
		var (
			v                         *archive.View3
			av                        *archive.Archive3
			noChangeView, noChangeArc bool
		)
		g, ctx := errgroup.WithContext(c)
		g.Go(func() (err error) {
			if v, err = s.vdao.ViewCache(ctx, a.Aid); err != nil {
				return
			}
			if v == nil {
				return
			}
			if action == _actForAdmin && v.Author.Name == accInfoReply.Info.Name && v.Author.Face == accInfoReply.Info.Face {
				noChangeView = true
				return
			}
			v.Author.Name = accInfoReply.Info.Name
			v.Author.Face = accInfoReply.Info.Face
			if err = s.vdao.UpViewCache(ctx, v); err != nil {
				err = errors.Wrapf(err, "%v", v)
				return
			}
			log.Info("account notify consumer view mid(%v) aid(%d) view(%v) name(%s) face(%s)", mid, a.Aid, v, accInfoReply.Info.Name, accInfoReply.Info.Face)
			return
		})
		g.Go(func() (err error) {
			if av, err = s.vdao.ArcCache(ctx, a.Aid); err != nil {
				return
			}
			if av == nil {
				return
			}
			if action == _actForAdmin && av.Author.Name == accInfoReply.Info.Name && av.Author.Face == accInfoReply.Info.Face {
				noChangeArc = true
				return
			}
			av.Author.Name = accInfoReply.Info.Name
			av.Author.Face = accInfoReply.Info.Face
			if s.vdao.UpArcCache(ctx, av); err != nil {
				err = errors.Wrapf(err, "%v", av)
				return
			}
			log.Info("account notify consumer archive mid(%v) aid(%d) archive(%v) name(%s) face(%s)", mid, a.Aid, av, accInfoReply.Info.Name, accInfoReply.Info.Face)
			return
		})
		if err = g.Wait(); err != nil {
			return
		}
		if noChangeView && noChangeArc {
			break
		}
	}
	return
}
