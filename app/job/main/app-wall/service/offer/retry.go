package offer

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/app-wall/model/offer"
	"go-common/library/log"
)

const (
	_upActiveRetry = 5
	_sleep         = 100 * time.Millisecond
)

func (s *Service) retryproc() {
	defer s.waiter.Done()
	var (
		bs  []byte
		err error
	)
	c := context.TODO()
	msg := &offer.Retry{}
	for {
		if s.closed {
			break
		}
		if bs, err = s.dao.PopFail(c); err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if len(bs) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if err = json.Unmarshal(bs, msg); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", bs, err)
			continue
		}
		log.Info("retry action(%s) data(%s)", msg.Action, bs)
		switch msg.Action {
		case offer.ActionActive:
			if msg.Data != nil {
				if err = retry(func() (err error) {
					return s.dao.Active(c, msg.Data.OS, msg.Data.IMEI, msg.Data.Androidid, msg.Data.Mac, "")
				}, _upActiveRetry, _sleep); err != nil {
					log.Error("%+v", err)
					if err = s.syncRetry(c, offer.ActionActive, msg.Data.OS, msg.Data.IMEI, msg.Data.Androidid, msg.Data.Mac); err != nil {
						log.Error("%+v", err)
					}
					return
				}
			}
		}
	}
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

func (s *Service) syncRetry(c context.Context, action, os, imei, androidid, mac string) (err error) {
	retry := &offer.Retry{Action: action, Data: &offer.Data{OS: os, IMEI: imei, Androidid: androidid, Mac: mac}}
	return s.dao.PushFail(c, retry)
}
