package region

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/region"
	locmdl "go-common/app/service/main/location/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	farm "github.com/dgryski/go-farm"
)

const (
	_initRegionKey   = "region_key_%d_%v"
	_initlanguage    = "hans"
	_initVersion     = "region_version"
	_regionAll       = int8(0)
	_isRegion        = int8(1)
	_isRank          = int8(2)
	_initRegionlimit = "%d_%v"
	_regionRepeat    = "r_%d_%d"
)

var (
	_emptyRegions = []*region.Region{}
	_isBangumi    = map[int]struct{}{
		13:  struct{}{},
		177: struct{}{},
		23:  struct{}{},
		11:  struct{}{},
	}
	_isBangumiIndex = map[int]struct{}{
		13:  struct{}{},
		23:  struct{}{},
		11:  struct{}{},
		167: struct{}{},
	}
	_regionlimit = map[int8]map[string]map[int]string{
		model.PlatIPhone: map[string]map[int]string{
			"65542_bilibili://cliparea": map[int]string{
				5960: "gt",
				6570: "lt",
			},
			"65541_bilibili://category/65541": map[int]string{
				5960: "gt",
				6570: "lt",
			},
			"65544_bilibili://albumarea": map[int]string{
				6090: "gt",
				6570: "lt",
			},
		},
	}
)

// Regions get regions.
func (s *Service) Regions(c context.Context, plat int8, build int, ver, mobiApp, device, language string) (rs []*region.Region, version string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	rs, version, err = s.getCache(c, plat, build, ver, ip, mobiApp, device, language, "", false)
	return
}

// Regions get regions.
func (s *Service) RegionsList(c context.Context, plat int8, build int, ver, mobiApp, device, language, entrance string) (rs []*region.Region, version string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	rs, version, err = s.getCache(c, plat, build, ver, ip, mobiApp, device, language, entrance, true)
	return
}

// getCache get region from cache.
func (s *Service) getCache(c context.Context, plat int8, build int, ver, ip, mobiApp, device, language, entrance string, isNew bool) (res []*region.Region, version string, err error) {
	if language == "" {
		language = _initlanguage
	}
	var (
		rs           = s.cache[fmt.Sprintf(_initRegionKey, plat, language)]
		child        = map[int][]*region.Region{}
		entranceShow = _isRegion
		ridlimit     = _regionlimit[plat]
		ridtmp       = map[string]struct{}{}
		pids         []string
		auths        map[string]*locmdl.Auth
	)
	switch entrance {
	case "region":
		entranceShow = _isRegion
	case "rank":
		entranceShow = _isRank
	}
	for _, rtmp := range rs {
		if rtmp.Area != "" {
			pids = append(pids, rtmp.Area)
		}
	}
	if len(pids) > 0 {
		auths, _ = s.loc.AuthPIDs(c, strings.Join(pids, ","), ip)
	}
Retry:
	for _, rtmp := range rs {
		r := &region.Region{}
		*r = *rtmp
		if _, isgbm := _isBangumi[r.Rid]; isgbm {
			r.IsBangumi = 1
		}
		if isNew {
			if r.Entrance != _regionAll && entranceShow != r.Entrance {
				continue
			}
		} else {
			switch r.Rid {
			case 65545, 65542, 65541, 65543, 65544, 65546:
				if mobiApp == "android" {
					continue
				}
			}
		}
		if r.Rid != 165 || ((mobiApp != "iphone" || device != "pad") || build <= 3590) {
			if model.InvalidBuild(build, r.Build, r.Condition) {
				continue
			}
		}
		key := fmt.Sprintf(_initRegionlimit, r.Rid, r.URI)
		if rlimit, ok := ridlimit[key]; ok {
			for blimit, climit := range rlimit {
				if model.InvalidBuild(build, blimit, climit) {
					continue Retry
				}
			}
		}
		if r.Rid == 65541 && (plat == model.PlatIPhone && build == 7040) {
			continue
		}
		if r.Rid == 65543 && ((plat == model.PlatIPhone && (build == 7070 || build == 7040 || build == 7030)) ||
			(plat == model.PlatAndroid && (build == 591182 || build == 591181 || build == 591178 || build == 591177))) {
			continue
		}
		if auth, ok := auths[r.Area]; ok && auth.Play == locmdl.Forbidden {
			log.Warn("s.invalid area(%v) ip(%v) error(%v)", r.Area, ip, err)
			continue
		}
		if isAudit := s.auditRegion(mobiApp, plat, build, r.Rid); isAudit {
			continue
		}
		rkey := fmt.Sprintf(_regionRepeat, r.Rid, r.Reid)
		if _, ok := ridtmp[rkey]; !ok {
			ridtmp[rkey] = struct{}{}
		} else {
			continue
		}
		if r.Reid != 0 {
			cl, ok := child[r.Reid]
			if !ok {
				cl = []*region.Region{}
			}
			cl = append(cl, r)
			child[r.Reid] = cl
		} else {
			res = append(res, r)
		}
	}
	if len(res) == 0 {
		res = _emptyRegions
	} else {
		for _, r := range res {
			r.Children = child[r.Rid]
		}
	}
	if version = s.hash(res); version == ver {
		err = ecode.NotModified
		res = nil
	}
	return
}

// NewRegionList get region from cache.
func (s *Service) NewRegionList(c context.Context, plat int8, build int, ver, mobiApp, device, language string) (res []*region.Region, version string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		hantlanguage = "hant"
	)
	if ok := model.IsOverseas(plat); ok && language != _initlanguage && language != hantlanguage {
		language = hantlanguage
	} else if language == "" {
		language = _initlanguage
	}
	var (
		rs     = s.cachelist[fmt.Sprintf(_initRegionKey, plat, language)]
		child  = map[int][]*region.Region{}
		ridtmp = map[string]struct{}{}
		pids   []string
		auths  map[string]*locmdl.Auth
	)
	for _, rtmp := range rs {
		if rtmp.Area != "" {
			pids = append(pids, rtmp.Area)
		}
	}
	if len(pids) > 0 {
		auths, _ = s.loc.AuthPIDs(c, strings.Join(pids, ","), ip)
	}
LOOP:
	for _, rtmp := range rs {
		r := &region.Region{}
		*r = *rtmp
		if _, isgbm := _isBangumiIndex[r.Rid]; isgbm {
			r.IsBangumi = 1
		}
		var tmpl, limitshow bool
		if limit, ok := s.limitCache[r.ID]; ok {
			for i, l := range limit {
				if i+1 <= len(limit)-1 {
					if ((l.Condition == "gt" && limit[i+1].Condition == "lt") && (l.Build < limit[i+1].Build)) ||
						((l.Condition == "lt" && limit[i+1].Condition == "gt") && (l.Build > limit[i+1].Build)) {
						if (l.Condition == "gt" && limit[i+1].Condition == "lt") &&
							(build > l.Build && build < limit[i+1].Build) {
							break
						} else if (l.Condition == "lt" && limit[i+1].Condition == "gt") &&
							(build < l.Build && build > limit[i+1].Build) {
							break
						} else {
							tmpl = true
							continue
						}
					}
				}
				if tmpl {
					if i == len(limit)-1 {
						limitshow = true
						// continue LOOP
						break
					}
					tmpl = false
					continue
				}
				if model.InvalidBuild(build, l.Build, l.Condition) {
					limitshow = true
					continue
					// continue LOOP
				} else {
					limitshow = false
					break
				}
			}
		}
		if limitshow {
			continue LOOP
		}
		if auth, ok := auths[r.Area]; ok && auth.Play == locmdl.Forbidden {
			log.Warn("s.invalid area(%v) ip(%v) error(%v)", r.Area, ip, err)
			continue
		}
		if isAudit := s.auditRegion(mobiApp, plat, build, r.Rid); isAudit {
			continue
		}
		if config, ok := s.configCache[r.ID]; ok {
			r.Config = config
		}
		key := fmt.Sprintf(_regionRepeat, r.Rid, r.Reid)
		if _, ok := ridtmp[key]; !ok {
			ridtmp[key] = struct{}{}
		} else {
			continue
		}
		if r.Reid != 0 {
			cl, ok := child[r.Reid]
			if !ok {
				cl = []*region.Region{}
			}
			cl = append(cl, r)
			child[r.Reid] = cl
		} else {
			res = append(res, r)
		}
	}
	if len(res) == 0 {
		res = _emptyRegions
	} else {
		for _, r := range res {
			r.Children = child[r.Rid]
		}
	}
	if version = s.hash(res); version == ver {
		err = ecode.NotModified
		res = nil
	}
	return
}

func (s *Service) hash(v []*region.Region) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return _initVersion
	}
	return strconv.FormatUint(farm.Hash64(bs), 10)
}

// loadRegion regions cache.
func (s *Service) loadRegion() {
	res, err := s.dao.All(context.TODO())
	if err != nil {
		log.Error("s.dao.All error(%v)", err)
		return
	}
	tmp := map[string][]*region.Region{}
	for _, v := range res {
		key := fmt.Sprintf(_initRegionKey, v.Plat, v.Language)
		tmp[key] = append(tmp[key], v)
	}
	if len(tmp) > 0 {
		s.cache = tmp
	}
	log.Info("region cacheproc success")
}

func (s *Service) loadRegionlist() {
	res, err := s.dao.AllList(context.TODO())
	if err != nil {
		log.Error("s.dao.All error(%v)", err)
		return
	}
	tmp := map[string][]*region.Region{}
	tmpRegion := map[string]map[int]*region.Region{}
	for _, v := range res {
		key := fmt.Sprintf(_initRegionKey, v.Plat, v.Language)
		tmp[key] = append(tmp[key], v)
		// region list map
		if r, ok := tmpRegion[key]; ok {
			r[v.Rid] = v
		} else {
			tmpRegion[key] = map[int]*region.Region{
				v.Rid: v,
			}
		}
	}
	if len(tmp) > 0 && len(tmpRegion) > 0 {
		s.cachelist = tmp
		s.regionListCache = tmpRegion
	}
	log.Info("region list cacheproc success")
	limit, err := s.dao.Limit(context.TODO())
	if err != nil {
		log.Error("s.dao.limit error(%v)", err)
		return
	}
	s.limitCache = limit
	log.Info("region limit cacheproc success")
	config, err := s.dao.Config(context.TODO())
	if err != nil {
		log.Error("s.dao.Config error(%v)", err)
		return
	}
	s.configCache = config
	log.Info("region config cacheproc success")
}
