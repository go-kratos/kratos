package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/sms/dao/chuanglan"
	"go-common/app/job/main/sms/dao/mengwang"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_callbackSize = 100
)

func (s *Service) dispatchCallback(provider int32) {
	switch provider {
	case smsmdl.ProviderChuangLan:
		s.waiter.Add(1)
		go s.callbackChuangLanproc()
	case smsmdl.ProviderMengWang:
		s.waiter.Add(1)
		go s.callbackMengWangproc()
	}
}

func (s *Service) callbackChuangLanproc() {
	defer s.waiter.Done()
	log.Info("callbackChuangLanproc start")
	group := errgroup.Group{}
	cli := chuanglan.NewClient(s.c)
	for {
		if s.closed {
			log.Info("callbackChuangLanproc exit")
			return
		}
		group.Go(func() error {
			callbacks, err := cli.Callback(context.Background(), s.c.Provider.ChuangLanSmsUser, s.c.Provider.ChuangLanSmsPwd, s.c.Provider.ChuangLanSmsCallbackURL, _callbackSize)
			if err != nil {
				time.Sleep(time.Second)
				return nil
			}
			s.sendChuangLanCallbacks(smsmdl.TypeSms, callbacks)
			return nil
		})
		group.Go(func() error {
			callbacks, err := cli.Callback(context.Background(), s.c.Provider.ChuangLanActUser, s.c.Provider.ChuangLanActPwd, s.c.Provider.ChuangLanActCallbackURL, _callbackSize)
			if err != nil {
				time.Sleep(time.Second)
				return nil
			}
			s.sendChuangLanCallbacks(smsmdl.TypeActSms, callbacks)
			return nil
		})
		group.Go(func() error {
			callbacks, err := cli.CallbackInternational(context.Background(), _callbackSize)
			if err != nil {
				time.Sleep(time.Second)
				return nil
			}
			s.sendChuangLanCallbacks(smsmdl.TypeSms, callbacks)
			return nil
		})
		group.Wait()
		time.Sleep(time.Second)
	}
}

func (s *Service) sendChuangLanCallbacks(typ int32, cbs []*chuanglan.Callback) (err error) {
	ts := time.Now().Unix()
	for _, cb := range cbs {
		if cb.NotifyTime != "" {
			if t, e := time.ParseInLocation("060102150405", cb.NotifyTime, time.Local); e != nil {
				log.Warn("sendChuangLanCallbacks(%+v) parse time error(%v)", cb, e)
			} else {
				ts = t.Unix()
			}
		}
		s.sendUserActionLog(&smsmdl.ModelUserActionLog{
			MsgID:    cb.MsgID,
			Mobile:   cb.Mobile,
			Status:   cb.Status,
			Desc:     cb.Desc,
			Provider: smsmdl.ProviderChuangLan,
			Type:     typ,
			Action:   smsmdl.UserActionCallback,
			Ts:       ts,
		})
	}
	return
}

func (s *Service) callbackMengWangproc() {
	defer s.waiter.Done()
	log.Info("callbackMengWangproc start")
	group := errgroup.Group{}
	cli := mengwang.NewClient(s.c)
	for {
		if s.closed {
			log.Info("callbackMengWangproc exit")
			return
		}
		group.Go(func() error {
			callbacks, err := cli.Callback(context.Background(), s.c.Provider.MengWangSmsUser, s.c.Provider.MengWangSmsPwd, s.c.Provider.MengWangSmsCallbackURL, _callbackSize)
			if err != nil {
				time.Sleep(time.Second)
				return nil
			}
			s.sendMengWangCallbacks(smsmdl.TypeSms, callbacks)
			return nil
		})
		group.Go(func() error {
			callbacks, err := cli.Callback(context.Background(), s.c.Provider.MengWangActUser, s.c.Provider.MengWangActPwd, s.c.Provider.MengWangActCallbackURL, _callbackSize)
			if err != nil {
				time.Sleep(time.Second)
				return nil
			}
			s.sendMengWangCallbacks(smsmdl.TypeActSms, callbacks)
			return nil
		})
		group.Go(func() error {
			callbacks, err := cli.Callback(context.Background(), s.c.Provider.MengWangInternationUser, s.c.Provider.MengWangInternationPwd, s.c.Provider.MengWangInternationalCallbackURL, _callbackSize)
			if err != nil {
				time.Sleep(time.Second)
				return nil
			}
			s.sendMengWangCallbacks(smsmdl.TypeSms, callbacks)
			return nil
		})
		group.Wait()
		time.Sleep(time.Second)
	}
}

func (s *Service) sendMengWangCallbacks(typ int32, cbs []*mengwang.Callback) (err error) {
	ts := time.Now().Unix()
	for _, cb := range cbs {
		if cb.ReportTime != "" {
			if t, e := time.ParseInLocation("2006-01-02 15:04:05", cb.ReportTime, time.Local); e != nil {
				log.Warn("sendMengWangCallbacks(%+v) parse time error(%v)", cb, e)
			} else {
				ts = t.Unix()
			}
		}
		s.sendUserActionLog(&smsmdl.ModelUserActionLog{
			MsgID:    strconv.FormatInt(cb.MsgID, 10),
			Mobile:   cb.Mobile,
			Status:   cb.Status,
			Desc:     cb.Desc,
			Provider: smsmdl.ProviderMengWang,
			Type:     typ,
			Action:   smsmdl.UserActionCallback,
			Ts:       ts,
		})
	}
	return
}
