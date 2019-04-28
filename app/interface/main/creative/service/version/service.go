package version

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/creative"
	vsMdl "go-common/app/interface/main/creative/model/version"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"time"
)

//Service struct
type Service struct {
	c        *conf.Config
	Creative *creative.Dao
	// cache
	VersionCache    []*vsMdl.Version
	VersionMapCache map[string][]*vsMdl.Version
	WebManagerTip   *vsMdl.Version // WebTip
	AppManagerTip   *vsMdl.Version // APPTip
	CusManagerTip   *vsMdl.Version // 客服Tip
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:        c,
		Creative: creative.New(c),
	}
	s.loadVersion()
	go s.loadproc()
	return s
}

func (s *Service) loadVersion() {
	var (
		vss         []*vsMdl.Version
		vsWebTip    *vsMdl.Version
		vsAppTip    *vsMdl.Version
		vsCustomTip *vsMdl.Version
		err         error
	)
	vss, err = s.Creative.AllByTypes(context.TODO(), vsMdl.FullVersions())
	if err != nil {
		log.Error("s.Version.versions error(%v)", err)
		return
	}
	s.VersionCache = vss
	s.VersionMapCache, _ = s.versionMap(context.TODO())
	vsWebTip, err = s.Creative.LatestByType(context.TODO(), 7)
	if err != nil {
		log.Error("s.Creative.LatestByType, type=7, error(%v)", err)
		return
	}
	s.WebManagerTip = vsWebTip
	vsAppTip, err = s.Creative.LatestByType(context.TODO(), 8)
	if err != nil {
		log.Error("s.Creative.LatestByType, type=8, error(%v)", err)
		return
	}
	s.AppManagerTip = vsAppTip
	vsCustomTip, err = s.Creative.LatestByType(context.TODO(), 9)
	if err != nil {
		log.Error("s.Creative.LatestByType, type=9, error(%v)", err)
		return
	}
	s.CusManagerTip = vsCustomTip
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(2 * time.Minute)
		s.loadVersion()
	}
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.Creative.Ping(c); err != nil {
		log.Error("s.versionDao.PingDb err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.Creative.Close()
}
