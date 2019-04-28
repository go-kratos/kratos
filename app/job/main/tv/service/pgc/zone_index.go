package pgc

import (
	"context"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

const (
	_seasonPassed = 1
	_epPassed     = 3
	_cmsValid     = 1
	_notDeleted   = 0
)

// ZoneIdx finds out all the passed seasons in DB and then arrange them in a sorted set in Redis
func (s *Service) ZoneIdx() {
	var (
		_pgcZones = s.c.Cfg.PGCZonesID
		ctx       = context.Background()
	)
	for _, v := range _pgcZones {
		zoneSns, err := s.dao.PassedSn(ctx, v)
		if err != nil {
			log.Error("ZoneIdx - PassedSn %d Error %v", v, err)
			continue
		}
		if err = s.dao.Flush(ctx, v, zoneSns); err != nil {
			log.Error("ZoneIdx - Flush %d Error %v", v, err)
			continue
		}
	}
}

// listMtn maintains the list of zone index
func (s *Service) listMtn(oldSn *model.MediaSn, newSn *model.MediaSn) (err error) {
	if oldSn == nil {
		log.Info("ListMtn OldSn is Nil, NewSn is %v", newSn)
		oldSn = &model.MediaSn{}
	}
	if oldSn.Check == _seasonPassed && oldSn.IsDeleted == _notDeleted && oldSn.Valid == _cmsValid { // previously passed
		if !(newSn.Check == _seasonPassed && newSn.IsDeleted == _notDeleted && newSn.Valid == _cmsValid) { // not passed now
			if err = s.dao.ZRemIdx(ctx, newSn.Category, newSn.ID); err != nil {
				log.Error("listMtn - ZRemIdx - Category: %d, Sn: %s, Error: %v", newSn.Category, newSn, err)
				return
			}
			log.Info("Remove Sid %d From Zone %d", newSn.ID, newSn.Category)
		}
	} else { // previously not passed, or not exist
		if newSn.Check == _seasonPassed && newSn.IsDeleted == _notDeleted && newSn.Valid == _cmsValid { // passed now
			if err = s.dao.ZAddIdx(ctx, newSn.Category, newSn.Ctime, newSn.ID); err != nil {
				log.Error("listMtn - ZAddIdx - Category: %d, Sn: %s, Error: %v", newSn.Category, newSn, err)
				return
			}
			log.Info("Add Sid %d Into Zone %d", newSn.ID, newSn.Category)
		}
	}
	return
}
