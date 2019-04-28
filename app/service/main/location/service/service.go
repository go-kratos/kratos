package service

import (
	"context"
	"path"
	"time"

	accoutCli "go-common/app/service/main/account/api"
	"go-common/app/service/main/location/conf"
	"go-common/app/service/main/location/dao"
	"go-common/app/service/main/location/model"
	"go-common/library/log"
	"go-common/library/stat/prom"

	ipdb "github.com/ipipdotnet/ipdb-go"
	maxminddb "github.com/oschwald/maxminddb-golang"
	"github.com/pkg/errors"
)

// Service define resource service
type Service struct {
	c              *conf.Config
	zdb            *dao.Dao
	accountSvc     accoutCli.AccountClient
	anonym         *maxminddb.Reader
	anonymFileName string
	// cache
	policy      map[int64]map[int64]int64
	groupPolicy map[int64][]int64
	// groupid by zone_id
	groupAuthZone map[int64]map[int64]map[int64]int64
	missch        chan interface{}
	// prom
	missedPorm   *prom.Prom
	authCodePorm *prom.Prom
	innerIPPorm  *prom.Prom
	// new ip library
	v4       *ipdb.City
	v6       *ipdb.City
	version4 string
	version6 string
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		zdb:            dao.New(c),
		missch:         make(chan interface{}, 1024),
		policy:         make(map[int64]map[int64]int64),
		groupPolicy:    make(map[int64][]int64),
		groupAuthZone:  make(map[int64]map[int64]map[int64]int64),
		missedPorm:     prom.New().WithCounter("auth_missed", []string{"method"}),
		authCodePorm:   prom.New().WithCounter("auth_code", []string{"method"}),
		innerIPPorm:    prom.New().WithCounter("inner_ipaddr", []string{"method"}),
		anonymFileName: c.AnonymFileName,
	}
	accountSvc, err := accoutCli.NewClient(c.Account)
	if err != nil {
		panic(err)
	}
	s.accountSvc = accountSvc
	if err := s.loadIPs(); err != nil {
		panic(err)
	}
	if err := s.loadAnonymous(); err != nil {
		panic(err)
	}
	if err := s.LoadPolicy(); err != nil {
		panic(err)
	}
	go s.reloadproc()
	go s.cacheproc()
	go s.reloadanonymous()
	return
}

// loadIPs load ip from ip.txt.
func (s *Service) loadIPs() (err error) {
	var (
		v4File  = path.Join(s.c.FilePath, s.c.IPv4Name)
		v6File  = path.Join(s.c.FilePath, s.c.IPv6Name)
		tmpV4   *ipdb.City
		tmpV6   *ipdb.City
		verison *model.Version
	)
	if verison, err = s.zdb.CheckVersion(context.Background()); err != nil || verison == nil {
		log.Error("%v or version is nil", err)
		return
	}
	if s.version4 != verison.StableV4 {
		log.Info("star down IPv4 library verson(%v) filename(%v)", verison.StableV4, v4File)
		if err = s.zdb.DownIPLibrary(context.Background(), verison.StableV4, v4File); err != nil {
			log.Error("error(%v) verson(%v) filename(%v)", err, verison.StableV4, v4File)
			return
		}
		log.Info("success down IPv4 library")
		log.Info("start load IPv4 library filename(%s)!", v4File)
		if tmpV4, err = ipdb.NewCity(v4File); err != nil {
			log.Error("load new IPv4 library error(%v)", err)
			err = errors.WithStack(err)
			return
		}
		s.v4 = tmpV4
		s.version4 = verison.StableV4
		log.Info("success load IPv4 library")
	}
	if s.version6 != verison.StableV6 {
		log.Info("star down IPv6 library verson(%v) filename(%v)", verison.StableV6, v6File)
		if err = s.zdb.DownIPLibrary(context.Background(), verison.StableV6, v6File); err != nil {
			log.Error("error(%v) verson(%v) filename(%v)", err, verison.StableV6, v6File)
			return
		}
		log.Info("success down IPv6 library")
		log.Info("start load IPv6 library filename(%s)!", v6File)
		if tmpV6, err = ipdb.NewCity(v6File); err != nil {
			err = errors.WithStack(err)
			return
		}
		s.v6 = tmpV6
		s.version6 = verison.StableV6
		log.Info("success load IPv6 library")
	}
	return
}

// reloadanonymous reload anonymous data.
func (s *Service) reloadIPs() {
	for {
		time.Sleep(time.Second * 86400)
		s.loadIPs()
	}
}

// loadAnonymous load anonymous ip.
func (s *Service) loadAnonymous() (err error) {
	log.Info("start down load Anonymous")
	if err = s.zdb.DownloadAnonym(); err != nil {
		log.Error("down load Anonymous faild error(%v)", err)
		return
	}
	log.Info("success down load Anonymous")
	var (
		anonymFile = path.Join(s.c.FilePath, s.c.AnonymFileName)
		tmpAnonym  *maxminddb.Reader
	)
	log.Info("start load Anonymous IP library filename(%s)!", anonymFile)
	if tmpAnonym, err = s.NewAnonym(anonymFile); err != nil {
		err = errors.WithStack(err)
		return
	}
	s.anonym = tmpAnonym
	log.Info("success load Anonymous IP library")
	return
}

// reloadanonymous reload anonymous data.
func (s *Service) reloadanonymous() {
	for {
		if time.Now().Weekday().String() == "Monday" && time.Now().Hour() == 0 {
			s.loadAnonymous()
		}
		time.Sleep(time.Minute * 60)
	}
}

// LoadPolicy locad policy from db
func (s *Service) LoadPolicy() (err error) {
	log.Info("start load policy cache !")
	var (
		tmpPolicy        map[int64]map[int64]int64
		tmpGroupPolicy   map[int64][]int64
		tmpGroupAuthZone map[int64]map[int64]map[int64]int64
	)
	log.Info("start to load s.policy")
	if tmpPolicy, err = s.zdb.Policies(context.TODO()); err != nil {
		return
	} else if len(tmpPolicy) > 0 {
		s.policy = tmpPolicy
	}
	log.Info("start to load s.groupPolicy")
	if tmpGroupPolicy, err = s.zdb.GroupPolicies(context.TODO()); err != nil {
		log.Error("s.groupPolicies error(%+v)", err)
	} else if len(tmpGroupPolicy) > 0 {
		s.groupPolicy = tmpGroupPolicy
	}
	log.Info("start to load s.zoneGroup")
	if tmpGroupAuthZone, err = s.zdb.GroupAuthZone(context.TODO()); err != nil {
		log.Error("s.GroupAuthZone error(%+v)", err)
	} else if len(tmpGroupAuthZone) > 0 {
		s.groupAuthZone = tmpGroupAuthZone
	}
	return
}

// reloadproc reload data from db.
func (s *Service) reloadproc() {
	for {
		time.Sleep(time.Minute * 10)
		s.LoadPolicy()
	}
}

// Close dao.
func (s *Service) Close() {}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.zdb.Ping(c)
}

func (s *Service) addCache(d interface{}) {
	// asynchronous add rules to redis
	select {
	case s.missch <- d:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for add rules into redis.
func (s *Service) cacheproc() {
	for {
		d := <-s.missch
		switch d.(type) {
		case map[int64]map[int64]int64:
			v := d.(map[int64]map[int64]int64)
			if err := s.zdb.AddAuth(context.TODO(), v); err != nil {
				log.Error("s.zdb.AddAuth(%v) error(%+v)", v, err)
			}
		default:
			log.Warn("cacheproc can't process the type")
		}
	}
}
