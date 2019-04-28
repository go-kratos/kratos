package service

import (
	"context"
	"fmt"

	"go-common/app/infra/config/model"
	"go-common/library/log"
)

// Push push new ver to config-service.
func (s *Service) Push(c context.Context, treeID int64, env, zone, bver string, ver int64) (err error) {
	svr := svrFmt(treeID, env, zone)
	arg := &model.ArgConf{
		App:      svr,
		BuildVer: bver,
		Ver:      ver,
	}
	if err = s.confSvr.PushV4(c, arg); err != nil {
		log.Error("PushV4(%v) error(%v)", arg, err)
	}
	return
}

// SetToken set token to config-service.
func (s *Service) SetToken(c context.Context, treeID int64, env, zone, token string) (err error) {
	svr := svrFmt(treeID, env, zone)
	arg := &model.ArgToken{
		App:   svr,
		Token: token,
	}
	if err = s.confSvr.SetTokenV4(c, arg); err != nil {
		log.Error("SetToken(%v) error(%v)", arg, err)
	}
	return
}

// Hosts get hosts.
func (s *Service) Hosts(c context.Context, treeID int64, appName, env, zone string) (hosts []*model.Host, err error) {
	svr := svrFmt(treeID, env, zone)
	if hosts, err = s.confSvr.Hosts(c, svr); err != nil {
		log.Error("Hosts(%v) error(%v)", svr, err)
		return
	}
	if len(hosts) == 0 {
		hosts = make([]*model.Host, 0)
		return
	}
	for _, host := range hosts {
		host.Service = appName
	}
	return
}

// ClearHost clear hosts.
func (s *Service) ClearHost(c context.Context, treeID int64, env, zone string) (err error) {
	svr := svrFmt(treeID, env, zone)
	if err = s.confSvr.ClearHost(c, svr); err != nil {
		log.Error("Hosts(%v) error(%v)", svr, err)
	}
	return
}

func svrFmt(treeID int64, env, zone string) string {
	return fmt.Sprintf("%d_%s_%s", treeID, env, zone)
}

// PushForce push new ver to config-service.
func (s *Service) PushForce(c context.Context, treeID int64, env, zone, bver string, ver int64, hosts map[string]string, sType int8) (err error) {
	svr := svrFmt(treeID, env, zone)
	arg := &model.ArgConf{
		App:      svr,
		BuildVer: bver,
		Ver:      ver,
		Env:      env,
		Hosts:    hosts,
		SType:    sType,
	}
	if err = s.confSvr.Force(c, arg); err != nil {
		log.Error("PushForce(%v) error(%v)", arg, err)
	}
	return
}
