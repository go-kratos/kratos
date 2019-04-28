package service

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_homepageID   = 0
	_jpType       = 1
	recomCategory = 1
)

// load homepage data
func (s *Service) loadHome(ctx context.Context) (err error) {
	if len(s.ZoneData) == 0 { // pick area's data
		log.Error("[loadHP] Can't Pick Zone Top Data!!")
		return
	}
	homepage := model.Homepage{}
	if homepage.Recom, err = s.HomeRecom(ctx); err != nil {
		log.Error("[loadHP] Can't Load Home Recom Data !")
		return
	} // recom part
	s.buildHeaderSids(homepage.Recom)
	homepage.Lists, homepage.Latest = s.HomeList() // list part
	s.HomeData = &homepage
	return
}

// HomeFollow picks the homepage with follow data
func (s *Service) HomeFollow(ctx context.Context, req *model.ReqHomeFollow) (data *model.Homepage, err error) {
	if s.HomeData == nil {
		log.Error("HomeFollow Data is Nil ")
		err = ecode.ServiceUnavailable
		return
	}
	data = &(*s.HomeData)
	if req.AccessKey != "" {
		data.Follow = s.FollowData(ctx, req.AccessKey)
	}
	var newRecom []*model.Card // old version, filter UGC data
	for _, v := range data.Recom {
		if !v.IsUGC() {
			newRecom = append(newRecom, v)
		}
	}
	data.Recom = newRecom
	return
}

func (s *Service) buildHeaderSids(Recom []*model.Card) {
	if len(Recom) == 0 {
		log.Error("Homepage Recom is Empty!")
		return
	}
	newMap := make(map[int]int)
	for _, v := range Recom {
		newMap[v.SeasonID] = 1
	}
	if len(newMap) > 0 {
		s.HeaderSids = newMap
	}
}

// pgcIndexs picks and treats PGC Index Data
func (s *Service) pgcIndexs(ctx context.Context) (indexCards []*model.Card) {
	var moduleOrder = s.conf.Cfg.ZonesInfo.ZonesName
	pgcData, err := s.dao.HeaderData(ctx, s.TVAppInfo) // pick data from PGC API
	if err != nil {
		log.Error("[loadHP] Can't Pick PGC/AI Data, Err: %v", err)
		return
	}
	for _, v := range moduleOrder { // arrange data according to hard-code order
		if value, ok := pgcData[v]; ok {
			for _, card := range value {
				if card.NewEP != nil {
					card.Type = _typePGC // define card type
					indexCards = append(indexCards, card)
				}
			}
		} else {
			log.Error("PGC Data Missing %s data", v)
		}
	}
	if err := s.cardIntervSn(indexCards); err != nil { // replace cover & title by CMS data
		log.Error("[cardIntervSn] ERROR [%v]", err)
	}
	return
}

// HomeRecom for the Homepage header, merge the rank data and the intervention data and gets the final header data
func (s *Service) HomeRecom(ctx context.Context) (homeRecom []*model.Card, err error) {
	var (
		hsize     = s.ZonesInfo[_homepageID].Top
		intervReq = &model.ReqZoneInterv{
			RankType: _homepageID,
			Category: recomCategory,
			Limit:    hsize,
		}
		interv []*model.Card
	)
	indexData := s.pgcIndexs(ctx)                  // 1. Treat PGC data
	resp, err := s.dao.ZoneIntervs(ctx, intervReq) // 2. Treat Interventions
	if err != nil {
		log.Error("[LoadPages] Can't Pick Intervention Data, Err: %v", err)
		return
	}
	if interv, err = s.intervToCards(ctx, resp); err != nil {
		log.Error("[LoadPages] Can't Combine Intervention Data, Err: %v", err)
		return
	}
	homeRecom = mergeSlice(interv, indexData)
	homeRecom = duplicate(homeRecom) // remove duplicated data
	homeRecom = cutSlice(homeRecom, hsize)
	return
}

// hideIndexShow for the configured zones to hide the index show, we modify their list
func (s *Service) hideIndexShow(list map[string][]*model.Card) {
	if len(s.conf.Homepage.HideIndexShow) == 0 {
		return
	}
	for _, zone := range s.conf.Homepage.HideIndexShow {
		if _, ok := list[zone]; !ok {
			continue
		} else {
			for k, v := range list[zone] { // hide index show without the influence to zone page
				var newCard = *v
				newCard.NewEP = &model.NewEP{
					ID:        v.NewEP.ID,
					Index:     v.NewEP.Index,
					Cover:     v.NewEP.Cover,
					IndexShow: "",
				}
				list[zone][k] = &newCard
			}
			log.Info("Hide Index Show for Zone %s", zone)
		}
	}
}

// HomeList gets the five zones' list for the homepage
func (s *Service) HomeList() (list map[string][]*model.Card, latest []*model.Card) {
	list = make(map[string][]*model.Card)
	conf := s.ZonesInfo[_homepageID]
	var (
		idList   = map[int][]*model.Card{} // the zone list, key is ID
		listSize = conf.Bottom             // the homepage list size
	)
	if len(s.ZoneData) == 0 {
		log.Error("ZoneData is Empty!")
		return
	}
	// remove the items alraedy in the headers
	for k, v := range s.ZoneData {
		// jp resources fill the latest part
		if k == _jpType {
			latest, list[s.ZonesInfo[k].Name] = s.HomeJP(v)
			continue
		}
		idList[k] = []*model.Card{}
		for _, vcard := range v {
			if len(idList[k]) >= listSize { // we only pick 10 data for the lists
				break
			}
			if _, ok := s.HeaderSids[vcard.SeasonID]; !ok {
				idList[k] = append(idList[k], vcard)
			}
		}
		list[s.ZonesInfo[k].Name] = idList[k]
	}
	// hide index show for configured list
	s.hideIndexShow(list)
	return
}

// HomeJP picks the JP resources to fill the latest part and the list part
func (s *Service) HomeJP(zoneData []*model.Card) (latest []*model.Card, list []*model.Card) {
	var (
		Intervs    []*model.Card // homepage middle interventions
		err        error
		middleSize = s.ZonesInfo[_homepageID].Middle
		listSize   = s.ZonesInfo[_homepageID].Bottom
		middleM    = s.ZonesInfo[_homepageID].MiddleM
	)
	// get homepage latest's intervention
	if middleM != 0 {
		if Intervs, err = s.modPGCIntervs(ctx, middleM, middleSize); err != nil {
			return
		}
	} else {
		if Intervs, err = s.getIntervs(context.TODO(), _homepageID, _latest, middleSize); err != nil {
			return
		}
	}
	// remove duplicated
	allCards := mergeSlice(Intervs, zoneData)
	allCards = duplicate(allCards)
	latestSource := []*model.Card{}
	// pick enough data for middle and list by checking duplication with the header data
	for _, vcard := range allCards {
		if len(latestSource) >= middleSize+listSize { // we pick enough data for the middle and the list of jp
			break
		}
		if _, ok := s.HeaderSids[vcard.SeasonID]; !ok { // remove duplicated with the header data
			latestSource = append(latestSource, vcard)
		}
	}
	// cut the latest part and add them into headerSids for duplication check
	latest = cutSlice(latestSource, middleSize)
	for _, v := range latest {
		s.HeaderSids[v.SeasonID] = 1
	}
	// cut JP list part
	if len(latest) < middleSize {
		list = []*model.Card{}
		return
	}
	list = cutSlice(latestSource[middleSize:], listSize)
	return
}
