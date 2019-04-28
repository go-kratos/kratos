package service

import (
	"time"

	"go-common/library/log"
)

const (
	_platIOS     = "ios"
	_platAndroid = "android"
	_platAll     = ""
	_vIPad       = "ipad"
	_vIPhone     = "iphone"
	_vAndroid    = "android"
	_vAndroidB   = "android_b"
)

func (s *Service) push(resIDs map[string]int64) {
	var (
		now     = time.Now().Unix()
		timeout = s.c.Cfg.Push.Timeout
		err     error
		msg     string
	)
	for resID, timeV := range resIDs {
		finish := false
		needPush := false
		// distinguish whether the resource is ready to push. calc finish or timeout
		if finish, err = s.pushDao.DiffFinish(ctx, resID); err != nil { // check whether diff cal finish
			continue
		}
		if now-timeV > timeout { // check whether it's already timeout
			needPush = true
			log.Info("CallPush [%v] Because of Timeout", resID)
		} else if finish {
			needPush = true
			log.Info("CallPush [%v] Because of DiffFinish", resID)
		} else {
			log.Info("CallPush Jump [%v]", resID)
			continue
		}
		// prepare api call
		if msg, err = s.pushDao.PushMsg(ctx, resID); err != nil { // prepare msg
			continue
		}
		if needPush {
			if err = s.pushDao.CallRefresh(ctx); err != nil {
				log.Error("CallPush [%d] app-resource refresh error [%v]", resID, err)
				continue
			}
			time.Sleep(time.Duration(s.c.Cfg.Push.Pause))
			if err = s.pushDao.CallPush(ctx, s.platform(resID), msg, ""); err != nil {
				log.Error("CallPush [%v] Error [%v]", resID, err)
				continue
			}
			log.Info("CallPush [%v] Succ, Platform: %s, Delete Key", resID)
			if err = s.pushDao.ZRem(ctx, resID); err != nil {
				continue
			}
		}
	}
}

// distinguish the resource's platform info
func (s *Service) platform(resID string) (platform string) {
	var (
		err          error
		ios, android bool
		mobiAPPs     []string
	)
	platform = _platAll // default value
	if mobiAPPs, err = s.pushDao.Platform(ctx, resID); err != nil {
		return
	}
	for _, value := range mobiAPPs {
		switch value {
		case _vAndroid:
			android = true
		case _vAndroidB:
			android = true
		case _vIPad:
			ios = true
		case _vIPhone:
			ios = true
		default:
			log.Error("ResourceID %d, Limit Wrong Value %s", resID, value)
		}
	}
	if ios && !android {
		return _platIOS
	}
	if !ios && android {
		return _platAndroid
	}
	return // other case like all or none, just return the default value
}

func (s *Service) pushproc() {
	var (
		resIDs map[string]int64
		err    error
	)
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("DB closed!")
			return
		}
		time.Sleep(time.Duration(s.c.Cfg.Push.Fre))
		// pick to push resIDs from redis
		if resIDs, err = s.pushDao.ZrangeList(ctx); err != nil {
			log.Error("Get ToPush List Err %v", err)
			continue
		}
		if len(resIDs) == 0 {
			log.Info("No ToPush Data, Sleep")
			continue
		}
		// push the data
		log.Info("ToPush Treat Data: %d", len(resIDs))
		s.push(resIDs)
	}
}
