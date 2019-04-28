package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/tv/model"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_homeFocus       = 1
	_fiveFocus       = 2
	_sixFocus        = 3
	_verticalOneList = 4
	_verticalTwoList = 5
	_horizontalList  = 6
	_followMod       = 7
)

func (s *Service) loadMods(ctx context.Context) (err error) {
	var (
		pageMods []*model.Module
		newMap   = make(map[int][]*model.Module)
		zoneCfg  = s.conf.Cfg.ZonesInfo
		pages    = append([]int{}, zoneCfg.PGCZonesID...)
	)
	for _, v := range zoneCfg.UGCZonesID {
		pages = append(pages, int(v))
	}
	// load home & other module pages
	if newMap[_homepageID], err = s.ModHome(); err != nil {
		return
	}
	for _, v := range s.RegionInfo {
		pages = append(pages, v.PageID)
	}
	pages = remRep(pages)
	log.Info("loadModsRegion Len Pages %d, RegionInfo %d", len(pages), len(s.RegionInfo))
	for _, v := range pages {
		if pageMods, err = s.pageData(ctx, v); err != nil {
			log.Error("LoadModPage Data PID %d, Err %v", v, err)
			return
		}
		newMap[v] = pageMods
	}
	if len(newMap) > 0 {
		s.ModPages = newMap
	}
	return
}

// ModHome load modularized homepage
func (s *Service) ModHome() (homepage []*model.Module, err error) {
	var (
		home     = []*model.Module{}
		pid      = _homepageID
		pageMods []*model.Module
	)
	if pageMods, err = s.pageData(ctx, pid); err != nil {
		log.Error("LoadModPage Data PID %d, Err %v", pid, err)
		return
	}
	// pick old logic homepage recom
	if len(s.HomeData.Recom) != 0 {
		home = append(home, &model.Module{
			Type:   _homeFocus,
			PageID: pid,
			Data:   cardTransform(s.HomeData.Recom),
		})
	}
	homepage = mergeSliceM(home, pageMods) // add home mods data into homepage
	return
}

// PageFollow serves the http level, it picks the follow location and fill in with follow data and then output
func (s *Service) PageFollow(c context.Context, req *model.ReqPageFollow) (res []*model.Module, err error) {
	var (
		ok       bool
		resMods  []*model.Module
		cfgBuild = s.conf.Cfg.PGCFilterBuild
	)
	if resMods, ok = s.ModPages[req.PageID]; !ok {
		err = ecode.ServiceUnavailable
		log.Error("ModPage %d, Err %v", req.PageID, err)
		return
	}
	// not logged in AND no need to filter ugc
	if req.AccessKey == "" && req.Build > cfgBuild {
		res = resMods
		return
	}
	//decorate follow data + build filter logic
	for _, v := range resMods {
		modV := *v
		if v.Type == _followMod && req.AccessKey != "" {
			modV.Data = followToMod(s.FollowData(c, req.AccessKey))
			res = append(res, &modV)
			continue
		}
		if req.Build <= cfgBuild { // if old version, only output pgc data
			if v.IsUGC() {
				continue
			}
			if len(v.Data) > 0 {
				var newData = []*model.ModCard{}
				for _, tp := range v.Data {
					if tp.IsUGC() {
						continue
					}
					newData = append(newData, tp)
				}
				log.Info("ModID %d, Data Length %d, After UGC Filter Length %d", v.ID, len(v.Data), len(newData))
				modV.Data = newData
			}
		}
		res = append(res, &modV)
	}
	return
}

// dupRecom removes the recom sids from the ugc or pgc data
func dupRecom(recomSids map[int]int, CardList map[int][]*model.Card) {
	if len(recomSids) > 0 { // remove recom data from PGC Data
		for srcType, srcList := range CardList {
			var filterS []*model.Card
			for _, source := range srcList {
				if _, ok := recomSids[source.SeasonID]; !ok {
					filterS = append(filterS, source)
				}
			}
			CardList[srcType] = filterS
		}
	}
}

// modMapping is used to map new zone idx to old zone idx pages
func (s *Service) modMapping(mod *model.Module) {
	var (
		tp  *arcwar.Tp
		cfg = s.conf.Cfg.ZonesInfo
		pid int32
		err error
	)
	mod.MoreTreat()                                                   // simple mapping with the page ID
	if !mod.OnHomepage() || (mod.OnHomepage() && !mod.JumpNewIdx()) { // if not homepage or it's homepage but it jumps to old idx and zone page
		return
	}
	if !mod.IsUGC() { // pgc jump to the source's old index page
		mod.MorePage = mod.Source
		return
	}
	if tp, err = s.arcDao.TypeInfo(int32(mod.Source)); err != nil { // if source can't be found, jump to the default Idx ID
		mod.MorePage = cfg.OldIdxJump
		return
	}
	if tp.Pid == 0 { // as it's UGC, let's find it's first level
		pid = tp.ID
	} else {
		pid = tp.Pid
	}
	if oldIdxCat, ok := cfg.OldIdxMapping[fmt.Sprintf("%d", pid)]; ok { // if we can mapping, we map to the old index page
		mod.MorePage = oldIdxCat
	} else { // otherwise we jump to the default index page
		mod.MorePage = cfg.OldIdxJump
	}
}

// modData loads the data of a module
func (s *Service) pageData(ctx context.Context, pageID int) (res []*model.Module, err error) {
	var (
		mods      []*model.Module
		respMod   *model.RespModInterv
		recomSids = make(map[int]int)
		recomAids = make(map[int]int)
		PGCListM  = copyCardList(s.PGCOrigins)
		UGCListM  = copyCardList(s.UGCOrigins)
	)
	if mods, err = s.dao.ModPage(ctx, pageID); err != nil {
		log.Error("ModPage ID %d, Error %v", pageID, err)
		return
	}
	for _, v := range mods {
		s.modMapping(v) // treat the new index page mapping logic
	}
	// if there is no mod, no need to continue the logic
	if len(mods) == 0 {
		log.Error("ModPage ID %d, Modules Empty", pageID)
		return
	}
	// remove the page's intervention and the home recom from the source data
	for _, mod := range mods { // get all the mods' intervention
		if respMod, err = s.dao.ModIntervs(ctx, mod.ID, mod.Capacity); err != nil || respMod == nil {
			log.Error("modPGCIntervs ModID %d, Capacity %d, Err %v", mod.ID, mod.Capacity, err)
			continue
		}
		if len(respMod.SIDs) > 0 {
			for _, sid := range respMod.SIDs {
				recomSids[int(sid)] = 1
			}
		}
		if len(respMod.AIDs) > 0 {
			for _, avid := range respMod.AIDs {
				recomAids[int(avid)] = 1
			}
		}
	}
	if pageID == _homepageID { // if it's the homepage, we need also to filter the home recom
		for _, focus := range s.HomeData.Recom {
			if focus.IsUGC() {
				recomAids[focus.SeasonID] = 1
			} else {
				recomSids[focus.SeasonID] = 1
			}
		}
	}
	dupRecom(recomSids, PGCListM) // filter pgc
	dupRecom(recomAids, UGCListM) // filter ugc
	// build each module's data
	for _, mod := range mods {
		switch mod.Type {
		case _sixFocus, _fiveFocus, _verticalOneList, _verticalTwoList, _horizontalList:
			if err = s.modData(ctx, &model.ReqModData{
				Mod:      mod,
				PGCListM: PGCListM,
				UGCListM: UGCListM,
			}); err != nil {
				log.Error("Load PageID %d, Mod %v, Data error %v", pageID, mod, err)
				return
			}
			res = append(res, mod)
		case _followMod:
			res = append(res, mod)
		default: // ignore invalid module
			log.Error("Invalid ModID %d, Type %d", mod.ID, mod.Type)
		}
	}
	return
}

// focusData loads fiveFocus or sixFocus's data
func (s *Service) modData(ctx context.Context, req *model.ReqModData) (err error) {
	mod := req.Mod
	var (
		intervs     []*model.Card
		capacity    int
		pid         = mod.PageID
		backupCards []*model.Card
		ok          bool
	)
	// capacity logic
	if mod.Type == _fiveFocus {
		capacity = 5
	} else if mod.Type == _sixFocus {
		capacity = 6
	} else { // lists
		capacity = mod.Capacity
	}
	// get interventions
	if intervs, err = s.modIntervs(ctx, mod.ID, capacity); err != nil {
		return
	}
	mod.Data = cardTransform(intervs)
	// pick up backup data
	if mod.IsUGC() {
		if backupCards, ok = req.UGCListM[mod.Source]; !ok {
			log.Error("UGCListM Page %d Source %d is Empty!", pid, mod.Source)
		}
	} else {
		if backupCards, ok = req.PGCListM[mod.Source]; !ok {
			log.Error("PGCListM Page %d Source %d is Empty!", pid, mod.Source)
			return
		}
	}
	// merge intervention with pgc source data
	allCards := mergeSlice(intervs, backupCards)
	allCards = duplicate(allCards)
	mod.Data = cardTransform(cutSlice(allCards, capacity))
	// remove used cards from the source data
	backupCards = allCards[len(mod.Data):]
	if mod.IsUGC() {
		req.UGCListM[mod.Source] = backupCards
		usedCards := allCards[0:len(mod.Data)]
		s.parentDup(usedCards, req)
	} else {
		req.PGCListM[mod.Source] = backupCards
	}
	return
}

// parentDup is dedicated for second level ugc types, pick the used cards
// and find its father list, and remove the used card from the father list also
func (s *Service) parentDup(usedCards []*model.Card, req *model.ReqModData) {
	if len(usedCards) == 0 { // if none of cards has been used
		return
	}
	var (
		secondTid = int32(req.Mod.Source)
		err       error
		tinfo     *arcwar.Tp
		children  []*arcwar.Tp
	)
	if tinfo, err = s.arcDao.TypeInfo(secondTid); err != nil {
		log.Error("parentDup TypeInfo Tid %d, Err %v", secondTid, err)
		return
	}
	usedCardMap := make(map[int]int) // usedCards build map
	for _, v := range usedCards {
		usedCardMap[v.SeasonID] = 1
	}
	if tinfo.Pid == 0 { // if it's first level
		if children, err = s.arcDao.TypeChildren(tinfo.ID); err != nil {
			log.Error("parentDup Pid %d Cant found children", tinfo.ID)
			return
		}
		for _, child := range children {
			childList, ex := req.UGCListM[int(child.ID)]
			if !ex {
				log.Warn("parentDup Pid %d ChildID %d is not in UGClistM", tinfo.ID, child.ID)
				continue
			}
			req.UGCListM[int(child.ID)] = dupList(childList, usedCardMap)
		}
		return
	}
	parentSrc, exist := req.UGCListM[int(tinfo.Pid)] // if it's second level type
	if !exist || len(parentSrc) == 0 {               // if parent list doesn't exist, just return
		return
	}
	req.UGCListM[int(tinfo.Pid)] = dupList(parentSrc, usedCardMap)
}

// dupList travels the list, remove the targetIDs cards
func dupList(list []*model.Card, targetIDs map[int]int) (newlist []*model.Card) {
	for _, vv := range list {
		if _, dup := targetIDs[vv.SeasonID]; !dup {
			newlist = append(newlist, vv)
		}
	}
	return
}

// copy the pgc data for each page, to guarantee they are all independent for each page
func copyCardList(source map[int][]*model.Card) (copied map[int][]*model.Card) {
	copied = make(map[int][]*model.Card)
	for k, v := range source {
		copied[k] = []*model.Card{}
		for _, vcard := range v {
			copied[k] = append(copied[k], &(*vcard))
		}
	}
	return
}

// modPGCIntervs gets the interventions of a module, only treat PGC data
func (s *Service) modPGCIntervs(c context.Context, modID int, nbLimit int) (modCards []*model.Card, err error) {
	var (
		respInterv *model.RespModInterv
	)
	if respInterv, err = s.dao.ModIntervs(c, modID, nbLimit); err != nil || respInterv == nil {
		log.Error("[modPGCIntervs] ModID: %d, Limit %d, Can't Pick Intervention Data, Err: %v", modID, nbLimit, err)
		return
	}
	if len(respInterv.Ranks) == 0 {
		log.Warn("[LoadPages] Mod %d, NbLimit %d, Intervention Empty", modID, nbLimit)
		return // empty result
	}
	modCards, _ = s.transformCards(respInterv.SIDs)
	return
}

// modIntervs gets the interventions of a module, treat both PGC & UGC Data
func (s *Service) modIntervs(c context.Context, modID int, nbLimit int) (modCards []*model.Card, err error) {
	var (
		resp *model.RespModInterv
	)
	if resp, err = s.dao.ModIntervs(c, modID, nbLimit); err != nil || resp == nil {
		log.Error("[modIntervs] ModID: %d, Limit %d, Can't Pick Intervention Data, Err: %v", modID, nbLimit, err)
		return
	}
	if len(resp.Ranks) == 0 {
		log.Warn("[LoadPages] Mod %d, NbLimit %d, Intervention Empty", modID, nbLimit)
		return // empty result
	}
	return s.intervToCards(c, resp)
}

func (s *Service) intervToCards(ctx context.Context, resp *model.RespModInterv) (modCards []*model.Card, err error) {
	if len(resp.Ranks) == 0 {
		return
	}
	pgcCards, pgcCardsMap := s.transformCards(resp.SIDs) // transform PGC
	if len(pgcCards) == len(resp.Ranks) {                // if all the ranks are pgc type
		return pgcCards, err
	}
	var ugcCardsMap map[int64]*model.ArcCMS
	if ugcCardsMap, err = s.cmsDao.LoadArcsMediaMap(ctx, resp.AIDs); err != nil { // transform UGC
		log.Error("[modIntervs] Can't Pick MediaCache Data, Aids: %v, Err: %v", resp.AIDs, err)
		return
	}
	for _, v := range resp.Ranks {
		if v.IsUGC() {
			if arc, ok := ugcCardsMap[v.ContID]; ok {
				modCards = append(modCards, arc.ToCard())
			} else {
				log.Warn("modIntervs, ContID:%d, ContType:%d, Not found", v.ContID, v.ContType)
			}
		} else {
			if sn, ok := pgcCardsMap[int(v.ContID)]; ok {
				modCards = append(modCards, sn)
			} else {
				log.Warn("modIntervs, ContID:%d, ContType:%d, Not found", v.ContID, v.ContType)
			}
		}
	}
	return
}

// remRep delete repeat .
func remRep(slc []int) (result []int) {
	tempMap := map[int]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}
