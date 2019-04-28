package app

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/app"
	"go-common/app/interface/main/creative/dao/message"
	appDML "go-common/app/interface/main/creative/model/app"
	msgDML "go-common/app/interface/main/creative/model/message"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
)

//Service struct.
type Service struct {
	c                         *conf.Config
	app                       *app.Dao
	msg                       *message.Dao
	PortalIntro, PortalNotice []*appDML.PortalMeta
	CameraCfg                 map[string]interface{}
	p                         *service.Public
}

func (s *Service) initCameraCfg() {
	s.CameraCfg = map[string]interface{}{
		"videoup_min_sec": 5,
		"videoup_max_sec": 300,
		"dyna_min_sec":    5,
		"dyna_max_sec":    300,
		"coo_min_sec":     5,
		"coo_max_sec":     300,
	}
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:            c,
		msg:          message.New(c),
		app:          app.New(c),
		PortalIntro:  make([]*appDML.PortalMeta, 0),
		PortalNotice: make([]*appDML.PortalMeta, 0),
		p:            p,
	}
	s.initCameraCfg()
	s.loadPortal()
	go s.loadproc()
	return s
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.app.Ping(c); err != nil {
		log.Error("s.appDao.PingDb err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.app.Close()
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(2 * time.Minute)
		s.loadPortal()
	}
}

//load db
func (s *Service) loadPortal() {
	intro, err := s.app.Portals(context.TODO(), appDML.PortalIntro)
	if err != nil {
		log.Error("s.app.intro error(%v)", err)
		return
	}
	s.PortalIntro = intro
	// 创作激励 + 征稿公告
	notice, err := s.app.Portals(context.TODO(), appDML.PortalIntro)
	if err != nil {
		log.Error("s.app.notice error(%v)", err)
		return
	}
	s.PortalNotice = notice
}

// TopMsg fn.
func (s *Service) TopMsg(c context.Context, mid int64, build int, os, app, ak, ck, ip string) (data []*msgDML.Message, err error) {
	if build > 5332000 && os == "android" {
		return
	} else if build > 8220 && app == "iphone" {
		return
	} else if build > 7339 && build < 8000 && app == "iphone_b" {
		return
	}
	data, err = s.msg.GetUpList(c, mid, ak, ck, ip)
	topLen := 1
	if len(data) > topLen {
		data = data[:topLen]
	}
	return
}
