package static

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	eggdao "go-common/app/interface/main/app-resource/dao/egg"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/static"
	"go-common/library/ecode"
	"go-common/library/log"

	farm "github.com/dgryski/go-farm"
)

const (
	_initVersion = "static_version"
)

var (
	_emptyStatics = []*static.Static{}
)

// Service static service.
type Service struct {
	dao        *eggdao.Dao
	tick       time.Duration
	cache      map[int8][]*static.Static
	staticPath string
}

// New new a static service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao:        eggdao.New(c),
		tick:       time.Duration(c.Tick),
		cache:      map[int8][]*static.Static{},
		staticPath: c.StaticJSONFile,
	}
	now := time.Now()
	s.loadCache(now)
	go s.loadCachepro()
	return
}

// Static return statics
func (s *Service) Static(plat int8, build int, ver string, now time.Time) (res []*static.Static, version string, err error) {
	var (
		tmps = s.cache[plat]
	)
	for _, tmp := range tmps {
		if model.InvalidBuild(build, tmp.Build, tmp.Condition) {
			continue
		}
		res = append(res, tmp)
	}
	if len(res) == 0 {
		res = _emptyStatics
	}
	if version = s.hash(res); version == ver {
		err = ecode.NotModified
		res = nil
	}
	return
}

func (s *Service) hash(v []*static.Static) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return _initVersion
	}
	return strconv.FormatUint(farm.Hash64(bs), 10)
}

// loadCache update egg
func (s *Service) loadCache(now time.Time) {
	tmp, err := s.dao.Egg(context.TODO(), now)
	if err != nil {
		log.Error("s.dao.Egg error(%v)", err)
		return
	}
	s.cache = tmp
}

func (s *Service) loadCachepro() {
	for {
		time.Sleep(s.tick)
		s.loadCache(time.Now())
	}
}
