package version

import (
	"context"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	verdao "go-common/app/interface/main/app-resource/dao/version"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/version"
	"go-common/library/ecode"
	"go-common/library/log"

	farm "github.com/dgryski/go-farm"
)

const (
	_defaultChannel = "bili"
)

var (
	_emptyVersion   = []*version.Version{}
	_emptyVersionSo = []*version.VersionSo{}
)

// Service version service.
type Service struct {
	dao          *verdao.Dao
	cache        map[int8][]*version.Version
	upCache      map[int8]map[string][]*version.VersionUpdate
	uplimitCache map[int][]*version.UpdateLimit
	soCache      map[string][]*version.VersionSo
	increCache   map[int8]map[string][]*version.Incremental
	rnCache      map[string]map[string]*version.Rn
	tick         time.Duration
}

// New new a version service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao:          verdao.New(c),
		tick:         time.Duration(c.Tick),
		cache:        map[int8][]*version.Version{},
		upCache:      map[int8]map[string][]*version.VersionUpdate{},
		uplimitCache: map[int][]*version.UpdateLimit{},
		soCache:      map[string][]*version.VersionSo{},
		increCache:   map[int8]map[string][]*version.Incremental{},
		rnCache:      map[string]map[string]*version.Rn{},
	}
	s.load()
	go s.loadproc()
	return
}

// Version return version
func (s *Service) Version(plat int8) (res []*version.Version, err error) {
	if res = s.cache[plat]; res == nil {
		res = _emptyVersion
	}
	return
}

// VersionUpdate return version
func (s *Service) VersionUpdate(build int, plat int8, buvid, sdkint, channel, module, oldID string) (res *version.VersionUpdate, err error) {
	var (
		gvu, tmp []*version.VersionUpdate
	)
	if vup, ok := s.upCache[plat]; ok {
		if tmp, ok = vup[channel]; !ok || len(tmp) == 0 {
			if plat == model.PlatAndroidTVYST {
				err = ecode.NotModified
				return
			}
			if tmp, ok = vup[_defaultChannel]; !ok || len(tmp) == 0 {
				err = ecode.NotModified
				return
			}
		}
		for _, t := range tmp {
			tu := &version.VersionUpdate{}
			*tu = *t
			gvu = append(gvu, tu)
		}
	} else {
		err = ecode.NotModified
		return
	}
LOOP:
	for _, vu := range gvu {
		if build >= vu.Build {
			err = ecode.NotModified
			return
		}
		if vu.IsGray == 1 {
			if len(vu.SdkIntList) > 0 {
				if _, ok := vu.SdkIntList[sdkint]; !ok {
					continue LOOP
				}
			}
			if module != vu.Model && vu.Model != "" {
				continue LOOP
			}
			if buvid != "" {
				id := farm.Hash32([]byte(buvid))
				n := int(id % 100)
				if vu.BuvidStart > n || n > vu.BuvidEnd {
					continue LOOP
				}
			}
		}
		if limit, ok := s.uplimitCache[vu.Id]; ok {
			var tmpl bool
			for i, l := range limit {
				if i+1 <= len(limit)-1 {
					if ((l.Conditions == "gt" && limit[i+1].Conditions == "lt") && (l.BuildLimit < limit[i+1].BuildLimit)) ||
						((l.Conditions == "lt" && limit[i+1].Conditions == "gt") && (l.BuildLimit > limit[i+1].BuildLimit)) {
						if (l.Conditions == "gt" && limit[i+1].Conditions == "lt") &&
							(build > l.BuildLimit && build < limit[i+1].BuildLimit) {
							res = vu
							break LOOP
						} else if (l.Conditions == "lt" && limit[i+1].Conditions == "gt") &&
							(build < l.BuildLimit && build > limit[i+1].BuildLimit) {
							res = vu
							break LOOP
						} else {
							tmpl = true
							continue
						}
					}
				}
				if tmpl {
					tmpl = false
					continue
				}
				if model.InvalidBuild(build, l.BuildLimit, l.Conditions) {
					continue
				} else {
					res = vu
					break LOOP
				}
			}
		} else {
			res = vu
			break LOOP
		}
	}
	if res == nil {
		err = ecode.NotModified
		return
	}
	res.Incre = s.versionIncrementals(plat, res.Build, oldID)
	return
}

// versionIncrementals version incrementals
func (s *Service) versionIncrementals(plat int8, build int, oldID string) (ver *version.Incremental) {
	if v, ok := s.increCache[plat]; ok {
		if vers, ok := v[oldID]; ok {
			for _, value := range vers {
				if value.Build == build {
					ver = value
					return
				}
			}
		}
	}
	return
}

// VersionSo return version_so
func (s *Service) VersionSo(build, seed, sdkint int, name, model string) (vsdesc *version.VersionSoDesc, err error) {
	vSo := s.soCache[name]
	if len(vSo) == 0 {
		err = ecode.NotModified
		return
	}
	vsdesc = &version.VersionSoDesc{
		Package:     vSo[0].Package,
		Name:        vSo[0].Name,
		Description: vSo[0].Description,
		Clear:       vSo[0].Clear,
	}
	for _, value := range vSo {
		if value.Min_build > build || (seed > 0 && seed%100 >= value.Coverage) || (sdkint != value.Sdkint && value.Sdkint != 0) || (model != value.Model && value.Model != "" && value.Model != "*") {
			continue
		}
		vsdesc.Versions = append(vsdesc.Versions, value)
	}
	if vsdesc.Versions == nil {
		vsdesc.Versions = _emptyVersionSo
	}
	return
}

// VersionRn return version_rn
func (s *Service) VersionRn(version, deploymentKey, bundleID string) (vrn *version.Rn, err error) {
	if v, ok := s.rnCache[deploymentKey]; ok {
		if vrn, ok = v[version]; !ok {
			err = ecode.NotModified
		} else if vrn.BundleID == bundleID {
			err = ecode.NotModified
		}
	} else {
		err = ecode.NotModified
	}
	return
}

// load cache data
func (s *Service) load() {
	ver, err := s.dao.All(context.TODO())
	if err != nil {
		log.Error("version s.dao.All() error(%v)", err)
		return
	}
	s.cache = ver
	upver, err := s.dao.Updates(context.TODO())
	if err != nil {
		log.Error("version s.dao.GetUpdate() error(%v)", err)
		return
	}
	s.upCache = upver
	log.Info("version cacheproc success")
	sover, err := s.dao.Sos(context.TODO())
	if err != nil {
		log.Error("version s.dao.Sos() error(%v)", err)
		return
	}
	uplimit, err := s.dao.Limits(context.TODO())
	if err != nil {
		log.Error("version s.dao.Limits() error(%v)", err)
		return
	}
	s.uplimitCache = uplimit
	s.soCache = sover
	log.Info("versionso cacheproc success")
	increver, err := s.dao.Incrementals(context.TODO())
	if err != nil {
		log.Error("version s.dao.Incrementals error(%v)", err)
		return
	}
	s.increCache = increver
	log.Info("versionIncre cacheproc success")
	rn, err := s.dao.Rn(context.TODO())
	if err != nil {
		log.Error("version s.dao.Rn error(%v)", err)
		return
	}
	s.rnCache = rn
	log.Info("versionRn cacheproc success")
}

// cacheproc load cache data
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.load()
	}
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
