package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"go-common/app/interface/main/dm2/lib/xregex"
	"go-common/app/interface/main/dm2/model"
	arcMdl "go-common/app/service/main/archive/model/archive"
	assmdl "go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/zhenjl/cityhash"
)

var (
	_defaultVersion     = uint64(0)
	globalFilterVersion = uint64(time.Now().Nanosecond())
)

// AddUserFilters multi add user filter,fltMap is a map struct,key:filter content,value is comment.
func (s *Service) AddUserFilters(c context.Context, mid int64, fType int8, filters map[string]string) (res []*model.UserFilter, err error) {
	// copy map,because delete operation will be used in function
	fltMap := make(map[string]string)
	for k, v := range filters {
		fltMap[k] = v
	}
	res = make([]*model.UserFilter, 0)
	for filter := range fltMap {
		if fType == model.FilterTypeRegex {
			reg := xregex.New()
			if len([]rune(filter)) > model.FilterLenRegex {
				err = ecode.DMFilterTooLong
				return
			}
			if _, err = reg.Parse(filter); err != nil {
				log.Error("regex parse(filter:%+v) error(%v)", filter, err)
				err = ecode.DMFitlerIllegalRegex
				return
			}
		}
		if fType == model.FilterTypeText {
			if len([]rune(filter)) > model.FilterLenText {
				err = ecode.DMFilterTooLong
				return
			}
		}
	}
	data, err := s.dao.UserFilter(c, mid, fType)
	if err != nil {
		return
	}
	for filter := range fltMap {
		for _, v := range data {
			if fType == model.FilterTypeText && strings.ToLower(filter) == strings.ToLower(v.Filter) {
				res = append(res, v)
				delete(fltMap, filter) // delete repeat filter in filter map
				break
			} else if filter == v.Filter {
				res = append(res, v)
				delete(fltMap, filter) // delete repeat filter in filter map
				break
			}
		}
	}
	if len(fltMap) == 0 {
		return
	}
	tx, err := s.dao.BeginBiliDMTrans(c)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	count, err := s.dao.UserFilterCnt(c, tx, mid, fType)
	if err != nil {
		return
	}
	var (
		num   = len(fltMap)
		id    int64
		limit int
	)
	switch fType {
	case model.FilterTypeText:
		if count+num > model.FilterMaxUserText {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUserText
	case model.FilterTypeRegex:
		if count+num > model.FilterMaxUserReg {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUserReg
	case model.FilterTypeID:
		if count+num > model.FilterMaxUserID {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUserID
	}
	// the filter id in database will be used in delete method,so add user filter in for loop
	for filter, comment := range fltMap {
		id, err = s.dao.AddUserFilter(tx, mid, fType, filter, comment)
		if err != nil {
			return
		}
		res = append(res, &model.UserFilter{ID: id, Mid: mid, Type: fType, Filter: filter, Comment: comment})
	}
	if count == model.FilterNotExist { // not exist, insert into table count=num
		if _, err = s.dao.InsertUserFilterCnt(c, tx, mid, fType, num); err != nil {
			return
		}
	} else { // already exist, set count=count+1
		if _, err = s.dao.UpdateUserFilterCnt(c, tx, mid, fType, 1, int64(limit)); err != nil {
			return
		}
	}
	// synchronized delete cache
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelUserFilterCache(ctx, mid)
	})
	return
}

// UserFilters return user filters
func (s *Service) UserFilters(c context.Context, mid int64) (res []*model.UserFilter, err error) {
	if res, err = s.dao.UserFilterCache(c, mid); err != nil {
		err = nil // NOTE load from db if cache error
	} else if len(res) > 0 {
		return
	}
	if res, err = s.dao.UserFilters(c, mid); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.AddUserFilterCache(ctx, mid, res)
	})
	return
}

// DelUserFilters delete user filters
func (s *Service) DelUserFilters(c context.Context, mid int64, idss []int64) (affect int64, err error) {
	var (
		idMap      = make(map[int8][]int64)
		aft, limit int64
	)
	res, err := s.dao.UserFiltersByID(c, mid, idss)
	if err != nil {
		return
	}
	for _, f := range res {
		idMap[f.Type] = append(idMap[f.Type], f.ID)
	}
	tx, err := s.dao.BeginBiliDMTrans(c)
	if err != nil {
		log.Error("tx.BeginBiliDMTrans() error(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	for fType, ids := range idMap {
		if aft, err = s.dao.DelUserFilter(tx, mid, ids); err != nil {
			return
		}
		switch fType {
		case model.FilterTypeText:
			limit = model.FilterMaxUserText
		case model.FilterTypeRegex:
			limit = model.FilterMaxUserReg
		case model.FilterTypeID:
			limit = model.FilterMaxUserID
		}
		if _, err = s.dao.UpdateUserFilterCnt(c, tx, mid, fType, -aft, limit+1); err != nil {
			return
		}
		affect = affect + aft
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelUserFilterCache(ctx, mid)
	})
	return
}

// AddUpFilters add up rule.fltMap is a map struct,key:filter content,value is comment.
func (s *Service) AddUpFilters(c context.Context, mid int64, fType int8, filters map[string]string) (err error) {
	fltMap := make(map[string]string)
	// copy map,because delete operation will be used in function
	for k, v := range filters {
		fltMap[k] = v
	}
	for filter := range fltMap {
		if fType == model.FilterTypeRegex {
			if len([]rune(filter)) > model.FilterLenRegex {
				err = ecode.DMFilterTooLong
				return
			}
			reg := xregex.New()
			if _, err = reg.Parse(filter); err != nil {
				log.Error("filter(%s) parse error(%v)", filter, err)
				err = ecode.DMFitlerIllegalRegex
				return
			}
		}
		if fType == model.FilterTypeText {
			if len([]rune(filter)) > model.FilterLenText {
				err = ecode.DMFilterTooLong
				return
			}
		}
	}
	res, err := s.dao.UpFilter(c, mid, fType)
	if err != nil {
		return
	}
	hash := model.Hash(mid, 0)
	for filter := range fltMap {
		if fType == model.FilterTypeID && filter == hash { // 忽略拉黑自己
			delete(fltMap, filter)
		}
		for _, f := range res {
			if fType == model.FilterTypeText && strings.ToLower(filter) == strings.ToLower(f.Filter) {
				delete(fltMap, filter)
				break
			} else if filter == f.Filter {
				delete(fltMap, filter)
				break
			}
		}
	}
	var (
		limit int
		num   = len(fltMap)
	)
	if num == 0 {
		return
	}
	tx, err := s.dao.BeginBiliDMTrans(c)
	if err != nil {
		log.Error("tx.BeginBiliDMTrans() error(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	count, err := s.dao.UpFilterCnt(c, tx, mid, fType)
	if err != nil {
		return
	}
	switch fType {
	case model.FilterTypeText:
		if count+num > model.FilterMaxUpText {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUpText
	case model.FilterTypeRegex:
		if count+num > model.FilterMaxUpReg {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUpReg
	case model.FilterTypeID:
		if count+num > model.FilterMaxUpID {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUpID
	}
	if _, err = s.dao.MultiAddUpFilter(tx, mid, fType, fltMap); err != nil {
		return
	}
	if count == model.FilterNotExist { // not exist, insert
		if _, err = s.dao.InsertUpFilterCnt(c, tx, mid, fType, num); err != nil {
			return
		}
	} else { // already exist, set count=count+fNum
		if _, err = s.dao.UpdateUpFilterCnt(c, tx, mid, fType, num, limit); err != nil {
			return
		}
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelUpFilterCache(ctx, mid)
	})
	return
}

// AddUpFilterID block user by upper or assist.fltMap is a map struct,key:user hashid,value is dm msg.
func (s *Service) AddUpFilterID(c context.Context, mid, oid int64, fltMap map[string]string) (err error) {
	var (
		isAssist bool
		sub      *model.Subject
		fType    = model.FilterTypeID
	)
	if sub, err = s.subject(c, model.SubTypeVideo, oid); err != nil {
		return
	}
	if !s.isUpper(sub.Mid, mid) {
		if err = s.isAssist(c, sub.Mid, mid); err != nil {
			return
		}
	}
	res, err := s.dao.UpFilter(c, sub.Mid, fType)
	if err != nil {
		return
	}
	hash := model.Hash(mid, 0)
	for filter := range fltMap {
		if filter == hash { //忽略拉黑自己
			delete(fltMap, filter)
		}
		for _, f := range res {
			if filter == f.Filter {
				delete(fltMap, filter)
				break
			}
		}
	}
	var (
		limit int
		num   = len(fltMap)
	)
	if num == 0 {
		return
	}
	tx, err := s.dao.BeginBiliDMTrans(c)
	if err != nil {
		log.Error("tx.BeginBiliDMTrans() error(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	count, err := s.dao.UpFilterCnt(c, tx, sub.Mid, fType)
	if err != nil {
		return
	}
	switch fType {
	case model.FilterTypeText:
		if count+num > model.FilterMaxUpText {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUpText
	case model.FilterTypeRegex:
		if count+num > model.FilterMaxUpReg {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUpReg
	case model.FilterTypeID:
		if count+num > model.FilterMaxUpID {
			err = ecode.DMFilterOverMax
			return
		}
		limit = model.FilterMaxUpID
	}
	if _, err = s.dao.MultiAddUpFilter(tx, sub.Mid, fType, fltMap); err != nil {
		return
	}
	if count == model.FilterNotExist { // not exist, insert
		if _, err = s.dao.InsertUpFilterCnt(c, tx, sub.Mid, fType, num); err != nil {
			return
		}
	} else { // already exist, set count=count+fNum
		if _, err = s.dao.UpdateUpFilterCnt(c, tx, sub.Mid, fType, num, limit); err != nil {
			return
		}
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelUpFilterCache(ctx, sub.Mid)
	})
	if isAssist {
		for filter, comment := range fltMap {
			if len(comment) > 50 {
				comment = fmt.Sprintf("%s...", comment[:50])
			}
			arg := &assmdl.ArgAssistLogAdd{
				Mid:       sub.Mid,
				AssistMid: mid,
				Type:      assmdl.TypeDm,
				Action:    assmdl.ActDisUser,
				SubjectID: sub.Mid,
				ObjectID:  filter,
				Detail:    comment,
			}
			select {
			case s.assistLogChan <- arg:
			default:
				log.Error("assistLogChan is full")
			}
		}
	}
	return
}

// UpFilters return up filters
func (s *Service) UpFilters(c context.Context, mid int64) (res []*model.UpFilter, err error) {
	if res, err = s.dao.UpFilterCache(c, mid); err != nil {
		err = nil // load from db if cache error
	} else if len(res) > 0 {
		return
	}
	if res, err = s.dao.UpFilters(c, mid); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.AddUpFilterCache(ctx, mid, res)
	})
	return
}

// BanUsers ban user by upper or assist.
func (s *Service) BanUsers(c context.Context, mid, oid int64, dmids []int64) (err error) {
	var (
		isAssist bool
		fltMap   = make(map[string]string)
	)
	sub, err := s.subject(c, model.SubTypeVideo, oid)
	if err != nil {
		return
	}
	if !s.isUpper(sub.Mid, mid) {
		if err = s.isAssist(c, sub.Mid, mid); err != nil {
			return
		}
	}
	idxMap, _, err := s.dao.IndexsByid(c, sub.Type, oid, dmids)
	if err != nil || len(idxMap) == 0 {
		return
	}
	ctsmap, err := s.dao.Contents(c, oid, dmids)
	if err != nil {
		return
	}
	for _, idx := range idxMap {
		hashID := model.Hash(idx.Mid, 0)
		if _, ok := fltMap[hashID]; !ok {
			var comment string
			if v, ok := ctsmap[idx.ID]; ok {
				comment = v.Msg
			}
			fltMap[hashID] = comment
		}
	}
	if len(fltMap) == 0 {
		return
	}
	if err = s.AddUpFilters(c, sub.Mid, model.FilterTypeID, fltMap); err != nil {
		return
	}
	if isAssist {
		for filter, comment := range fltMap {
			if len(comment) > 50 {
				comment = fmt.Sprintf("%s...", comment[:50])
			}
			arg := &assmdl.ArgAssistLogAdd{
				Mid:       sub.Mid,
				AssistMid: mid,
				Type:      assmdl.TypeDm,
				Action:    assmdl.ActDisUser,
				SubjectID: sub.Mid,
				ObjectID:  filter,
				Detail:    comment,
			}
			select {
			case s.assistLogChan <- arg:
			default:
				log.Error("assistLogChan is full")
			}
		}
	}
	return
}

// CancelBanUsers cancel up filter by assist.
func (s *Service) CancelBanUsers(c context.Context, mid, aid int64, filters []string) (err error) {
	var (
		isAssist bool
		arg      = arcMdl.ArgAid2{Aid: aid}
	)
	res, err := s.arcRPC.Archive3(c, &arg)
	if err != nil {
		log.Error("s.arcRPC.Archive3(%v) error(%v)", arg, err)
		return
	}
	if !s.isUpper(res.Author.Mid, mid) {
		if err = s.isAssist(c, res.Author.Mid, mid); err != nil {
			return
		}
	}
	if _, err = s.EditUpFilters(c, res.Author.Mid, model.FilterTypeID, model.FilterUnActive, filters); err != nil {
		return
	}
	if isAssist {
		for _, filter := range filters {
			arg := &assmdl.ArgAssistLogAdd{
				Mid:       res.Author.Mid,
				AssistMid: mid,
				Type:      assmdl.TypeDm,
				Action:    assmdl.ActCancelDisUser,
				SubjectID: res.Author.Mid,
				ObjectID:  filter,
				Detail:    "cancel ban users",
			}
			select {
			case s.assistLogChan <- arg:
			default:
				log.Error("assistLogChan is full")
			}
		}
	}
	return
}

// EditUpFilters edit up filters.
func (s *Service) EditUpFilters(c context.Context, mid int64, fType, active int8, filters []string) (affect int64, err error) {
	var limit int
	tx, err := s.dao.BeginBiliDMTrans(c)
	if err != nil {
		log.Error("tx.BeginBiliDMTrans error(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if affect, err = s.dao.UpdateUpFilter(tx, mid, fType, active, filters); err != nil {
		return
	}
	switch fType {
	case model.FilterTypeText:
		limit = model.FilterMaxUpText
	case model.FilterTypeRegex:
		limit = model.FilterMaxUpReg
	case model.FilterTypeID:
		limit = model.FilterMaxUpID
	}
	if active == model.FilterUnActive {
		affect = -affect
		limit = 10000
	}
	if _, err = s.dao.UpdateUpFilterCnt(c, tx, mid, fType, int(affect), limit+1); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelUpFilterCache(ctx, mid)
	})
	return
}

// GlobalFilterVersion return global filter version
func (s *Service) GlobalFilterVersion() uint64 {
	return globalFilterVersion
}

// AddGlobalFilter add global filter
func (s *Service) AddGlobalFilter(c context.Context, fType int8, filter string) (ret *model.GlobalFilter, err error) {
	ret = &model.GlobalFilter{}
	if fType != model.FilterTypeText && fType != model.FilterTypeRegex {
		err = ecode.DMFilterIllegalType
		return
	}
	// regex varify
	if fType == model.FilterTypeRegex {
		if len([]rune(filter)) > model.FilterLenRegex {
			err = ecode.DMFilterTooLong
			return
		}
		reg := xregex.New()
		if _, err = reg.Parse(filter); err != nil {
			log.Error("filter(%s) parse error(%v)", filter, err)
			err = ecode.DMFitlerIllegalRegex
			return
		}
	}
	if fType == model.FilterTypeText {
		if len([]rune(filter)) > model.FilterLenText {
			err = ecode.DMFilterTooLong
			return
		}
	}
	res, err := s.dao.GlobalFilter(c, fType, filter)
	if err != nil {
		return
	}
	for _, f := range res {
		switch fType {
		case model.FilterTypeText:
			if strings.ToLower(filter) == strings.ToLower(f.Filter) {
				ret = f
				err = ecode.DMFilterExist
				return
			}
		case model.FilterTypeRegex, model.FilterTypeID:
			if filter == f.Filter {
				ret = f
				err = ecode.DMFilterExist
				return
			}
		}
	}
	if ret.ID, err = s.dao.AddGlobalFilter(c, fType, filter); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelGlobalFilterCache(ctx)
	})
	// update global filter version
	atomic.StoreUint64(&globalFilterVersion, _defaultVersion)
	return
}

// GlobalFilters return global filters
func (s *Service) GlobalFilters(c context.Context) (res []*model.GlobalFilter, err error) {
	var (
		done bool
		sid  int64
		data []byte
		f    []*model.GlobalFilter
	)
	if res, err = s.dao.GlobalFilterCache(c); err != nil {
		err = nil
	}
	if res != nil {
		return
	}
	for !done && (len(res) < 50000) {
		if f, err = s.dao.GlobalFilters(c, sid, 1000); err != nil {
			return
		}
		if len(f) == 0 {
			break
		}
		if len(f) < 1000 {
			done = true
		}
		sid = f[len(f)-1].ID + 1
		res = append(res, f...)
	}
	if len(res) > 50000 {
		res = res[:50000]
	}
	// set empty cache
	if len(res) == 0 {
		res = []*model.GlobalFilter{}
		atomic.StoreUint64(&globalFilterVersion, _defaultVersion)
	} else {
		// get global filter version
		if data, err = json.Marshal(res); err != nil {
			log.Error("json.Marshal() error(%v)", err)
			return
		}
		atomic.StoreUint64(&globalFilterVersion, cityhash.CityHash64(data, 16)%math.MaxInt64)
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.AddGlobalFilterCache(ctx, res)
	})
	return
}

// DelGlobalFilters delete global filters
func (s *Service) DelGlobalFilters(c context.Context, ids []int64) (affect int64, err error) {
	if affect, err = s.dao.DelGlobalFilters(c, ids); err != nil {
		return
	}
	if err = s.dao.DelGlobalFilterCache(c); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelGlobalFilterCache(ctx)
	})
	// update global filter version
	atomic.StoreUint64(&globalFilterVersion, _defaultVersion)
	return
}
