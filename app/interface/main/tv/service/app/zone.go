package service

import (
	"context"
	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

const (
	_rank    = 1
	_list    = 2
	_latest  = 3 // added No.3 type of module
	_typePGC = 1
)

// preparePGCData loads
func (s *Service) preparePGCData(c context.Context) {
	var (
		zones      = s.conf.Cfg.ZonesInfo.PGCZonesID
		originData []*model.Card
		err        error
	)
	for _, v := range zones {
		if originData, err = s.pgcData(c, v); err != nil {
			log.Error("LoadPGCList Data for Zone %d, Err %v", v, err)
			continue
		}
		if len(originData) == 0 {
			log.Error("LoadPGCList Espically Data for Zone %d, Empty", v)
			continue
		}
		s.PGCOrigins[v] = originData
	}
}

// load all the zone's data and save them in memroy
func (s *Service) zonesData(c context.Context) (err error) {
	zones := s.conf.Cfg.ZonesInfo.PGCZonesID
	for _, v := range zones {
		if err = s.loadZone(c, v); err != nil {
			log.Error("[ZoneData] Error (%v)", err)
			return
		}
	}
	return
}

// load zone data
func (s *Service) loadZone(c context.Context, sType int) (err error) {
	var (
		top, middle, bottom []*model.Card
		conf                = s.ZonesInfo[sType]
		zoneData            = []*model.Card{}
		pgcList             = copyCardList(s.PGCOrigins)
	)
	s.ZoneSids[sType] = make(map[int]int) //re-init the zone header sids for removing duplication
	// re-init the rank data's map
	s.RankData[sType] = make(map[string][]*model.Card)
	// top module logic
	if conf.Top != 0 {
		reqTop := &model.ReqZone{
			SType:       sType,
			IntervType:  _rank,
			LengthLimit: conf.Top,
			IntervM:     conf.TopM,
			PGCListM:    pgcList,
		}
		if top, err = s.zoneLogic(c, reqTop); err != nil {
			log.Error("ZoneTop %d Error %v", sType, err)
			return
		}
		zoneData = mergeSlice(zoneData, top)
		s.RankData[sType]["rank"] = top
	}
	// middle module logic
	if conf.Middle != 0 {
		reqMiddle := &model.ReqZone{
			SType:       sType,
			IntervType:  _latest,
			LengthLimit: conf.Middle,
			IntervM:     conf.MiddleM,
			PGCListM:    pgcList,
		}
		if middle, err = s.zoneLogic(c, reqMiddle); err != nil {
			log.Error("ZoneMiddle %d Error %v", sType, err)
			return
		}
		zoneData = mergeSlice(zoneData, middle)
		s.RankData[sType]["latest"] = middle
	}
	// bottom module logic
	if conf.Bottom != 0 {
		reqBottom := &model.ReqZone{
			SType:       sType,
			IntervType:  _list,
			LengthLimit: conf.Bottom,
			IntervM:     0,
			PGCListM:    pgcList,
		}
		if bottom, err = s.zoneLogic(c, reqBottom); err != nil {
			log.Error("ZoneBottom %d Error %v", sType, err)
		}
		zoneData = mergeSlice(zoneData, bottom)
		s.RankData[sType]["list"] = bottom
	}
	s.ZoneData[sType] = zoneData
	return
}

// pgcData gets the origin pgc data and intervene with TV CMS data
func (s *Service) pgcData(c context.Context, seasonType int) (intervened []*model.Card, err error) {
	intervened, err = s.dao.ChannelData(c, seasonType, s.TVAppInfo)
	if err != nil {
		log.Error("[LoadPGCList] Can't Pick PGC/AI Data, Zone %d, Err: %v", seasonType, err)
		return
	}
	if err2 := s.cardIntervSn(intervened); err2 != nil {
		log.Error("[cardIntervSn] ERROR [%v]", err2)
	}
	return
}

// getIntervs gets the specified type of intervention ( top, middle or botton ) with the number limit and transform the intervention to cards
func (s *Service) getIntervs(c context.Context, sType int, intervType int, nbLimit int) (resCards []*model.Card, err error) {
	var resp *model.RespModInterv
	if resp, err = s.dao.ZoneIntervs(c, &model.ReqZoneInterv{ // for home & zone old logic, we only pick PGC data
		RankType: sType,
		Category: intervType,
		Limit:    nbLimit,
	}); err != nil {
		log.Error("[loadZone] Can't Pick Intervention Data, Err: %v", err)
		return
	}
	if len(resp.SIDs) == 0 {
		log.Warn("[LoadPages] Zone %d, Category %d, Intervention Empty", sType, intervType)
		return // empty result
	}
	resCards, _ = s.transformCards(resp.SIDs)
	return
}

// zoneTop returns the zop module data
func (s *Service) zoneLogic(c context.Context, req *model.ReqZone) (modData []*model.Card, err error) {
	var (
		Intervs []*model.Card
		pgcData = s.PGCOrigins[req.SType]
		lenLmt  = req.LengthLimit
		sType   = req.SType
		intervM = req.IntervM
	)
	if intervM != 0 {
		if Intervs, err = s.modPGCIntervs(c, intervM, lenLmt); err != nil {
			return
		}
	} else {
		if Intervs, err = s.getIntervs(c, sType, req.IntervType, lenLmt); err != nil {
			return
		}
	}
	// merge top intervention and pgc data, then remove duplicated
	allCards := mergeSlice(Intervs, pgcData)
	allCards = duplicate(allCards)
	for _, v := range allCards {
		if len(modData) >= lenLmt {
			break
		}
		if _, ok := s.ZoneSids[sType][v.SeasonID]; !ok {
			modData = append(modData, v)
			s.ZoneSids[sType][v.SeasonID] = 1
		}
	}
	// remove the used pgc cards
	req.PGCListM[sType] = allCards[len(modData):]
	return
}
