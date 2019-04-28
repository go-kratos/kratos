package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/app/model"
	"go-common/app/job/main/app/model/space"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

func (s *Service) retryproc() {
	var (
		bs  []byte
		err error
	)
	c := context.Background()
	retry := &model.Retry{}
	for {
		if s.closed {
			break
		}
		if bs, err = s.vdao.PopFail(c); err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if len(bs) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if err = json.Unmarshal(bs, retry); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", bs, err)
			continue
		}
		log.Info("retry action(%s) data(%s)", retry.Action, bs)
		switch retry.Action {
		case model.ActionUpStat:
			if retry.Data.Aid != 0 {
				s.retryStatCache(retry.Data.Aid)
			}
		case model.ActionUpView:
			if retry.Data.Aid != 0 {
				s.upViewCache([]int64{retry.Data.Aid})
			}
		case model.ActionUpContribute:
			if retry.Data.Mid != 0 {
				s.contributeUpdate(retry.Data.Mid, retry.Data.Attrs, retry.Data.Items)
			}
		case model.ActionUpContributeAid:
			if retry.Data.Mid != 0 {
				s.contributeCache(retry.Data.Mid, retry.Data.Attrs, retry.Data.Time, retry.Data.IP)
			}
		case model.ActionUpViewContribute:
			if retry.Data.Mid != 0 {
				if retry.Data.Mid != 0 {
					s.upViewContribute(retry.Data.Mid)
				}
			}
		case model.ActionUpAccount:
			if retry.Data.Mid != 0 {
				s.upNotifyArc(retry.Data.Mid, retry.Data.Action)
			}
		}
	}
	s.waiter.Done()
}

func retry(callback func() error, retry int, sleep time.Duration) (err error) {
	for i := 0; i < retry; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(sleep)
	}
	return
}

func (s *Service) syncRetry(c context.Context, action string, mid, aid int64, attrs *space.Attrs, items []*space.Item, time xtime.Time, ip string) (err error) {
	retry := &model.Retry{Action: action}
	retry.Data.Mid = mid
	retry.Data.Aid = aid
	retry.Data.Attrs = attrs
	retry.Data.Items = items
	retry.Data.Time = time
	retry.Data.IP = ip
	return s.vdao.PushFail(c, retry)
}

func (s *Service) retryStatCache(aid int64) {
	var (
		st  *api.Stat
		err error
	)
	c := context.Background()
	defer func() {
		if err != nil {
			log.Error("%+v", err)
			retry := &model.Retry{Action: model.ActionUpStat}
			retry.Data.Aid = aid
			if err = s.vdao.PushFail(c, retry); err != nil {
				log.Error("%+v", err)
			}
			return
		}
		log.Info("retry update stat cache aid(%d) st(%+v) success", aid, st)
	}()
	arg := &archive.ArgAid2{Aid: aid}
	if st, err = s.arcRPC.Stat3(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if st == nil {
		return
	}
	err = s.vdao.UpStatCache(c, st)
}
