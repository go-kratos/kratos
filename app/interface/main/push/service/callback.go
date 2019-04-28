package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/push/dao"
	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/dao/jpush"
	"go-common/app/service/main/push/dao/mi"
	"go-common/app/service/main/push/dao/oppo"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/log"
)

func (s *Service) callbackproc() {
	defer s.waiter.Done()
	var data []*pushmdl.Callback
	for {
		v, ok := <-s.callbackCh
		if !ok {
			log.Info("callbackCh has been closed.")
			if len(data) > 0 {
				s.sendCallback(data)
			}
			return
		}
		data = append(data, v)
		if len(data) >= s.c.Push.CallbackSize {
			s.sendCallback(data)
			data = []*pushmdl.Callback{}
		}
	}
}

func (s *Service) sendCallback(v []*pushmdl.Callback) (err error) {
	for i := 0; i < 3; i++ {
		if err = s.dao.PubCallback(context.TODO(), v); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	return
}

func (s *Service) addCallbackChan(cb *pushmdl.Callback) (err error) {
	if s.closed {
		log.Warn("addCallbackChan, channel is closed")
		return
	}
	select {
	case s.callbackCh <- cb:
	default:
		err = errors.New("callbackCh full")
		log.Error("callbackCh full. data(%+v)", cb)
		dao.PromError("callbackCh full")
	}
	return
}

// CallbackXiaomiRegid xiaomi regid callback.
func (s *Service) CallbackXiaomiRegid(c context.Context, cb *mi.RegidCallback) (err error) {
	// 小米token注册回调，暂时没用
	log.Info("s.CallbackXiaomiRegid(%+v)", cb)
	return
}

// CallbackHuawei huawei callback.
func (s *Service) CallbackHuawei(c context.Context, hcb *huawei.Callback) (err error) {
	for _, v := range hcb.Statuses {
		log.Info("huawei callback task(%s) token(%s)", v.BiTag, v.Token)
		appid, _ := strconv.ParseInt(v.AppID, 10, 64)
		cb := &pushmdl.Callback{
			Task:     v.BiTag,
			APP:      appid,
			Platform: pushmdl.PlatformHuawei,
			Pid:      pushmdl.MobiAndroid,
			Token:    v.Token,
			Extra:    &pushmdl.CallbackExtra{Status: v.Status},
		}
		s.addCallbackChan(cb)
	}
	return
}

// CallbackXiaomi xiaomi callback.
func (s *Service) CallbackXiaomi(c context.Context, m map[string]*mi.Callback) (err error) {
	for _, v := range m {
		log.Info("callback xiaomi task(%s)", v.Jobkey)
		barStatus := mi.CallbackBarStatusEnable
		if v.BarStatus == mi.CallbackBarStatusDisableStr {
			barStatus = mi.CallbackBarStatusDisable
		} else if v.BarStatus == mi.CallbackBarStatusUnknownStr {
			barStatus = mi.CallbackBarStatusUnknown
		}
		sp := strings.Split(v.Targets, ",")
		appid, _ := strconv.ParseInt(v.Param, 10, 64)
		for _, t := range sp {
			if t == "" {
				continue
			}
			cb := &pushmdl.Callback{
				Task:     v.Jobkey,
				APP:      appid,
				Platform: pushmdl.PlatformXiaomi,
				Pid:      pushmdl.MobiAndroid,
				Token:    t,
				Extra:    &pushmdl.CallbackExtra{Status: barStatus},
			}
			log.Info("xiaomi callback task(%s) token(%s)", v.Jobkey, t)
			s.addCallbackChan(cb)
		}
	}
	return
}

// CallbackOppo oppo callback.
func (s *Service) CallbackOppo(c context.Context, task string, cbs []*oppo.Callback) (err error) {
	for _, v := range cbs {
		for _, t := range strings.Split(v.Tokens, ",") {
			log.Info("oppo callback task(%s) token(%s)", task, t)
			cb := &pushmdl.Callback{
				Task:     task,
				Platform: pushmdl.PlatformOppo,
				Pid:      pushmdl.MobiAndroid,
				Token:    t,
			}
			s.addCallbackChan(cb)
		}
	}
	return
}

// CallbackJpush jpush callback batch.
func (s *Service) CallbackJpush(c context.Context, cbs []*jpush.CallbackReply) (err error) {
	for _, cb := range cbs {
		var (
			task  string
			appid int64
		)
		if cb.Params != nil {
			task = cb.Params["task"]
			appid, _ = strconv.ParseInt(cb.Params["appid"], 10, 64)
		}
		log.Info("jpush callback task(%s) token(%s) channel(%d)", task, cb.Token, cb.Channel)
		status := jpush.StatusSwitchOn
		if !cb.Switch {
			status = jpush.StatusSwitchOff
		}
		s.addCallbackChan(&pushmdl.Callback{
			Task:     task,
			APP:      appid,
			Platform: pushmdl.PlatformJpush,
			Pid:      pushmdl.MobiAndroid,
			Token:    cb.Token,
			Extra:    &pushmdl.CallbackExtra{Status: status, Channel: cb.Channel},
		})
	}
	return
}

// CallbackIOS ios arrived callback.
func (s *Service) CallbackIOS(c context.Context, task, token string, pid int) (err error) {
	cb := &pushmdl.Callback{
		Task:  task,
		Pid:   pid,
		Token: token,
	}
	err = s.addCallbackChan(cb)
	return
}

// CallbackClick click callback.
func (s *Service) CallbackClick(c context.Context, cb *pushmdl.Callback) error {
	return s.addCallbackChan(cb)
}
