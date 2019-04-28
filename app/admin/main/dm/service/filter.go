package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
)

// UpFilters return up filters
func (s *Service) UpFilters(c context.Context, mid, ftype, pn, ps int64) (res []*model.UpFilter, total int64, err error) {
	//type all
	if ftype == int64(model.FilterTypeAll) {
		if res, total, err = s.dao.UpFiltersAll(c, mid, pn, ps); err != nil {
			log.Error("s.dao.UpFiltersAll(mid:%d) error(%v)", mid, err)
		}
		return
	}
	if res, total, err = s.dao.UpFilters(c, mid, ftype, pn, ps); err != nil {
		log.Error("s.dao.UpFilters(mid:%d, type:%d) error(%v)", mid, ftype, err)
	}
	return
}

// EditUpFilters edit up filters.
func (s *Service) EditUpFilters(c context.Context, id, mid int64, fType, active int8) (affect int64, err error) {
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
	if affect, err = s.dao.UpdateUpFilter(tx, mid, id, active); err != nil {
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
	if _, err = s.dao.UpdateUpFilterCnt(tx, mid, fType, int(affect), limit+1); err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.dao.DelUpFilterCache(ctx, mid, 0)
	})

	return
}
