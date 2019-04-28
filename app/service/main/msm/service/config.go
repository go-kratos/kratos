package service

import (
	"context"

	"go-common/app/infra/config/model"
	"go-common/library/log"
)

// Push push new ver to config-service.
func (s *Service) Push(c context.Context, app, bver, env string, ver int64) (err error) {
	arg := &model.ArgConf{
		App:      app,
		BuildVer: bver,
		Ver:      ver,
		Env:      env,
	}
	if err = s.confSvr.Push(c, arg); err != nil {
		log.Error("push(%v) error(%v)", arg, err)
	}
	return
}

// SetToken set  token to config-service.
func (s *Service) SetToken(c context.Context, app, env, token string) (err error) {
	arg := &model.ArgToken{
		App:   app,
		Token: token,
		Env:   env,
	}
	if err = s.confSvr.SetToken(c, arg); err != nil {
		log.Error("SetToken(%v) error(%v)", arg, err)
	}
	return
}
