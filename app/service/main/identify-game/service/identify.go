package service

import (
	"context"
	"errors"
	"strings"

	"go-common/app/service/main/identify-game/api/grpc/v1"
	"go-common/app/service/main/identify-game/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_segmentation = "_"
)

var (
	_noLogin = &model.AccessInfo{
		Mid:     -1,
		Expires: 60,
	}

	// ErrDispatcherError dispatcher error
	ErrDispatcherError = errors.New("dispatcher route map is error")
)

// Oauth verify user info by accesskey.
func (s *Service) Oauth(c context.Context, accesskey, from string) (res *model.AccessInfo, err error) {
	cache := true
	if len(accesskey) > 32 {
		cache = false
	}
	if cache {
		if res, err = s.d.AccessCache(c, accesskey); err != nil {
			cache = false
		}
	}
	if res != nil {
		if res.Mid == _noLogin.Mid {
			err = ecode.NoLogin
		}
		return
	}
	target := s.target(accesskey)
	if from != "" && target != s.c.Dispatcher.Name {
		s.dispatcherErrStats.Incr("dispatcher_error")
		err = ErrDispatcherError
		log.Error("Oauth dispatcher routMap is error. token:%s, from:%s", accesskey, from)
		return
	}
	if res, err = s.d.AccessToken(c, accesskey, target); err != nil {
		ec := ecode.Cause(err)
		if ec != ecode.NoLogin && ec != ecode.AccessKeyErr {
			return
		}
		if cache {
			s.addCache(func() {
				s.d.SetAccessCache(context.TODO(), accesskey, _noLogin)
			})
		}
		return
	}
	if cache && res != nil {
		s.addCache(func() {
			s.d.SetAccessCache(context.TODO(), accesskey, res)
		})
	}
	return
}

// RenewToken prolong user accesskey.
func (s *Service) RenewToken(c context.Context, accesskey, from string) (res *model.RenewInfo, err error) {
	target := s.target(accesskey)
	if from != "" && target != s.c.Dispatcher.Name {
		s.dispatcherErrStats.Incr("dispatcher_error")
		err = ErrDispatcherError
		log.Error("RenewToken dispatcher routMap is error. token:%s, from:%s", accesskey, from)
		return
	}
	return s.d.RenewToken(c, accesskey, target)
}

func (s *Service) target(accesskey string) (res string) {
	index := strings.Index(accesskey, _segmentation)
	if index < 0 {
		res = s.c.Dispatcher.Name
		return
	}
	res = accesskey[index+1:]
	return
}

// DelCache for clean cache
func (s *Service) DelCache(c context.Context, accesskey string) error {
	return s.d.DelAccessCache(c, accesskey)
}

// GetCookieByToken get mid by token, get cookie by Cookie
func (s *Service) GetCookieByToken(c context.Context, accesskey, from string) (cookies *v1.CreateCookieReply, err error) {
	accInfo, err := s.Oauth(c, accesskey, from)
	if err != nil {
		return
	}
	cookies, err = s.d.GetCookieByMid(c, accInfo.Mid)
	return
}
