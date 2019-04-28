package feed

import (
	"context"
	"strconv"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/bplus"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/model"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// Actives2 return actives
func (s *Service) Actives2(c context.Context, id, mid int64, mobiApp string, plat int8, build, forceHost int, now time.Time) (items []card.Handler, cover string, isBnj bool, bnjDays int, err error) {
	if id == s.c.Bnj.TabID {
		isBnj = true
		nt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		bt, _ := time.Parse("2006-01-02", s.c.Bnj.BeginTime)
		bnjDays = int(bt.Sub(nt).Hours() / 24)
		if bnjDays < 0 {
			bnjDays = 0
		}
	}
	rs := s.tabCache[id]
	if items, err = s.dealTab2(c, rs, mid, mobiApp, plat, build, forceHost, now); err != nil {
		log.Error("s.dealTab(%v) error(%v)", rs, err)
		return
	}
	cover = s.coverCache[id]
	return
}

func (s *Service) dealTab2(c context.Context, rs []*operate.Active, mid int64, mobiApp string, plat int8, build, forceHost int, now time.Time) (is []card.Handler, err error) {
	if len(rs) == 0 {
		is = []card.Handler{}
		return
	}
	var (
		paids, aids, tids, roomIDs, sids, metaIDs, picIDs []int64
		pam, am                                           map[int64]*archive.ArchiveWithPlayer
		rm                                                map[int64]*live.Room
		sm                                                map[int64]*bangumi.Season
		metam                                             map[int64]*article.Meta
		tagm                                              map[int64]*tag.Tag
		picm                                              map[int64]*bplus.Picture
	)
	convergem := map[int64]*operate.Card{}
	specialm := map[int64]*operate.Card{}
	downloadm := map[int64]*operate.Download{}
	for _, r := range rs {
		switch r.Type {
		case model.GotoPlayer:
			if r.Pid != 0 {
				paids = append(paids, r.Pid)
			}
		case model.GotoPlayerLive:
			if r.Pid != 0 {
				roomIDs = append(roomIDs, r.Pid)
			}
		case model.GotoTabTagRcmd:
			if r.Pid != 0 {
				var taids []int64
				if taids, err = s.rcmd.TagTop(c, mid, r.Pid, r.Limit); err != nil {
					log.Error("%+v", err)
					err = nil
					continue
				}
				tids = append(tids, r.Pid)
				r.Items = make([]*operate.Active, 0, len(taids))
				for _, aid := range taids {
					item := &operate.Active{Pid: aid, Goto: model.GotoAv, Param: strconv.FormatInt(aid, 10)}
					r.Items = append(r.Items, item)
					aids = append(aids, aid)
				}
			}
		case model.GotoConverge:
			cardm, aid, roomID, metaID := s.convergeCard(c, 3, r.Pid)
			for id, card := range cardm {
				convergem[id] = card
			}
			aids = append(aids, aid...)
			roomIDs = append(roomIDs, roomID...)
			metaIDs = append(metaIDs, metaID...)
		case model.GotoTabContentRcmd:
			for _, item := range r.Items {
				if item.Pid == 0 {
					continue
				}
				switch item.Goto {
				case cdm.GotoAv:
					aids = append(aids, item.Pid)
				case cdm.GotoLive:
					roomIDs = append(roomIDs, item.Pid)
				case cdm.GotoBangumi:
					sids = append(sids, item.Pid)
				case cdm.GotoGame:
					if card, ok := s.downloadCache[item.Pid]; ok {
						downloadm[item.Pid] = card
					}
				case cdm.GotoArticle:
					metaIDs = append(metaIDs, item.Pid)
				case cdm.GotoSpecial:
					cardm := s.specialCard(c, item.Pid)
					for id, card := range cardm {
						specialm[id] = card
					}
				case cdm.GotoPicture:
					// 版本过滤5.37为新卡片
					if (plat == model.PlatIPhone && build > 8300) || (plat == model.PlatAndroid && build > 5365000) {
						picIDs = append(picIDs, item.Pid)
					}
				}
			}
		case model.GotoSpecial:
			cardm := s.specialCard(c, r.Pid)
			for id, card := range cardm {
				specialm[id] = card
			}
		}
	}
	g, ctx := errgroup.WithContext(c)
	if len(tids) != 0 {
		g.Go(func() (err error) {
			if tagm, err = s.tg.InfoByIDs(c, 0, tids); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(aids) != 0 {
		g.Go(func() (err error) {
			if am, err = s.ArchivesWithPlayer(ctx, aids, 0, "", 0, 0, 0, 0); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(paids) != 0 {
		g.Go(func() (err error) {
			if pam, err = s.ArchivesWithPlayer(ctx, paids, 32, mobiApp, 0, 0, forceHost, build); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(roomIDs) != 0 {
		g.Go(func() (err error) {
			if rm, err = s.lv.AppMRoom(ctx, roomIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(sids) != 0 {
		g.Go(func() (err error) {
			if sm, err = s.bgm.Seasons(ctx, sids, now); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(metaIDs) != 0 {
		g.Go(func() (err error) {
			if metam, err = s.art.Articles(ctx, metaIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(picIDs) != 0 {
		g.Go(func() (err error) {
			if picm, err = s.bplus.DynamicDetail(ctx, picIDs...); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	is = make([]card.Handler, 0, len(rs))
	for _, r := range rs {
		var main interface{}
		cardGoto := cdm.CardGt(r.Type)
		op := &operate.Card{}
		op.From(cardGoto, r.Pid, 0, plat, build)
		// 版本过滤
		hasThreePoint := (plat == model.PlatIPhone && build >= 8240) || (plat == model.PlatAndroid && build > 5341000)
		if hasThreePoint {
			op.FromSwitch(cdm.SwitchFeedIndexTabThreePoint)
		}
		h := card.Handle(plat, cardGoto, "", cdm.ColumnSvrDouble, nil, tagm, nil, nil, nil)
		if h == nil {
			continue
		}
		switch r.Type {
		case model.GotoPlayer:
			main = pam
		case model.GotoPlayerLive:
			main = rm
		case model.GotoSpecial:
			op = specialm[r.Pid]
		case model.GotoConverge:
			main = map[cdm.Gt]interface{}{cdm.GotoAv: am, cdm.GotoLive: rm, cdm.GotoArticle: metam}
			op = convergem[r.Pid]
		case model.GotoBanner:
			op.FromActiveBanner(r.Items, "")
		case model.GotoTabNews:
			op.FromActive(r)
		case model.GotoTabContentRcmd:
			main = map[cdm.Gt]interface{}{cdm.GotoAv: am, cdm.GotoGame: downloadm, cdm.GotoBangumi: sm, cdm.GotoLive: rm, cdm.GotoArticle: metam, cdm.GotoSpecial: specialm, cdm.GotoPicture: picm}
			op.FromActive(r)
		case model.GotoTabEntrance:
			op.FromActive(r)
		case model.GotoTabTagRcmd:
			main = map[cdm.Gt]interface{}{cdm.GotoAv: am}
			op.FromActive(r)
			op.Items = make([]*operate.Card, 0, len(r.Items))
			for _, item := range r.Items {
				if item != nil {
					op.Items = append(op.Items, &operate.Card{ID: item.Pid, Goto: item.Goto})
				}
			}
		}
		h.From(main, op)
		if !h.Get().Right {
			continue
		}
		if hasThreePoint {
			h.Get().TabThreePointWatchLater()
		}
		is = append(is, h)
	}
	return
}
