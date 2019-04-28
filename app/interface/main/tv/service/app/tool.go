package service

import (
	"strconv"
	"time"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// mergeSlice merges two slices and return a new slice
func mergeSlice(s1 []*model.Card, s2 []*model.Card) []*model.Card {
	slice := make([]*model.Card, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

// mergeSliceM merges two slices and return a new slice
func mergeSliceM(s1 []*model.Module, s2 []*model.Module) []*model.Module {
	slice := make([]*model.Module, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

// transform string to int
func atoi(number string) (result int) {
	result, _ = strconv.Atoi(number)
	return result
}

// dupliInt is used to merge the base and the interv, then remove duplicated numbers
func dupliInt(base []int64, interv []int64) (merged []int64) {
	var dupMap = make(map[int64]int)
	for _, v := range append(interv, base...) {
		length := len(dupMap)
		dupMap[v] = 1
		if len(dupMap) != length {
			merged = append(merged, v)
		}
	}
	return
}

// removeInt is used to remove interv data from base, and return the rest base data
func removeInt(base []int64, interv []int64) (restBase []int64) {
	var intervMap = make(map[int64]int)
	for _, v := range interv {
		intervMap[v] = 1
	}
	for _, v := range base {
		if _, ok := intervMap[v]; !ok {
			restBase = append(restBase, v)
		}
	}
	return
}

// applyInterv is used for ES index page intervention
func applyInterv(base []int64, interv []int64, pn int) (res []int64) {
	if len(interv) == 0 {
		return base
	}
	if pn == 1 { // first page, put the intervs in the beginning
		return dupliInt(base, interv)
	}
	return removeInt(base, interv) // other pages, remove the intervs from the ES result
}

// remove duplicated Cards
func duplicate(a []*model.Card) (ret []*model.Card) {
	resUGC := make(map[int]int)
	resPGC := make(map[int]int)
	for _, v := range a {
		if v.IsUGC() {
			if _, ok := resUGC[v.SeasonID]; ok {
				log.Warn("[UGC] v.SeasonID %d is duplicated", v.SeasonID)
				continue
			}
			resUGC[v.SeasonID] = 1
			ret = append(ret, v)
		} else {
			if _, ok := resPGC[v.SeasonID]; ok {
				log.Warn("v.SeasonID %d is duplicated", v.SeasonID)
				continue
			}
			resPGC[v.SeasonID] = 1
			ret = append(ret, v)
		}
	}
	return
}

// cut the slice with the limit
func cutSlice(source []*model.Card, lengthLimit int) (res []*model.Card) {
	if len(source) <= lengthLimit {
		return source
	}
	return source[0:lengthLimit]
}

// transform card data to mod_card data
func cardTransform(source []*model.Card) (target []*model.ModCard) {
	for _, v := range source {
		target = append(target, &model.ModCard{
			Card: *v,
		})
	}
	return
}

// followToMod transforms follow structure data to module structure data
func followToMod(source []*model.Follow) (target []*model.ModCard) {
	for _, v := range source {
		target = append(target, &model.ModCard{
			Card: model.Card{
				SeasonID: atoi(v.SeasonID),
				Title:    v.Title,
				Cover:    v.Cover,
				NewEP: &model.NewEP{
					ID:        int64(atoi(v.NewEP.EpisodeID)),
					Index:     v.NewEP.Index,
					IndexShow: v.NewEP.IndexTitle,
					Cover:     v.NewEP.Cover,
				},
				CornerMark: v.CornerMark,
			},
			LastEPIndex:   v.UserSeason.LastEPIndex,
			NewestEPIndex: v.NewestEPIndex,
			TotalCount:    v.TotalCount,
			IsFinish:      v.IsFinish,
		})
	}
	return
}

// elapsed record the function's execution time
func elapsed(funcName string) func() {
	start := time.Now()
	return func() {
		log.Info("[Elapsed] %s took %v\n", funcName, time.Since(start))
	}
}

// transformCards rewrite, use season & ep's cache to build the cards, instead of pick them from DB
func (s *Service) transformCards(sids []int64) (target []*model.Card, targetMap map[int]*model.Card) {
	var (
		seasons     []*model.SeasonCMS
		eps         map[int64]*model.EpCMS
		newestEPIDs []int64
		err         error
	)
	targetMap = make(map[int]*model.Card)
	if seasons, newestEPIDs, err = s.cmsDao.LoadSnsCMS(ctx, sids); err != nil {
		log.Error("transformCards - LoadSnsCMS - Sids %v, Err %v", sids, err)
		return
	}
	if eps, err = s.cmsDao.LoadEpsCMS(ctx, newestEPIDs); err != nil {
		log.Error("transformCard - LoadEpsCMS - Epids %v, Err %v", newestEPIDs, err)
		return
	}
	for _, val := range seasons {
		newcard := &model.Card{
			SeasonID: int(val.SeasonID),
			Title:    val.Title,
			Cover:    val.Cover,
			Type:     _typePGC,
			NewEP: &model.NewEP{
				ID: val.NewestEPID,
			},
		}
		if val.NeedVip() { // card add vip corner mark
			newcard.CornerMark = &(*s.conf.Cfg.SnVipCorner)
		}
		if val.NewestEPID == 0 {
			log.Error("transformCard - NewestEPID of SeasonID %v is Empty", val.SeasonID)
		} else {
			// epMeta info
			if epval, ok := eps[val.NewestEPID]; ok {
				newcard.NewEP.Index = epval.Title
				newcard.NewEP.Cover = epval.Cover
			} else {
				log.Error("transformCard - EpMetas of Epid is Empty", val.NewestEPID)
			}
			// epIndex show
			if len(s.PGCIndexShow) > 0 { // maybe first launch it's empty, wait 2 minutes it will be ready
				if indexShow, ok := s.PGCIndexShow[int64(val.SeasonID)]; ok && len(indexShow) > 0 {
					newcard.NewEP.IndexShow = indexShow
				} else {
					log.Error("transformCard - Missing Index_show For Sid:%v", val.SeasonID)
				}
			}
		}
		target = append(target, newcard)
		targetMap[newcard.SeasonID] = newcard
	}
	return
}
