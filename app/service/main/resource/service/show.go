package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

// loadCardCache load all card cache
func (s *Service) loadCardCache() {
	now := time.Now()
	hdm, err := s.show.PosRecs(context.TODO(), now)
	if err != nil {
		log.Error("s.show.PosRecs error(%v)", err)
		return
	}
	itm, aids, err := s.show.RecContents(context.TODO(), now)
	if err != nil {
		log.Error("s.show.RecContents error(%v)", err)
		return
	}
	tmpItem := map[int]map[int64]*model.ShowItem{}
	for recid, aid := range aids {
		tmpItem[recid] = s.fromCardAids(context.TODO(), aid)
	}
	tmp := s.mergeCard(context.TODO(), hdm, itm, tmpItem, now)
	s.cardCache = tmp
}

func (s *Service) loadSideBarCache() {
	now := time.Now()
	sidebar, limits, err := s.show.SideBar(context.TODO(), now)
	if err != nil {
		log.Error("s.show.SideBar error(%v)", err)
		return
	}
	s.sideBarCache = sidebar
	s.sideBarLimitCache = limits
}

// SideBars get side bars
func (s *Service) SideBars(c context.Context) (res *model.SideBars) {
	res = &model.SideBars{
		SideBar: s.sideBarCache,
		Limit:   s.sideBarLimitCache,
	}
	return res
}

// RegionCard get voice card.
func (s *Service) RegionCard(c context.Context, plat int8, build int) (res *model.Head, err error) {
	res = &model.Head{}
	sw := s.cardCache[plat]
	if sw == nil {
		return
	}
	if model.InvalidBuild(build, sw.Build, sw.Condition) {
		return
	}
	*res = *sw
	res.FillBuildURI(plat, build)
	return
}

// fromCardAids get Aids.
func (s *Service) fromCardAids(c context.Context, aids []int64) (data map[int64]*model.ShowItem) {
	var (
		arc *api.Arc
		ok  bool
	)
	as, err := s.arcRPC.Archives3(c, &archive.ArgAids2{Aids: aids})
	if err != nil {
		log.Error("s.arcRPC.Archives3 error(%v)", err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arcRPC.Archives3(%v) length is 0", aids)
		return
	}
	data = map[int64]*model.ShowItem{}
	for _, aid := range aids {
		if arc, ok = as[aid]; ok {
			if !arc.IsNormal() {
				continue
			}
			i := &model.ShowItem{}
			i.FromArchivePB(arc)
			data[aid] = i
		}
	}
	return
}

// mergeCard merge Card
func (s *Service) mergeCard(c context.Context, hdm map[int8][]*model.Card, itm map[int][]*model.Content, tmpItems map[int]map[int64]*model.ShowItem, now time.Time) (res map[int8]*model.Head) {
	res = map[int8]*model.Head{}
	for plat, hds := range hdm {
		for _, hd := range hds {
			var (
				sis []*model.ShowItem
			)
			its, ok := itm[hd.ID]
			if !ok {
				its = []*model.Content{}
			}
			tmpItem, ok := tmpItems[hd.ID]
			if !ok {
				tmpItem = map[int64]*model.ShowItem{}
			}
			switch hd.Type {
			case 1:
				for _, ci := range its {
					si := s.fillCardItem(ci, tmpItem)
					if si.Title != "" {
						sis = append(sis, si)
					}
				}
			default:
				continue
			}
			if len(sis) == 0 {
				continue
			}
			sw := &model.Head{
				CardID:    hd.ID,
				Title:     hd.Title,
				Type:      hd.TypeStr,
				Build:     hd.Build,
				Condition: hd.Condition,
				Plat:      hd.Plat,
			}
			if hd.Cover != "" {
				sw.Cover = hd.Cover
			}
			switch sw.Type {
			case model.GotoDaily:
				sw.Date = now.Unix()
				sw.Param = hd.Rvalue
				sw.URI = hd.URI
				sw.Goto = hd.Goto
			}
			sw.Body = sis
			res[plat] = sw
		}
	}
	return
}

// fillCardItem fill card
func (s *Service) fillCardItem(csi *model.Content, tsi map[int64]*model.ShowItem) (si *model.ShowItem) {
	si = &model.ShowItem{}
	switch csi.Type {
	case model.CardGotoAv:
		si.Goto = model.GotoAv
		si.Param = csi.Value
	}
	si.URI = model.FillURI(si.Goto, si.Param)
	if si.Goto == model.GotoAv {
		aid, err := strconv.ParseInt(si.Param, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", si.Param, err)
		} else {
			if it, ok := tsi[aid]; ok {
				si = it
				if csi.Title != "" {
					si.Title = csi.Title
				}
			} else {
				si = &model.ShowItem{}
			}
		}
	}
	return
}

// Audit all audit config.
func (s *Service) Audit(c context.Context) map[string][]int {
	return s.auditCache
}
