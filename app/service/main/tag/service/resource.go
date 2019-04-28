package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/service/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ResTagMap res tag .
func (s *Service) ResTagMap(c context.Context, arg *model.ArgResTag) (res map[int64]*model.Resource, err error) {
	if res, err = s.dao.ResTagMapCache(c, arg.Oid, arg.Type); err != nil || res != nil {
		return
	}
	if res, err = s.dao.ResTagMap(c, arg.Oid, arg.Type); err != nil {
		return
	}
	rs := make([]*model.Resource, 0, len(res))
	for _, v := range res {
		r := &model.Resource{}
		*r = *v
		rs = append(rs, r)
	}
	s.cacheCh.Save(func() {
		s.dao.AddResTagCache(context.Background(), arg.Oid, arg.Type, rs)
	})
	return
}

// MutiResTagMap get muti res tag.
func (s *Service) MutiResTagMap(c context.Context, arg *model.ArgMutiResTag) (res map[int64][]*model.Resource, err error) {
	var missed []int64
	if res, missed, err = s.dao.ResTagMapCaches(c, arg.Oids, arg.Type); err != nil || len(missed) == 0 {
		return
	}
	missedResTagMap, err := s.dao.ResTagsMap(c, missed, arg.Type)
	if err != nil {
		return
	}
	for oid, missedResTags := range missedResTagMap {
		rs := make([]*model.Resource, 0, len(missedResTags))
		for _, v := range missedResTags {
			r := &model.Resource{}
			*r = *v
			rs = append(rs, r)
		}
		res[oid] = rs
	}
	s.cacheCh.Save(func() {
		s.dao.AddResTagMapCaches(context.Background(), arg.Type, missedResTagMap)
	})
	return
}

// ResTags .
func (s *Service) ResTags(c context.Context, oid int64, typ int32, mid int64) (rts []*model.Resource, err error) {
	var (
		tids   []int64
		tagMap map[int64]*model.Tag
		subMap map[int64]int32
		rtMap  map[int64]*model.Resource
	)
	if rtMap, err = s.resTags(c, oid, typ); err != nil {
		return
	}
	if len(rtMap) == 0 {
		return
	}
	for _, rt := range rtMap {
		tids = append(tids, rt.Tid)
	}
	if tagMap, err = s.tagMap(c, tids); err != nil {
		return
	}
	if len(tagMap) == 0 {
		return
	}
	if mid > 0 {
		subMap, _ = s.isSubTids(c, mid, tids)
	}
	for _, rt := range rtMap {
		var ok bool
		if rt.Tag, ok = tagMap[rt.Tid]; !ok {
			continue
		}
		if mid > 0 {
			if _, ok := subMap[rt.Tid]; ok {
				rt.Tag.Attention = subMap[rt.Tid]
			}
		}
		rts = append(rts, rt)
	}
	sort.Sort(model.Resources(rts))
	return
}

func (s *Service) resTags(c context.Context, oid int64, typ int32) (rtMap map[int64]*model.Resource, err error) {
	var (
		ok             bool
		tids, missTids []int64
	)
	if ok, err = s.dao.ExpireOidCache(c, oid, typ); err != nil {
		return
	}
	if ok {
		if tids, err = s.dao.OidCache(c, oid, typ); err != nil {
			return
		}
		if len(tids) == 0 {
			return
		}
		if len(tids) == 1 {
			if tids[0] == int64(-1) {
				return
			}
		}
		if rtMap, missTids, err = s.dao.ResourceMapCache(c, oid, typ, tids); err != nil {
			return
		}
		if len(missTids) == 0 {
			return
		}
	}
	if rtMap, err = s.dao.ResourceMap(c, oid, typ); err != nil {
		return
	}
	// must deep copy because rtMap data rate.
	newMap := make(map[int64]*model.Resource, len(rtMap))
	for key, value := range rtMap {
		rt := &model.Resource{}
		*rt = *value
		newMap[key] = rt
	}
	s.cacheCh.Save(func() {
		s.dao.AddOidMapCache(context.Background(), oid, typ, newMap)
		s.dao.AddResourceMapCache(context.Background(), newMap)
	})
	return
}

// ResTagLog .
func (s *Service) ResTagLog(c context.Context, oid int64, typ int32, mid int64, pn, ps int, ip string) (atls []*model.ResourceLog, err error) {
	var (
		start = ps * (pn - 1)
	)
	if atls, err = s.dao.ResourceLogs(c, oid, typ, start, ps); err != nil {
		return
	}
	reportMap, err := s.dao.Report(c, oid, typ)
	if err != nil {
		return
	}
	for _, v := range atls {
		k := fmt.Sprintf("%d_%d_%d_%d_%d", v.Oid, v.Type, v.Tid, v.Mid, v.Action)
		if r, ok := reportMap[k]; ok && r != nil {
			if r.State == 1 || r.State == 4 {
				v.State = 1 // 已处理
			} else {
				v.State = 2 // 已举报
			}
		}
	}
	return
}

func (s *Service) addResourceTag(c context.Context, rt *model.Resource, action, logState int32, now time.Time) (err error) {
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	var affected int64
	if affected, err = s.dao.TxAddResource(tx, rt); err != nil || affected == 0 {
		tx.Rollback()
		return
	}
	if affected, err = s.dao.TxUpTagBindCount(tx, rt.Tid, 1); err != nil || affected == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.cache.Save(func() {
		l := &model.ResourceLog{
			Oid:    rt.Oid,
			Type:   rt.Type,
			Mid:    rt.Mid,
			Tid:    rt.Tid,
			Role:   rt.Role,
			Action: action,
			State:  logState,
		}
		t, _ := s.tag(c, rt.Tid)
		s.dao.DelResTagCache(context.Background(), rt.Oid, rt.Type)
		s.dao.AddOidCache(context.Background(), rt.Oid, rt.Type, rt)
		if t != nil && rt.Type == model.ResTypeArchive {
			s.dao.AddResourceLog(context.Background(), t.Name, l)
		}
		s.dao.DelCountCache(context.Background(), rt.Tid)
	})
	return
}

func (s *Service) delResourceTag(c context.Context, rt *model.Resource, action, logState int32, now time.Time) (err error) {
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	var affected int64
	if affected, err = s.dao.TxDelResource(tx, rt); err != nil || affected == 0 {
		tx.Rollback()
		return
	}
	if affected, err = s.dao.TxUpTagBindCount(tx, rt.Tid, -1); err != nil || affected == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.cache.Save(func() {
		l := &model.ResourceLog{
			Oid:    rt.Oid,
			Type:   rt.Type,
			Mid:    rt.Mid,
			Tid:    rt.Tid,
			Role:   rt.Role,
			Action: action,
			State:  logState,
		}
		t, _ := s.tag(c, rt.Tid)
		s.dao.DelResTagCache(context.Background(), rt.Oid, rt.Type)
		s.dao.ZremOidCache(context.Background(), rt.Oid, rt.Tid, rt.Type)
		if t != nil && rt.Type == model.ResTypeArchive {
			s.dao.AddResourceLog(context.Background(), t.Name, l)
		}
		s.dao.DelCountCache(context.Background(), rt.Tid)
	})
	return
}

// PlatformUserBind . role == user|upper
func (s *Service) PlatformUserBind(c context.Context, oid, mid, tid int64, typ, role int32, action int32, ip string) (err error) {
	var (
		rtMap map[int64]*model.Resource
		now   = time.Now()
	)
	rtMap, err = s.resTags(c, oid, typ)
	if err != nil {
		return
	}
	resource, ok := rtMap[tid]
	rt := &model.Resource{
		Oid:  oid,
		Type: typ,
		Tid:  tid,
		Mid:  mid,
		Role: role,
	}
	switch action {
	case model.ResTagLogAdd:
		if ok {
			return ecode.TagAlreadyExist
		}
		if model.ResRoleUpper != role {
			if err = s.checkUserAdd(c, oid, tid, 0, mid, typ, ip); err != nil {
				return
			}
			if _, err = s.checkResTag(c, oid, tid, model.SpamActionAdd, typ); err != nil {
				return
			}
		}
		return s.addResourceTag(c, rt, model.ResTagLogAdd, model.ResTagLogOpen, now)
	case model.ResTagLogDel:
		if !ok {
			return ecode.TagNotExist
		}
		//
		if resource.Attr == 1 { // TODO lock tag
			return ecode.TagArcTagisLocked
		}
		if model.ResRoleUpper != role {
			if err = s.checkUserDel(c, oid, tid, 0, mid, typ, ip); err != nil {
				return
			}
			if _, err = s.checkResTag(c, oid, tid, model.SpamActionDel, typ); err != nil {
				return
			}
		}
		return s.delResourceTag(c, rt, model.ResTagLogDel, model.ResTagLogOpen, now)
	default:
		log.Warn("PlatformUserBind(%d,%d,%d,%d,%d,%d,%s)", oid, mid, tid, typ, role, action, ip)
	}
	return
}

// platformBind .
func (s *Service) platformBind(c context.Context, oid, mid int64, tids []int64, typ, role, logState int32, ip string) (err error) {
	var (
		addTids, delTids []int64
		rtMap            map[int64]*model.Resource
		tidMap           = make(map[int64]bool)
	)
	rtMap, err = s.resTags(c, oid, typ)
	if err != nil {
		return
	}
	for _, tid := range tids {
		tidMap[tid] = true
		if len(rtMap) > 0 {
			if _, ok := rtMap[tid]; ok {
				continue
			}
		}
		addTids = append(addTids, tid)
	}
	for k := range rtMap {
		if _, ok := tidMap[k]; ok {
			continue
		}
		delTids = append(delTids, k)
	}
	now := time.Now()
	for _, tid := range addTids {
		rt := &model.Resource{
			Oid:  oid,
			Type: typ,
			Tid:  tid,
			Mid:  mid,
			Role: role,
		}
		s.addResourceTag(c, rt, model.ResTagLogAdd, logState, now)
	}
	for _, tid := range delTids {
		rt := &model.Resource{
			Oid:  oid,
			Type: typ,
			Tid:  tid,
			Mid:  mid,
			Role: role,
		}
		s.delResourceTag(c, rt, model.ResTagLogDel, logState, now)
	}
	s.cacheCh.Save(func() {
		s.dao.DelResTagCache(context.Background(), oid, typ)
	})
	return
}

func (s *Service) platformDefaultBind(c context.Context, oid, mid int64, tids []int64, tp, role int32, ip string) (err error) {
	var (
		bindTids     []int64
		bindResource []*model.Resource
		delResource  []*model.Resource
		tidMap       = make(map[int64]struct{}, len(tids))
	)
	resMap, err := s.dao.ResourceDefault(c, oid, tp)
	if err != nil {
		return
	}
	for _, tid := range tids {
		tidMap[tid] = struct{}{}
	}
	for _, tid := range tids {
		if _, ok := resMap[tid]; !ok {
			bindTids = append(bindTids, tid)
			continue
		}
		delete(resMap, tid)
	}
	for _, r := range resMap {
		delResource = append(delResource, r)
	}
	if len(delResource) == 0 && len(bindTids) == 0 {
		return
	}
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	for _, tid := range bindTids {
		rt := &model.Resource{
			Oid:   oid,
			Type:  tp,
			Tid:   tid,
			Mid:   mid,
			Role:  role,
			State: model.ResStateDefault,
		}
		bindResource = append(bindResource, rt)
	}
	if len(delResource) > 0 {
		for _, r := range delResource {
			if _, err = s.dao.TxDelResource(tx, r); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	if len(bindResource) > 0 {
		for _, rt := range bindResource {
			if _, err = s.dao.TxAddDefaultResource(tx, rt); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	if err = tx.Commit(); err == nil {
		s.cacheCh.Save(func() {
			s.dao.DelResTagCache(context.Background(), oid, tp)
		})
	}
	return
}

// PlatformUpBind .
func (s *Service) PlatformUpBind(c context.Context, oid, mid int64, tids []int64, typ int32, ip string) (err error) {
	return s.platformBind(c, oid, mid, tids, typ, model.ResRoleUpper, model.ResTagLogOpen, ip)
}

// PlatformAdminBind .
func (s *Service) PlatformAdminBind(c context.Context, oid, mid int64, tids []int64, typ int32, ip string) (err error) {
	return s.platformBind(c, oid, mid, tids, typ, model.ResRoleAdmin, model.ResTagLogClose, ip)
}

// DefaultAdminBind DefaultAdminBind.
func (s *Service) DefaultAdminBind(c context.Context, oid, mid int64, tids []int64, tp int32, ip string) (err error) {
	return s.platformDefaultBind(c, oid, mid, tids, tp, model.ResRoleAdmin, ip)
}

// DefaultUpBind DefaultUpBind.
func (s *Service) DefaultUpBind(c context.Context, oid, mid int64, tids []int64, tp int32, ip string) (err error) {
	return s.platformDefaultBind(c, oid, mid, tids, tp, model.ResRoleUpper, ip)
}

// TODO 待优化
func (s *Service) resource(c context.Context, oid, tid int64, typ int32) (res *model.Resource, err error) {
	if res, err = s.dao.ResourceCache(c, oid, typ, tid); err != nil {
		return
	}
	if res != nil {
		return
	}
	if res, err = s.dao.Resource(c, oid, tid, typ); err != nil {
		return
	}
	if res == nil {
		err = ecode.TagArcTagNotExist
	}
	return
}

// ResOidsByTid .
func (s *Service) ResOidsByTid(c context.Context, tid, limit int64, typ int32, ip string) (res []int64, err error) {
	if limit == 0 || limit > 1000 {
		limit = 1000
	}
	_, res, err = s.dao.ResOidsByTid(c, tid, limit, typ)
	return
}
