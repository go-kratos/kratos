package filter

import (
	"context"
	"math"

	"strings"
	"sync"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	"go-common/app/service/main/filter/service/area"
	"go-common/app/service/main/filter/service/regexp"
	"go-common/library/log"
)

const (
	_filterModeRegexp = 0
	_filterModeString = 1
)

type FilterAreas func(c context.Context, srcMask int64, area string) (ars []*model.FilterAreaInfo, err error)

// Filter filter service def.
type Filter struct {
	// 基础库
	baselibFilters map[int64]*model.Filter // map[sourceMask]filter
	// 业务库
	bizlibFilters map[int64]map[string]*model.Filter // map[sourceMask]map[area]filter
	sync.RWMutex
}

// New create new instance of Filter.
func New() (filter *Filter) {
	return new(Filter)
}

// Load load all filters when init or reload.
func (f *Filter) Load(c context.Context, loader FilterAreas, area *area.Area) (err error) {
	var (
		rules          []*model.FilterAreaInfo
		flt            *model.Filter
		baselibFilters = make(map[int64]*model.Filter)
		bizlibFilters  = make(map[int64]map[string]*model.Filter)
		areaNames      = area.AreaNames()
	)
	// load baselib
	for _, srcMask := range conf.Conf.Property.SourceMask {
		if rules, err = loader(c, srcMask, "common"); err != nil {
			return
		}
		if flt, err = f.createFilter(rules, "common", srcMask); err != nil {
			return
		}
		if flt != nil {
			baselibFilters[srcMask] = flt
		}
	}
	// load bizlib
	for _, srcMask := range conf.Conf.Property.SourceMask {
		bizlibFilters[srcMask] = make(map[string]*model.Filter)
		for _, area := range areaNames {
			if area == "common" {
				continue
			}
			if rules, err = loader(c, srcMask, area); err != nil {
				return
			}
			if flt, err = f.createFilter(rules, area, srcMask); err != nil {
				return
			}
			if flt != nil {
				bizlibFilters[srcMask][area] = flt
			}
		}
	}
	f.Lock()
	f.baselibFilters = baselibFilters
	f.bizlibFilters = bizlibFilters
	f.Unlock()
	return
}

func (f *Filter) createFilter(rules []*model.FilterAreaInfo, area string, source int64) (filter *model.Filter, err error) {
	rl := len(rules)
	log.Info("createFilter area(%s), source(%d), role len(%d)", area, source, rl)
	if rl == 0 {
		return
	}
	filter = &model.Filter{
		Matcher: actriearea.NewMatcher(),
	}
	for _, rule := range rules {
		switch rule.Mode {
		case _filterModeRegexp:
			var re = &model.Regexp{
				TypeIDs: rule.TpIDs,
				Level:   rule.Level(),
				Fid:     rule.ID,
				Area:    area,
			}
			if re.Reg, err = regexp.Compile(rule.Filter); err != nil {
				log.Errorv(context.TODO(), log.KV("err", err.Error()), log.KV("msg", "load filter regexp err"), log.KV("regexp", rule.Filter), log.KV("area", rule.Area))
				continue
			}
			filter.Regs = append(filter.Regs, re)
		case _filterModeString:
			filter.Matcher.Insert(strings.ToLower(rule.Filter), rule.Level(), rule.TpIDs, rule.ID)
		default:
			log.Error("skip load unknown mode (%d) rule (%+v)", rule.Mode, rule)
			continue
		}
	}
	filter.Matcher.Build()
	return
}

func (f *Filter) GetFiltersByArea(area string) (filters []*model.Filter) {
	filters = make([]*model.Filter, 0)
	f.RLock()
	defer f.RUnlock()

	if area == "common" {
		for _, filter := range f.baselibFilters {
			filters = append(filters, filter)
		}
	} else {
		for _, filterMap := range f.bizlibFilters {
			if filter, ok := filterMap[area]; ok {
				filters = append(filters, filter)
			}
		}
	}
	return
}

// GetFilters return hitting filters by param1 & param2.
// param1 means wh baselib
func (f *Filter) GetFilters(area string, baseHit bool) (filters []*model.Filter) {
	f.RLock()
	defer f.RUnlock()
	var (
		bizlibParam  int64 = math.MaxInt64
		baselibParam int64 = math.MaxInt64
	)

	// get baselib filters
	if baseHit {
		for srcMask, filter := range f.baselibFilters {
			if srcMask == 0x00 || srcMask&baselibParam > 0 {
				filters = append(filters, filter)
				continue
			}
		}
	}
	// get bizlib filters
	for srcMask, areaFilters := range f.bizlibFilters {
		if srcMask == 0x00 || srcMask&bizlibParam > 0 {
			if filter, ok := areaFilters[area]; ok {
				filters = append(filters, filter)
			}
		}
	}
	return filters
}
