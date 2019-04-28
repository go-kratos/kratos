package notice

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	locdao "go-common/app/interface/main/app-resource/dao/location"
	ntcdao "go-common/app/interface/main/app-resource/dao/notice"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/notice"
	locmdl "go-common/app/service/main/location/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"go-common/library/net/metadata"

	"github.com/dgryski/go-farm"
)

const (
	_initNoticeKey = "notice_key_%d_%d"
	_initNoticeVer = "notice_version"
)

var (
	_emptyNotice = &notice.Notice{}
)

// Service notice service.
type Service struct {
	dao *ntcdao.Dao
	loc *locdao.Dao
	// tick
	tick time.Duration
	// cache
	cache map[string][]*notice.Notice
}

// New new a notice service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: ntcdao.New(c),
		loc: locdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// cache
		cache: map[string][]*notice.Notice{},
	}
	s.load(time.Now())
	go s.loadproc()
	return
}

// Notice return Notice to json
func (s *Service) Notice(c context.Context, plat int8, build, typeInt int, ver string) (res *notice.Notice, version string, err error) {
	var (
		ip    = metadata.String(c, metadata.RemoteIP)
		pids  []string
		auths map[string]*locmdl.Auth
	)
	for _, ntc := range s.cache[fmt.Sprintf(_initNoticeKey, plat, typeInt)] {
		if model.InvalidBuild(build, ntc.Build, ntc.Condition) {
			continue
		}
		if ntc.Area != "" {
			pids = append(pids, ntc.Area)
		}
	}
	if len(pids) > 0 {
		auths, _ = s.loc.AuthPIDs(c, strings.Join(pids, ","), ip)
	}
	for _, ntc := range s.cache[fmt.Sprintf(_initNoticeKey, plat, typeInt)] {
		if model.InvalidBuild(build, ntc.Build, ntc.Condition) {
			continue
		}
		if auth, ok := auths[ntc.Area]; ok && auth.Play == locmdl.Forbidden {
			log.Warn("s.invalid area(%v) ip(%v) error(%v)", ntc.Area, ip, err)
			continue
		}
		res = ntc
		break
	}
	if res == nil {
		res = _emptyNotice
	}
	if version = s.hash(res); ver == version {
		err = ecode.NotModified
		res = nil
	}
	return
}

// load
func (s *Service) load(now time.Time) {
	// get notice
	ntcs, err := s.dao.All(context.TODO(), now)
	if err != nil {
		log.Error("s.dao.GetAll() error(%v)", err)
		return
	}
	// copy cache
	tmp := map[string][]*notice.Notice{}
	for _, v := range ntcs {
		key := fmt.Sprintf(_initNoticeKey, v.Plat, v.Type)
		tmp[key] = append(tmp[key], v)
	}
	s.cache = tmp
	log.Info("notice cacheproc success")
}

func (s *Service) hash(v *notice.Notice) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return _initNoticeVer
	}
	return strconv.FormatUint(farm.Hash64(bs), 10)
}

// cacheproc load cache data
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.load(time.Now())
	}
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
