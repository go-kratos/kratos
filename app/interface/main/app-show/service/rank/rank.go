package rank

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-show/conf"
	accdao "go-common/app/interface/main/app-show/dao/account"
	arcdao "go-common/app/interface/main/app-show/dao/archive"
	adtdao "go-common/app/interface/main/app-show/dao/audit"
	rcmmndao "go-common/app/interface/main/app-show/dao/recommend"
	rgdao "go-common/app/interface/main/app-show/dao/region"
	reldao "go-common/app/interface/main/app-show/dao/relation"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/region"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_initRank = "rank_key_%s_%d"
)

var (
	// 番剧 动画，音乐，舞蹈，游戏，科技，娱乐，鬼畜，电影，时尚, 生活，连载番剧（二级分区），国漫，影视，纪录片，国创相关，数码
	_tids = []int{13, 1, 3, 129, 4, 36, 5, 119, 23, 155, 160, 11, 33, 167, 181, 177, 168, 188}
	// region.ShowItem
	_emptyShowItems = []*region.ShowItem{}
	// _pgctids        = map[int]struct{}{
	// 	177: struct{}{},
	// }
)

type Service struct {
	c *conf.Config
	// region
	rdao *rgdao.Dao
	// rcmmnd
	rcmmnd *rcmmndao.Dao
	// archive
	arc *arcdao.Dao
	// audit
	adt *adtdao.Dao
	// account
	accd *accdao.Dao
	// relation
	reldao *reldao.Dao
	// tick
	tick time.Duration
	// ranking
	rankCache     map[string][]*region.ShowItem
	rankOseaCache map[string][]*region.ShowItem
	// audit cache
	auditCache map[string]map[int]struct{} // audit mobi_app builds
}

// New new a region service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		rdao:   rgdao.New(c),
		rcmmnd: rcmmndao.New(c),
		// archive
		arc: arcdao.New(c),
		// audit
		adt: adtdao.New(c),
		// account
		accd: accdao.New(c),
		// relation
		reldao: reldao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// ranking
		rankCache:     map[string][]*region.ShowItem{},
		rankOseaCache: map[string][]*region.ShowItem{},
		// audit cache
		auditCache: map[string]map[int]struct{}{},
	}
	time.Sleep(time.Second * 2)
	s.load()
	s.loadAuditCache()
	go s.loadproc()
	return
}

// RankShow
func (s *Service) RankShow(c context.Context, plat int8, rid, pn, ps int, mid int64, order string) (res []*region.ShowItem) {
	var (
		key             = fmt.Sprintf(_initRank, order, rid)
		tmp             []*region.ShowItem
		authorMids      []int64
		authorMidExist  = map[int64]struct{}{}
		authorRelations map[int64]*account.Relation
		authorStats     map[int64]*relation.Stat
		authorCards     map[int64]*account.Card
		err             error
	)
	if model.IsOverseas(plat) {
		tmp = s.rankOseaCache[key]
	} else {
		tmp = s.rankCache[key]
	}
	start := (pn - 1) * ps
	end := start + ps
	if len(tmp) == 0 || start > len(tmp) {
		res = _emptyShowItems
		return
	}
	if end < len(tmp) {
		tmp = tmp[start:end]
	} else {
		tmp = tmp[start:]
	}
	for _, t := range tmp {
		i := &region.ShowItem{}
		*i = *t
		res = append(res, i)
		// up mid
		if _, ok := authorMidExist[i.Mid]; !ok && i.Mid > 0 {
			authorMids = append(authorMids, i.Mid)
			authorMidExist[i.Mid] = struct{}{}
		}
	}
	if len(authorMids) > 0 {
		g, ctx := errgroup.WithContext(c)
		g.Go(func() error {
			if authorCards, err = s.accd.Cards3(ctx, authorMids); err != nil {
				log.Error("s.accd.Cards3 error(%v)", err)
			}
			return nil
		})
		if mid > 0 {
			g.Go(func() error {
				if authorRelations, err = s.accd.Relations3(ctx, mid, authorMids); err != nil {
					log.Error("s.accd.Relations2 error(%v)", err)
				}
				return nil
			})
		}
		g.Go(func() error {
			if authorStats, err = s.reldao.Stats(ctx, authorMids); err != nil {
				log.Error("s.reldao.Stats error(%v)", err)
			}
			return nil
		})
		if err = g.Wait(); err != nil {
			log.Error("RankUser errgroup.WithContext error(%v)", err)
		}
	}
	for _, i := range res {
		if len(authorRelations) > 0 {
			if relations, ok := authorRelations[i.Mid]; ok {
				if relations.Following {
					i.Attribute = 1
				}
			}
		}
		if len(authorStats) > 0 {
			if stats, ok := authorStats[i.Mid]; ok {
				i.Follower = int(stats.Follower)
			}
		}
		if len(authorCards) > 0 {
			if info, ok := authorCards[i.Mid]; ok {
				ov := &region.OfficialVerify{}
				ov.FromOfficialVerify(info.Official)
				i.OfficialVerify = ov
			}
		}
		if !model.IsIPad(plat) {
			if i.RedirectURL != "" {
				i.URI = i.RedirectURL
				i.Goto = model.GotoBangumi
			}
		}
	}
	return
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.load()
		s.loadAuditCache()
	}
}

// load load Rank all
func (s *Service) load() {
	var (
		tmp     = map[string][]*region.ShowItem{}
		tmpOsea = map[string][]*region.ShowItem{}
	)
	for _, rid := range _tids {
		if rid == 33 {
			aids, others, scores, err := s.rcmmnd.RankAppBangumi(context.TODO())
			key := fmt.Sprintf(_initRank, "bangumi", 0)
			if err != nil || len(aids) < 5 {
				log.Error("s.rcmmnd.RankAppBangumi len lt 20 OR error(%v)", err)
				tmp[key], tmpOsea[key] = s.rankCache[key], s.rankOseaCache[key]
				continue
			}
			tmp[key], tmpOsea[key] = s.fromRankAids(context.TODO(), aids, others, scores)
			log.Info("loadRankBangumi success")
		} else {
			aids, others, scores, err := s.rcmmnd.RankAppRegion(context.TODO(), rid)
			key := fmt.Sprintf(_initRank, "all", rid)
			if err != nil || len(aids) < 5 {
				log.Error("s.rcmmnd.RankAppRegion rid (%v) len lt 20 OR error(%v)", rid, err)
				tmp[key], tmpOsea[key] = s.rankCache[key], s.rankOseaCache[key]
				continue
			}
			tmp[key], tmpOsea[key] = s.fromRankAids(context.TODO(), aids, others, scores)
			log.Info("loadRankRegion(%s_%d) success", "all", rid)
		}
	}
	aids, others, scores, err := s.rcmmnd.RankAppAll(context.TODO())
	key := fmt.Sprintf(_initRank, "all", 0)
	if err != nil || len(aids) < 5 {
		log.Error("s.rcmmnd.RankAppAll(%s) len lt 20 OR error(%v)", "all", err)
		return
	}
	tmp[key], tmpOsea[key] = s.fromRankAids(context.TODO(), aids, others, scores)
	log.Info("loadRank(%s) success", "all")
	aids, others, scores, err = s.rcmmnd.RankAppOrigin(context.TODO())
	key = fmt.Sprintf(_initRank, "origin", 0)
	if err != nil || len(aids) < 5 {
		log.Error("s.rcmmnd.RankAppOrigin(%s) len lt 20 OR error(%v)", "all", err)
		return
	}
	tmp[key], tmpOsea[key] = s.fromRankAids(context.TODO(), aids, others, scores)
	log.Info("loadRank(%s) success", "origin")
	if len(tmp) > 0 {
		s.rankCache = tmp
	}
	if len(tmpOsea) > 0 {
		s.rankOseaCache = tmpOsea
	}
}

// fromRankAids
func (s *Service) fromRankAids(ctx context.Context, aids []int64, others, scores map[int64]int64) (sis, sisOsea []*region.ShowItem) {
	var (
		aid  int64
		as   map[int64]*api.Arc
		arc  *api.Arc
		ok   bool
		err  error
		paid int64
	)
	if as, err = s.arc.ArchivesPB(ctx, aids); err != nil {
		log.Error("s.arc.ArchivesPB error(%v)", err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", aids)
		return
	}
	child := map[int64][]*region.ShowItem{}
	childOsea := map[int64][]*region.ShowItem{}
	for _, aid = range aids {
		if arc, ok = as[aid]; ok {
			if paid, ok = others[arc.Aid]; ok {
				i := &region.ShowItem{}
				i.FromArchivePBRank(arc, scores)
				child[paid] = append(child[paid], i)
				if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
					childOsea[paid] = append(childOsea[paid], i)
				}
			}
		}
	}
	for _, aid = range aids {
		if arc, ok = as[aid]; ok {
			if _, ok = others[arc.Aid]; !ok {
				i := &region.ShowItem{}
				i.FromArchivePBRank(arc, scores)
				if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
					if tmpchild, ok := childOsea[arc.Aid]; ok {
						i.Children = tmpchild
					}
					sisOsea = append(sisOsea, i)
				}
				if tmpchild, ok := child[arc.Aid]; ok {
					i.Children = tmpchild
				}
				sis = append(sis, i)
			}
		}
	}
	return
}
