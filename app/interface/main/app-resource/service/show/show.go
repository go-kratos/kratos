package show

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/abtest"
	"go-common/app/interface/main/app-resource/model/show"
	"go-common/library/ecode"
	"go-common/library/log"

	farm "github.com/dgryski/go-farm"
)

const (
	_initTabKey          = "tab_%d_%s"
	_initVersion         = "showtab_version"
	_defaultLanguageHans = "hans"
	_defaultLanguageHant = "hant"
)

var (
	_showAbtest = map[string]string{
		"bilibili://pegasus/hottopic": "home_tabbar_server_1",
	}
	_deafaultTab = map[string]*show.Tab{
		"bilibili://pegasus/promo": &show.Tab{
			DefaultSelected: 1,
		},
	}
)

// Tabs show tabs
func (s *Service) Tabs(c context.Context, plat int8, build int, buvid, ver, mobiApp, language string, mid int64) (res map[string][]*show.Tab, version string, a *abtest.List, err error) {
	if key := fmt.Sprintf(_initTabKey, plat, language); len(s.tabCache[fmt.Sprintf(key)]) == 0 || language == "" {
		if model.IsOverseas(plat) {
			var key = fmt.Sprintf(_initTabKey, plat, _defaultLanguageHant)
			if len(s.tabCache[fmt.Sprintf(key)]) > 0 {
				language = _defaultLanguageHant
			} else {
				language = _defaultLanguageHans
			}
		} else {
			language = _defaultLanguageHans
		}
	}
	var (
		key     = fmt.Sprintf(_initTabKey, plat, language)
		tmptabs = []*show.Tab{}
	)
	res = map[string][]*show.Tab{}
	if tabs, ok := s.tabCache[key]; ok {
	LOOP:
		for _, v := range tabs {
			for _, l := range s.limitsCahce[v.ID] {
				if model.InvalidBuild(build, l.Build, l.Condition) {
					continue LOOP
				}
			}
			if !s.c.ShowHotAll {
				if ab, ok := s.abtestCache[v.Group]; ok {
					if _, ok := s.showTabMids[mid]; !ab.AbTestIn(buvid) && !ok {
						continue LOOP
					}
					a = &abtest.List{}
					a.ListChange(ab)
				}
			}
			tmptabs = append(tmptabs, v)
		}
	}
	if !s.auditTab(mobiApp, build, plat) {
		if menus := s.menus(plat, build); len(menus) > 0 {
			tmptabs = append(tmptabs, menus...)
		}
	}
	for _, v := range tmptabs {
		t := &show.Tab{}
		*t = *v
		t.Pos = len(res[v.ModuleStr]) + 1
		res[v.ModuleStr] = append(res[v.ModuleStr], t)
	}
	if version = s.hash(res); version == ver {
		err = ecode.NotModified
		res = nil
	}
	return
}

func (s *Service) menus(plat int8, build int) (res []*show.Tab) {
	memuCache := s.menuCache
LOOP:
	for _, m := range memuCache {
		if vs, ok := m.Versions[model.PlatAPPBuleChange(plat)]; ok {
			for _, v := range vs {
				if model.InvalidBuild(build, v.Build, v.Condition) {
					continue LOOP
				}
			}
			t := &show.Tab{}
			t.TabMenuChange(m)
			res = append(res, t)
		}
	}
	return
}

func (s *Service) hash(v map[string][]*show.Tab) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return _initVersion
	}
	return strconv.FormatUint(farm.Hash64(bs), 10)
}
