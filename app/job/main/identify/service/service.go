package service

import (
	"context"

	"go-common/app/job/main/identify/conf"
	"go-common/app/job/main/identify/dao"
	"go-common/app/job/main/identify/model"
	"go-common/library/cache/memcache"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
)

const (
	_tokenTable   = "aso_app_perm"
	_cookieTable  = "aso_cookie_token"
	_insertAction = "insert"
	_delteAction  = "delete"
)

var (
	_gameAppID = [3]int64{432, 876, 849}
)

// Service is a identify service.
type Service struct {
	c           *conf.Config
	d           *dao.Dao
	identifySub *databus.Databus
	authDataBus *databus.Databus
	// mc
	poolm map[string]*memcache.Pool
	// databus group
	authGroup     *databusutil.Group
	identifyGroup *databusutil.Group

	cookieCh []chan *model.AuthCookie
	tokenCh  []chan *model.AuthToken
}

// New new a identify service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		d:           dao.New(c),
		identifySub: databus.New(c.DataBus.IdentifySub),
		authDataBus: databus.New(c.DataBus.AuthDataBus),
		cookieCh:    make([]chan *model.AuthCookie, c.CheckConf.ChanNum),
		tokenCh:     make([]chan *model.AuthToken, c.CheckConf.ChanNum),
	}
	if len(s.c.Memcaches) > 0 {
		pm := make(map[string]*memcache.Pool, len(s.c.Memcaches))
		for name, mcc := range s.c.Memcaches {
			p := memcache.NewPool(mcc.Config)
			pm[name] = p
		}
		s.poolm = pm
	}

	s.authGroup = databusutil.NewGroup(c.Databusutil, s.authDataBus.Messages())
	s.authGroup.New = s.new
	s.authGroup.Split = s.spilt
	s.authGroup.Do = s.processAuthBinlog2
	s.authGroup.Start()

	s.identifyGroup = databusutil.NewGroup(c.Databusutil, s.identifySub.Messages())
	s.identifyGroup.New = s.identifyNew
	s.identifyGroup.Split = s.identifySplit
	s.identifyGroup.Do = s.processIdentifyInfo
	s.identifyGroup.Start()

	if c.CheckConf.Switch {
		for i := 0; i < c.CheckConf.ChanNum; i++ {
			cookie := make(chan *model.AuthCookie, c.CheckConf.ChanSize)
			token := make(chan *model.AuthToken, c.CheckConf.ChanSize)
			s.cookieCh[i] = cookie
			s.tokenCh[i] = token
			go s.checkCookie(cookie)
			go s.checkToken(token)
		}
		go s.queryCookieDeleted()
		go s.queryTokenDeleted()
	}

	return
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return nil
}

// Close close.
func (s *Service) Close() (err error) {
	s.identifySub.Close()
	s.authDataBus.Close()
	s.authGroup.Close()
	s.identifyGroup.Close()
	return nil
}
