package service

import (
	"context"
	"encoding/json"
	"math"
	"sync/atomic"
	"time"

	"go-common/app/interface/main/dm/model"
	dm2Mdl "go-common/app/interface/main/dm2/model"
	"go-common/library/log"

	"github.com/zhenjl/cityhash"
)

var (
	_defaultVersion   = uint64(0)
	globalRuleVersion = uint64(time.Now().Nanosecond())
)

// AddUserRule add user rule
func (s *Service) AddUserRule(c context.Context, fType int8, mid int64, filters []string, comment string) (res []*dm2Mdl.UserFilter, err error) {
	arg := &dm2Mdl.ArgAddUserFilters{
		Type:    fType,
		Mid:     mid,
		Filters: filters,
		Comment: comment,
	}
	if res, err = s.dmRPC.AddUserFilters(c, arg); err != nil {
		log.Error("dmRPC.AddUserFilters(%+v) error(%v)", arg, err)
	}
	return
}

// UserRules return user rules
func (s *Service) UserRules(c context.Context, mid int64) (res []*dm2Mdl.UserFilter, err error) {
	arg := &dm2Mdl.ArgMid{Mid: mid}
	if res, err = s.dmRPC.UserFilters(c, arg); err != nil {
		log.Error("dmRPC.UserFilters(%+v) error(%v)", arg, err)
	}
	return
}

// DelUserRules delete user rules
func (s *Service) DelUserRules(c context.Context, mid int64, idss []int64) (affect int64, err error) {
	arg := &dm2Mdl.ArgDelUserFilters{Mid: mid, IDs: idss}
	if affect, err = s.dmRPC.DelUserFilters(c, arg); err != nil {
		log.Error("dmRPC.DelUserFilters(%+v) error(%v)", arg, err)
	}
	return
}

// GlobalRuleVersion return global rule version
func (s *Service) GlobalRuleVersion() uint64 {
	return globalRuleVersion
}

// AddGlobalRule add global rule
func (s *Service) AddGlobalRule(c context.Context, fType int8, filter string) (res *dm2Mdl.GlobalFilter, err error) {
	arg := &dm2Mdl.ArgAddGlobalFilter{Type: fType, Filter: filter}
	if res, err = s.dmRPC.AddGlobalFilter(c, arg); err != nil {
		log.Error("dmRPC.AddGlobalFilter(%+v) error(%v)", arg, err)
		return
	}
	atomic.StoreUint64(&globalRuleVersion, _defaultVersion)
	return
}

// GlobalRules return global rules
func (s *Service) GlobalRules(c context.Context) (res []*dm2Mdl.GlobalFilter, err error) {
	arg := &dm2Mdl.ArgGlobalFilters{}
	if res, err = s.dmRPC.GlobalFilters(c, arg); err != nil {
		log.Error("dmRPC.GlobalFilters(%+v) error(%v)", arg, err)
		return
	}
	if len(res) == 0 {
		atomic.StoreUint64(&globalRuleVersion, _defaultVersion)
	} else {
		var buf []byte
		if buf, err = json.Marshal(res); err != nil {
			log.Error("json.Marshal() error(%v)", err)
			return
		}
		atomic.StoreUint64(&globalRuleVersion, cityhash.CityHash64(buf, 16)%math.MaxInt64)
	}
	return
}

// DelGlobalRules delete global rules
func (s *Service) DelGlobalRules(c context.Context, ids []int64) (affect int64, err error) {
	arg := &dm2Mdl.ArgDelGlobalFilters{IDs: ids}
	if affect, err = s.dmRPC.DelGlobalFilters(c, arg); err != nil {
		log.Error("dmRPC.DelGlobalFilters(%+v) error(%v)", arg, err)
		return
	}
	// update global rule version
	atomic.StoreUint64(&globalRuleVersion, _defaultVersion)
	return
}

// FilterList get user filter list
func (s *Service) FilterList(c context.Context, mid, cid int64) (l *model.UserFilterList, err error) {
	l = new(model.UserFilterList)
	arg := &dm2Mdl.ArgUpFilters{Mid: mid}
	res, err := s.dmRPC.UpFilters(c, arg)
	if err != nil {
		log.Error("dmRPC.UpFilters(%v) error(%v)", arg, err)
		return
	}
	if len(res) == 0 {
		return
	}
	for _, f := range res {
		switch f.Type {
		case dm2Mdl.FilterTypeBottom:
			l.Bottom = f.Active
			continue
		case dm2Mdl.FilterTypeTop:
			l.Top = f.Active
			continue
		case dm2Mdl.FilterTypeRev:
			l.Reverse = f.Active
			continue
		}
		filter := &model.IndexFilter{
			ID:       f.ID,
			MID:      f.Mid,
			Filter:   f.Filter,
			Activate: f.Active,
			Regex:    f.Type,
			Ctime:    int64(f.Ctime),
		}
		l.Filter = append(l.Filter, filter)
	}
	return
}

// EditFilter edit up filter from creative center.
func (s *Service) EditFilter(c context.Context, cid, mid int64, filter string, fType, state int8) (err error) {
	if state == dm2Mdl.FilterActive {
		arg := &dm2Mdl.ArgAddUpFilters{
			Mid:     mid,
			Type:    fType,
			Filters: []string{filter},
		}
		if err = s.dmRPC.AddUpFilters(c, arg); err != nil {
			log.Error("dmRPC.AddUpFilters(%v) error(%v)", arg, err)
		}
	} else {
		arg := &dm2Mdl.ArgEditUpFilters{
			Type:    fType,
			Mid:     mid,
			Active:  dm2Mdl.FilterUnActive,
			Filters: []string{filter},
		}
		if _, err = s.dmRPC.EditUpFilters(c, arg); err != nil {
			log.Error("dmRPC.EditUpFilters(%v) error(%v)", arg, err)
		}
	}
	return
}
