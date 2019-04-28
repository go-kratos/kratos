package module

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	moduledao "go-common/app/interface/main/app-resource/dao/module"
	"go-common/app/interface/main/app-resource/model/module"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emptylist = []*module.ResourcePool{}
)

// Service module service.
type Service struct {
	dao             *moduledao.Dao
	tick            time.Duration
	resourceCache   map[string]*module.ResourcePool
	conditionsCache map[int]*module.Condition
}

// New new a module service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao:             moduledao.New(c),
		tick:            time.Duration(c.Tick),
		resourceCache:   map[string]*module.ResourcePool{},
		conditionsCache: make(map[int]*module.Condition),
	}
	s.loadCache()
	go s.loadproc()
	return
}

func (s *Service) FormCondition(versions []*module.Versions) (res map[string]map[string]int) {
	res = make(map[string]map[string]int)
	for _, pools := range versions {
		var (
			re map[string]int
			ok bool
		)
		for _, resource := range pools.Resource {
			if re, ok = res[pools.PoolName]; !ok {
				re = make(map[string]int)
				res[pools.PoolName] = re
			}
			var tmpVer int
			switch tmp := resource.Version.(type) {
			case string:
				tmpVer, _ = strconv.Atoi(tmp)
			case float64:
				tmpVer = int(tmp)
			}
			re[resource.ResourceName] = tmpVer
		}
	}
	return
}

// List get All or by poolname
func (s *Service) List(c context.Context, mobiApp, device, platform, poolName, env string, build, sysver, level, scale, arch int,
	versions []*module.Versions, now time.Time) (res []*module.ResourcePool) {
	var (
		resTmp      []*module.ResourcePool
		versionsMap = s.FormCondition(versions)
	)
	if poolName != "" {
		if pool, ok := s.resourceCache[poolName]; ok {
			resTmp = append(resTmp, pool)
		}
	} else {
		for _, l := range s.resourceCache {
			resTmp = append(resTmp, l)
		}
	}
	if len(resTmp) == 0 {
		res = _emptylist
		return
	}
	for _, resPool := range resTmp {
		var (
			existRes      = map[string]*module.Resource{}
			existResTotal = map[string]struct{}{}
			resPoolTmp    = &module.ResourcePool{Name: resPool.Name}
			ok            bool
		)
		for _, re := range resPool.Resources {
			if re == nil {
				continue
			}
			if !s.checkCondition(c, mobiApp, device, platform, env, build, sysver, level, scale, arch, re.Condition, now) {
				continue
			}
			var t *module.Resource
			if _, ok = existResTotal[re.Name]; ok {
				continue
			}
			if t, ok = existRes[re.Name]; ok {
				if re.Increment == module.Total {
					tmp := &module.Resource{}
					*tmp = *t
					tmp.TotalMD5 = re.MD5
					existResTotal[tmp.Name] = struct{}{}
					resPoolTmp.Resources = append(resPoolTmp.Resources, tmp)
					continue
				}
			} else {
				var (
					resVer map[string]int
					ver    int
				)
				if resVer, ok = versionsMap[resPool.Name]; ok {
					if ver, ok = resVer[re.Name]; ok {
						if re.Increment == module.Incremental && re.FromVer != ver {
							continue
						}
					} else if !ok && re.Increment == module.Incremental {
						continue
					}
				} else if !ok && re.Increment == module.Incremental {
					continue
				}
				tmp := &module.Resource{}
				*tmp = *re
				existRes[tmp.Name] = tmp
				if re.Increment == module.Total {
					tmp.TotalMD5 = re.MD5
					existResTotal[tmp.Name] = struct{}{}
					resPoolTmp.Resources = append(resPoolTmp.Resources, tmp)
				}
			}
		}
		if len(resPoolTmp.Resources) == 0 {
			continue
		}
		res = append(res, resPoolTmp)
	}
	return
}

// Resource get by poolname and resourcename
func (s *Service) Resource(c context.Context, mobiApp, device, platform, poolName, resourceName, env string,
	ver, build, sysver, level, scale, arch int, now time.Time) (res *module.Resource, err error) {
	if resPoolTmp, ok := s.resourceCache[poolName]; ok {
		if resPoolTmp == nil {
			err = ecode.NothingFound
			return
		}
		var (
			resTmp   *module.Resource
			existRes = map[string]struct{}{}
		)
		for _, resTmp = range resPoolTmp.Resources {
			if resTmp == nil {
				continue
			}
			if resTmp != nil && resTmp.Name == resourceName {
				if !s.checkCondition(c, mobiApp, device, platform, env, build, sysver, level, scale, arch, resTmp.Condition, now) {
					continue
				}
				if ver == 0 {
					if resTmp.Increment == module.Incremental {
						continue
					}
				} else {
					if resTmp.Increment == module.Incremental && resTmp.FromVer != ver {
						continue
					}
				}
				if resTmp.Increment == module.Total && resTmp.Version == ver {
					err = ecode.NotModified
					break
				}
				if _, ok := existRes[resTmp.Name]; !ok {
					res = &module.Resource{}
					*res = *resTmp
					existRes[resTmp.Name] = struct{}{}
				}
				if resTmp.Increment == module.Total && res != nil {
					res.TotalMD5 = resTmp.MD5
					break
				}
			}
		}
	}
	if err != nil {
		return
	}
	if res == nil {
		err = ecode.NothingFound
	}
	return
}

func (s *Service) checkCondition(c context.Context, mobiApp, device, platform, env string, build, sysver, level, scale, arch int, condition *module.Condition, now time.Time) bool {
	if condition == nil {
		return true
	}
	if env == module.EnvRelease && condition.Valid == 0 {
		return false
	} else if env == module.EnvTest && condition.ValidTest == 0 {
		return false
	} else if env == module.EnvDefault && condition.Default != 1 {
		return false
	}
	if !condition.STime.Time().IsZero() && now.Unix() < int64(condition.STime) {
		return false
	}
	if !condition.ETime.Time().IsZero() && now.Unix() > int64(condition.ETime) {
		return false
	}
NETX:
	for column, cv := range condition.Columns {
		switch column {
		case "plat": // whith list
			for _, v := range cv {
				if strings.TrimSpace(v.Value) == platform {
					continue NETX
				}
			}
			return false
		case "mobi_app": // whith list
			for _, v := range cv {
				if strings.TrimSpace(v.Value) == mobiApp {
					continue NETX
				}
			}
			return false
		case "device": // blace list
			for _, v := range cv {
				if strings.TrimSpace(v.Value) == device {
					return false
				}
			}
		case "build": // build < lt  gt > build ge >= build, le <= build
			for _, v := range cv {
				value, _ := strconv.Atoi(strings.TrimSpace(v.Value))
				if invalidModelBuild(build, value, v.Condition) {
					return false
				}
			}
		case "sysver":
			if sysver > 0 {
				for _, v := range cv {
					value, _ := strconv.Atoi(strings.TrimSpace(v.Value))
					if invalidModelBuild(sysver, value, v.Condition) {
						return false
					}
				}
			}
		case "scale": // whith list
			if scale > 0 {
				for _, v := range cv {
					value, _ := strconv.Atoi(strings.TrimSpace(v.Value))
					if value == scale {
						continue NETX
					}
				}
				return false
			}
		case "arch": // whith list
			if arch > 0 {
				for _, v := range cv {
					value, _ := strconv.Atoi(strings.TrimSpace(v.Value))
					if value == arch {
						continue NETX
					}
				}
				return false
			}
		}
	}
	return true
}

// ModuleUpdateCache update module cache
func (s *Service) ModuleUpdateCache() (err error) {
	err = s.loadCache()
	return
}

// load update cache
func (s *Service) loadCache() (err error) {
	configsTmp, err := s.dao.ResourceConfig(context.TODO())
	if err != nil {
		log.Error("s.dao.ResourceConfig error(%v)", err)
		return
	}
	limitTmp, err := s.dao.ResourceLimit(context.TODO())
	if err != nil {
		log.Error("s.dao.ResourceLimit error(%v)", err)
		return
	}
	for _, config := range configsTmp {
		if limit, ok := limitTmp[config.ID]; ok {
			config.Columns = limit
		}
	}
	s.conditionsCache = configsTmp
	log.Info("module conditions success")
	tmpResourceDev, err := s.dao.ModuleDev(context.TODO())
	if err != nil {
		log.Error("s.dao.ModuleDev error(%v)", err)
		return
	}
	tmpResources, err := s.dao.ModuleAll(context.TODO())
	if err != nil {
		log.Error("s.dao.ModuleAll error(%v)", err)
		return
	}
	tmpResourcePoolCaches := map[string]*module.ResourcePool{}
	for _, resPool := range tmpResourceDev {
		if resPool == nil {
			continue
		}
		var tmpResourcePoolCache = &module.ResourcePool{ID: resPool.ID, Name: resPool.Name}
		for _, res := range resPool.Resources {
			if res == nil {
				continue
			}
			if re, ok := tmpResources[res.ID]; ok {
				var tmpre []*module.Resource
				for _, r := range re {
					if r.URL == "" || r.MD5 == "" {
						continue
					}
					if c, ok := s.conditionsCache[r.ResID]; ok {
						r.Condition = c
						// all level
						if c != nil {
							for column, cv := range c.Columns {
								switch column {
								case "level":
									for _, v := range cv {
										value, _ := strconv.Atoi(strings.TrimSpace(v.Value))
										r.Level = value
									}
								}
							}
						}
						r.IsWifi = c.IsWifi
					}
					tmpre = append(tmpre, r)
				}
				if len(tmpre) == 0 {
					continue
				}
				tmpResourcePoolCache.Resources = append(tmpResourcePoolCache.Resources, tmpre...)
			}
		}
		tmpResourcePoolCaches[resPool.Name] = tmpResourcePoolCache
	}
	s.resourceCache = tmpResourcePoolCaches
	log.Info("module resources success")
	return
}

// cacheproc load cache data
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.loadCache()
	}
}

// invalidModelBuild model build
func invalidModelBuild(srcBuild, cfgBuild int, cfgCond string) bool {
	if cfgBuild != 0 && cfgCond != "" {
		switch cfgCond {
		case "lt":
			if cfgBuild <= srcBuild {
				return true
			}
		case "le":
			if cfgBuild < srcBuild {
				return true
			}
		case "ge":
			if cfgBuild > srcBuild {
				return true
			}
		case "gt":
			if cfgBuild >= srcBuild {
				return true
			}
		}
	}
	return false
}
