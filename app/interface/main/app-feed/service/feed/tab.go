package feed

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/app-card/model/card/bangumi"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	tag "go-common/app/interface/main/tag/model"
	article "go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

func (s *Service) Menus(c context.Context, plat int8, build int, now time.Time) (menus []*operate.Menu) {
	memuCache := s.menuCache
	menus = make([]*operate.Menu, 0, len(memuCache))
LOOP:
	for _, m := range memuCache {
		if vs, ok := m.Versions[plat]; ok {
			for _, v := range vs {
				if model.InvalidBuild(build, v.Build, v.Condition) {
					continue LOOP
				}
			}
			if m.Status == 1 && (m.STime == 0 || now.After(m.STime.Time())) && (m.ETime == 0 || now.Before(m.ETime.Time())) {
				if m.ID == s.c.Bnj.TabID {
					m.Img = s.c.Bnj.TabImg
				}
				menus = append(menus, m)
			}
		}
	}
	return
}

// Actives return actives
func (s *Service) Actives(c context.Context, id, mid int64, now time.Time) (items []*feed.Item, cover string, isBnj bool, bnjDays int, err error) {
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
	if items, err = s.dealTab(c, rs, mid, now); err != nil {
		log.Error("%+v", err)
		return
	}
	cover = s.coverCache[id]
	return
}

func (s *Service) dealTab(c context.Context, rs []*operate.Active, mid int64, now time.Time) (is []*feed.Item, err error) {
	if len(rs) == 0 {
		is = _emptyItem
		return
	}
	var (
		aids, tids, roomIDs, sids, metaIDs []int64
		am                                 map[int64]*archive.ArchiveWithPlayer
		rm                                 map[int64]*live.Room
		sm                                 map[int64]*bangumi.Season
		metam                              map[int64]*article.Meta
		tagm                               map[int64]*tag.Tag
	)
	convergem := map[int64]*operate.Converge{}
	downloadm := map[int64]*operate.Download{}
	for _, r := range rs {
		switch r.Type {
		case model.GotoPlayer:
			if r.Pid != 0 {
				aids = append(aids, r.Pid)
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
					item := &operate.Active{Pid: aid, Goto: model.GotoAv}
					r.Items = append(r.Items, item)
					aids = append(aids, aid)
				}
			}
		case model.GotoConverge:
			if card, ok := s.convergeCache[r.Pid]; ok {
				for _, item := range card.Items {
					switch item.Goto {
					case model.GotoAv:
						if item.Pid != 0 {
							aids = append(aids, item.Pid)
						}
					case model.GotoLive:
						if item.Pid != 0 {
							roomIDs = append(roomIDs, item.Pid)
						}
					case model.GotoArticle:
						if item.Pid != 0 {
							metaIDs = append(metaIDs, item.Pid)
						}
					}
				}
				convergem[r.Pid] = card
			}
		case model.GotoTabEntrance, model.GotoTabContentRcmd:
			for _, item := range r.Items {
				switch item.Goto {
				case model.GotoAv:
					if item.Pid != 0 {
						aids = append(aids, item.Pid)
					}
				case model.GotoLive:
					if item.Pid != 0 {
						roomIDs = append(roomIDs, item.Pid)
					}
				case model.GotoBangumi:
					if item.Pid != 0 {
						sids = append(sids, item.Pid)
					}
				case model.GotoGame:
					if card, ok := s.downloadCache[item.Pid]; ok {
						downloadm[item.Pid] = card
					}
				case model.GotoArticle:
					if item.Pid != 0 {
						metaIDs = append(metaIDs, item.Pid)
					}
				}
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
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	is = make([]*feed.Item, 0, len(rs))
	for _, r := range rs {
		i := &feed.Item{}
		switch r.Type {
		case model.GotoPlayer:
			if a, ok := am[r.Pid]; ok {
				i.FromPlayer(a)
				is = append(is, i)
			}
		case model.GotoPlayerLive:
			if room, ok := rm[r.Pid]; ok {
				i.FromPlayerLive(room)
				if i.Goto != "" {
					is = append(is, i)
				}
			}
		case model.GotoSpecial:
			if sc, ok := s.specialCache[r.Pid]; ok {
				i.FromSpecial(sc.ID, sc.Title, sc.Cover, sc.Desc, sc.ReValue, sc.ReType, sc.Badge, sc.Size)
			}
			if i.Goto != "" {
				is = append(is, i)
			}
		case model.GotoConverge:
			if cc, ok := convergem[r.Pid]; ok {
				i.FromConverge(cc, am, rm, metam)
				if i.Goto != "" {
					is = append(is, i)
				}
			}
		case model.GotoTabTagRcmd:
			i.FromTabTags(r, am, tagm)
			if i.Goto != "" {
				is = append(is, i)
			}
		case model.GotoTabEntrance, model.GotoTabContentRcmd:
			i.FromTabCards(r, am, downloadm, sm, rm, metam, s.specialCache)
			if i.Goto != "" {
				is = append(is, i)
			}
		case model.GotoBanner:
			i.FromTabBanner(r)
			if i.Goto != "" {
				is = append(is, i)
			}
		case model.GotoTabNews:
			i.FromNews(r)
			if i.Goto != "" {
				is = append(is, i)
			}
		}
	}
	if len(is) == 0 {
		is = _emptyItem
	}
	return
}

func (s *Service) loadTabCache() {
	c := context.TODO()
	menus, err := s.tab.Menus(c)
	if err != nil {
		log.Error("%+v", err)
	} else {
		s.menuCache = menus
	}
	acs, err := s.tab.Actives(c)
	if err != nil {
		log.Error("%+v", err)
	} else {
		s.tabCache, s.coverCache = mergeTab(acs)
	}
}

func mergeTab(acs []*operate.Active) (tabm map[int64][]*operate.Active, coverm map[int64]string) {
	coverm = make(map[int64]string, len(acs))
	parentm := make(map[int64]struct{}, len(acs))
	for _, ac := range acs {
		if ac.Type == model.GotoTabBackground {
			parentm[ac.ID] = struct{}{}
			coverm[ac.ID] = ac.Cover
		}
	}
	sort.Sort(operate.Actives(acs))
	tabm = make(map[int64][]*operate.Active, len(acs))
	for parentID := range parentm {
		for _, ac := range acs {
			if ac.ParentID == parentID {
				tabm[ac.ParentID] = append(tabm[ac.ParentID], ac)
			}
		}
	}
	return
}

func (s *Service) tabproc() {
	for {
		time.Sleep(time.Minute * 1)
		s.loadTabCache()
	}
}
