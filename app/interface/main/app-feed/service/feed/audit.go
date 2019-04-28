package feed

import (
	"context"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card"
	"go-common/app/interface/main/app-card/model/card/banner"
	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/interface/main/app-feed/model/feed"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

var (
	_auditBanners = []*banner.Banner{
		&banner.Banner{
			Title: "充电",
			Image: "http://i0.hdslb.com/bfs/archive/9ce8f6cdf76e6cbd50ce7db76262d5a35e594c79.png",
			Hash:  "3c4990d06c46de0080e3821fca6bedca",
			URI:   "bilibili://video/813060",
		},
	}
	// av2314237 已经删除，后续看情况处理
	// 已删除失效id已更换新的稿件id
	_aids = []int64{308040, 2431658, 2432648, 2427553, 539600, 1968681, 850424, 887861, 1960912, 1935680, 1406019,
		1985297, 1977493, 2312184, 2316891, 864845, 1986932, 880857, 875624, 744299}
)

// Audit check audit plat then return audit data.
func (s *Service) Audit(c context.Context, mobiApp string, plat int8, build int) (is []*feed.Item, ok bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			return s.auditData(c), true
		}
	}
	return
}

// Audit2 check audit plat and ip, then return audit data.
func (s *Service) Audit2(c context.Context, mobiApp string, plat int8, build int, column cdm.ColumnStatus) (is []card.Handler, ok bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			return s.auditData2(c, plat, column), true
		}
	}
	return
}

// auditData some data for audit.
func (s *Service) auditData(c context.Context) (is []*feed.Item) {
	i := &feed.Item{}
	i.FromBanner(_auditBanners, "")
	is = append(is, i)
	am, err := s.ArchivesWithPlayer(c, _aids, 0, "", 0, 0, 0, 0)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	for _, aid := range _aids {
		if a, ok := am[aid]; ok {
			i := &feed.Item{}
			i.FromAv(a)
			is = append(is, i)
		}
	}
	return
}

// auditData2 some data for audit.
func (s *Service) auditData2(c context.Context, plat int8, column cdm.ColumnStatus) (is []card.Handler) {
	i := card.Handle(plat, model.GotoBanner, "", column, nil, nil, nil, nil, nil)
	if i != nil {
		op := &operate.Card{}
		op.FromBanner(_auditBanners, "")
		i.From(nil, op)
		is = append(is, i)
	}
	am, err := s.arc.Archives(c, _aids)
	if err != nil {
		log.Error("%+v", err)
	}
	var main interface{}
	for _, aid := range _aids {
		if a, ok := am[aid]; ok {
			i := card.Handle(plat, model.GotoAv, "", column, nil, nil, nil, nil, nil)
			if i == nil {
				continue
			}
			op := &operate.Card{}
			op.From(cdm.CardGotoAv, aid, 0, 0, 0)
			main = map[int64]*archive.ArchiveWithPlayer{a.Aid: &archive.ArchiveWithPlayer{Archive3: archive.BuildArchive3(a)}}
			i.From(main, op)
			if !i.Get().Right {
				continue
			}
			if model.IsIPad(plat) {
				// ipad卡片不展示标签
				i.Get().DescButton = nil
			}
			is = append(is, i)
		}
	}
	return
}

func (s *Service) loadAuditCache() {
	as, err := s.adt.Audits(context.Background())
	if err != nil {
		log.Error("s.adt.Audits error(%v)", err)
		return
	}
	s.auditCache = as
}

// auditproc load audit cache.
func (s *Service) auditproc() {
	for {
		time.Sleep(s.tick)
		s.loadAuditCache()
	}
}
