package ugc

import (
	"context"

	"go-common/app/job/main/tv/model/ugc"
	"go-common/library/log"
)

// ZoneIdx finds out all the passed seasons in DB and then arrange them in a sorted set in Redis
func (s *Service) ZoneIdx() {
	var ctx = context.TODO()
	for catID, rel := range s.ugcTypesRel {
		var (
			firstTid  = rel.TID
			totalTids = []int64{}
		)
		for _, v := range s.arcTypes {
			if v.Pid == firstTid {
				totalTids = append(totalTids, int64(v.ID)) // add second level types
			}
		}
		totalTids = append(totalTids, int64(firstTid)) // add first level type
		IdxRanks, err := s.dao.PassedArcs(ctx, totalTids)
		if err != nil {
			log.Error("UgcZoneIdx - PassedArc TID %d Error %v", firstTid, err)
			continue
		}
		if err = s.appDao.Flush(ctx, int(catID), IdxRanks); err != nil {
			log.Error("UgcZoneIdx - Flush CatID %d Error %v", catID, err)
			continue
		}
	}
}

// loadTids loads the relation between typeIDs and category
func (s *Service) loadTids() {
	var newTids = make(map[int32]int32)
	for catID, rel := range s.ugcTypesRel {
		firstTid := rel.TID
		for _, v := range s.arcTypes {
			if v.Pid == firstTid {
				newTids[v.ID] = catID
			}
		}
		newTids[firstTid] = catID
	}
	if len(newTids) > 0 {
		s.ugcTypesCat = newTids
	}
}

// listMtn maintains the list of zone index
func (s *Service) listMtn(old *ugc.MarkArc, new *ugc.MarkArc) (err error) {
	if old == nil {
		log.Info("ListMtn Old is Nil, NewSn is %v", new)
		old = &ugc.MarkArc{}
	}
	if old.IsPass() && new.IsPass() && old.TypeID == new.TypeID { // no need to take action
		return
	}
	if !old.IsPass() && !new.IsPass() { // no need to take action
		return
	}
	if old.TypeID != 0 { // means old is not null
		if oldCat, ok := s.ugcTypesCat[old.TypeID]; ok { // if old one is in our list, remove it firstly
			if err = s.appDao.ZRemIdx(ctx, int(oldCat), old.AID); err != nil {
				log.Error("listMtn - ZRemIdx - Category: %d, Arc: %d, Error: %v", oldCat, old.AID, err)
				return
			}
		}
	}
	catID, ok := s.ugcTypesCat[new.TypeID]
	if !ok {
		log.Warn("TypeID %d Is Not our target, ignore", new.TypeID)
		return
	}
	if new.IsPass() { // passed now
		if err = s.appDao.ZAddIdx(ctx, int(catID), new.Ctime, new.AID); err != nil {
			log.Error("listMtn - ZAddIdx - Category: %d, Arc: %d, Error: %v", catID, new.AID, err)
			return
		}
		log.Info("Add Aid %d Into Zone %d", new.AID, catID)
	}
	return
}
