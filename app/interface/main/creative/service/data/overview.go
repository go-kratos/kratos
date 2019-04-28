package data

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/creative/model/data"
	"go-common/app/interface/main/creative/model/tag"

	"go-common/library/ecode"
	"go-common/library/log"
)

func beginningOfDay(t time.Time) time.Time {
	d := time.Duration(-t.Hour()) * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getTuesday(now time.Time) time.Time {
	t := beginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	d := time.Duration(-weekday+2) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getSunday(now time.Time) time.Time {
	t := beginningOfDay(now)
	weekday := int(t.Weekday())
	if weekday == 0 {
		return t
	}
	d := time.Duration(7-weekday) * 24 * time.Hour
	return t.Truncate(time.Hour).Add(d)
}

func getDate() (sd string) {
	t := time.Now()
	td := getTuesday(t).Add(12 * time.Hour)
	if t.Before(td) { //当前时间在本周二12点之前，则取上上周日的数据，否则取上周日的数据
		sd = getSunday(t.AddDate(0, 0, -14)).Format("20060102")
	} else {
		sd = getSunday(t.AddDate(0, 0, -7)).Format("20060102")
	}
	log.Info("current time (%s) tuesday (%s) sunday (%s)", t.Format("2006-01-02 15:04:05"), td, sd)
	return
}

// NewStat get stat from hbase.
func (s *Service) NewStat(c context.Context, mid int64, ip string) (r *data.Stat, err error) {
	hbaseDate1 := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	hbaseDate2 := time.Now().AddDate(0, 0, -2).Add(-12 * time.Hour).Format("20060102")
	var r1, r2 *data.UpBaseStat
	if r1, err = s.data.UpStat(c, mid, hbaseDate1); err != nil || r1 == nil {
		log.Error("s.data.NewStat error(%v) mid(%d) r1(%v) ip(%s)", err, mid, r1, ip)
		err = ecode.CreativeDataErr
		return
	}
	if r2, err = s.data.UpStat(c, mid, hbaseDate2); err != nil || r2 == nil {
		log.Error("s.data.NewStat error(%v) mid(%d) r2(%v) ip(%s)", err, mid, r2, ip)
		err = ecode.CreativeDataErr
		return
	}
	r = &data.Stat{}
	if r1 != nil {
		r.Play = r1.View
		r.Dm = r1.Dm
		r.Comment = r1.Reply
		r.Fan = r1.Fans
		r.Fav = r1.Fav
		r.Like = r1.Like
		r.Share = r1.Share
		r.Coin = r1.Coin
		r.Elec = r1.Elec
	}
	log.Info("s.data.UpStat hbaseDate1(%+v) mid(%d)", r1, mid)
	if r2 != nil {
		r.PlayLast = r2.View
		r.DmLast = r2.Dm
		r.CommentLast = r2.Reply
		r.FanLast = r2.Fans
		r.FavLast = r2.Fav
		r.LikeLast = r2.Like
		r.ShareLast = r2.Share
		r.CoinLast = r2.Coin
		r.ElecLast = r2.Elec
	}
	log.Info("s.data.UpStat hbaseDate2 (%+v) mid(%d)", r2, mid)
	pfl, err := s.acc.ProfileWithStat(c, mid)
	if err != nil {
		return
	}
	r.Fan = int64(pfl.Follower)
	return
}

// ViewerBase get up viewer base data.
func (s *Service) ViewerBase(c context.Context, mid int64) (res map[string]*data.ViewerBase, err error) {
	dt := getDate()
	// try cache
	if res, _ = s.data.ViewerBaseCache(c, mid, dt); res != nil {
		s.pCacheHit.Incr("viewer_base_cache")
		return
	}
	// from data source
	if res, err = s.data.ViewerBase(c, mid, dt); len(res) != 0 {
		s.pCacheMiss.Incr("viewer_base_cache")
		s.data.AddCache(func() {
			s.data.AddViewerBaseCache(context.Background(), mid, dt, res)
		})
	}
	return
}

// ViewerArea get up viewer area data.
func (s *Service) ViewerArea(c context.Context, mid int64) (res map[string]map[string]int64, err error) {
	dt := getDate()
	// try cache
	if res, _ = s.data.ViewerAreaCache(c, mid, dt); res != nil {
		s.pCacheHit.Incr("viewer_area_cache")
		return
	}
	// from data source
	if res, err = s.data.ViewerArea(c, mid, dt); len(res) != 0 {
		s.pCacheMiss.Incr("viewer_area_cache")
		s.data.AddCache(func() {
			s.data.AddViewerAreaCache(context.Background(), mid, dt, res)
		})
	}
	return
}

// CacheTrend get trend from mc.
func (s *Service) CacheTrend(c context.Context, mid int64) (res map[string]*data.ViewerTrend, err error) {
	dt := getDate()
	// try cache
	if res, err = s.data.TrendCache(c, mid, dt); err != nil {
		log.Error("trend s.data.TrendCache err(%v)", err)
		return
	}
	if len(res) != 0 {
		s.pCacheHit.Incr("trend_cache")
		return
	}
	// from data source
	if res, err = s.viewerTrend(c, mid, dt); err != nil {
		return
	}
	s.pCacheMiss.Incr("trend_cache")
	if len(res) != 0 {
		s.data.AddCache(func() {
			s.data.AddTrendCache(context.Background(), mid, dt, res)
		})
	}
	return
}

// ViewerTrend get up viewer trend data.
func (s *Service) viewerTrend(c context.Context, mid int64, dt string) (res map[string]*data.ViewerTrend, err error) {
	ut, err := s.data.ViewerTrend(c, mid, dt)
	if err != nil || ut == nil {
		log.Error("trend s.data.ViewerTrend err(%v)", err)
		return
	}
	f := []string{"fan", "not_fan"}
	skeys := make([]int, 0) //for tag sort.
	tgs := make([]int64, 0) // for request tag name.
	res = make(map[string]*data.ViewerTrend)
	for _, fk := range f {
		td := ut[fk]
		vt := &data.ViewerTrend{}
		if td == nil {
			vt.Ty = nil
			vt.Tag = nil
			res[fk] = vt
			continue
		}
		tg := make(map[int]string)   //return tag map to user.
		ty := make(map[string]int64) //return type map to user.
		//deal type for type name.
		if td.Ty != nil {
			for k, v := range td.Ty {
				ke := int16(k)
				if t, ok := s.p.TypeMapCache[ke]; ok {
					ty[t.Name] = v
				}
			}
		} else {
			ty = nil
		}
		// deal tag for tag name.
		if td.Tag != nil {
			for k, v := range td.Tag {
				tgs = append(tgs, v)
				skeys = append(skeys, k)
			}
			var tlist []*tag.Meta
			if tlist, err = s.dtag.TagList(c, tgs); err != nil {
				log.Error("trend s.dtag.TagList err(%v)", err)
			}
			tNameMap := make(map[int64]string)
			for _, v := range tlist {
				tNameMap[v.TagID] = v.TagName
			}
			for _, k := range skeys {
				if _, ok := tNameMap[td.Tag[k]]; ok {
					tg[k] = tNameMap[td.Tag[k]]
				}
			}
		} else {
			tg = nil
		}
		vt.Ty = ty
		vt.Tag = tg
		res[fk] = vt
	}
	return
}

// RelationFansDay get up viewer trend data.
func (s *Service) RelationFansDay(c context.Context, mid int64) (res map[string]map[string]int, err error) {
	dt := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	// try cache
	if res, _ = s.data.RelationFansDayCache(c, mid, dt); res != nil {
		s.pCacheHit.Incr("relation_fans_day_cache")
		return
	}
	// from data source
	if res, err = s.data.RelationFansDay(c, mid); len(res) != 0 {
		s.pCacheMiss.Incr("relation_fans_day_cache")
		s.data.AddCache(func() {
			s.data.AddRelationFansDayCache(context.Background(), mid, dt, res)
		})
	}
	return
}

// RelationFansHistory get relation history data by month.
func (s *Service) RelationFansHistory(c context.Context, mid int64, month string) (res map[string]map[string]int, err error) {
	if res, err = s.data.RelationFansHistory(c, mid, month); err != nil {
		log.Error("s.data.RelationFansHistory err(%v)", err)
	}
	return
}

// RelationFansMonth get up viewer trend data.
func (s *Service) RelationFansMonth(c context.Context, mid int64) (res map[string]map[string]int, err error) {
	dt := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	// try cache
	if res, _ = s.data.RelationFansMonthCache(c, mid, dt); res != nil {
		s.pCacheHit.Incr("relation_fans_month_cache")
		return
	}
	// from data source
	if res, err = s.data.RelationFansMonth(c, mid); len(res) != 0 {
		s.pCacheMiss.Incr("relation_fans_month_cache")
		s.data.AddCache(func() {
			s.data.AddRelationFansMonthCache(context.Background(), mid, dt, res)
		})
	}
	return
}

// ViewerActionHour get up viewer action hour data.
func (s *Service) ViewerActionHour(c context.Context, mid int64) (res map[string]*data.ViewerActionHour, err error) {
	dt := getDate()
	// try cache
	if res, _ = s.data.ViewerActionHourCache(c, mid, dt); res != nil {
		s.pCacheHit.Incr("viewer_action_hour_cache")
		return
	}
	// from data source
	if res, err = s.data.ViewerActionHour(c, mid, dt); len(res) != 0 {
		s.data.AddCache(func() {
			s.pCacheMiss.Incr("viewer_action_hour_cache")
			s.data.AddViewerActionHourCache(context.Background(), mid, dt, res)
		})
	}
	return
}

// UpIncr for Play/Dm/Reply/Fav/Share/Elec/Coin incr.
func (s *Service) UpIncr(c context.Context, mid int64, ty int8, ip string) (res map[string]*data.ViewerIncr, err error) {
	tyStr, _ := data.IncrTy(ty)
	res = make(map[string]*data.ViewerIncr)
	daytime := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour)
	datekey := daytime.Format("20060102")
	dt := daytime.Format("20060102")
	vic, _ := s.data.ViewerIncrCache(c, mid, tyStr, dt)
	if vic != nil {
		s.pCacheHit.Incr("viewer_incr_cache")
		res[datekey] = vic
		log.Info("s.data.ViewerIncrCache mid(%d) cache(%v) err(%v)", mid, vic, err)
		return
	}
	incr, _ := s.data.UpIncr(c, mid, ty, dt)
	if incr == nil {
		res[datekey] = nil
		log.Info("s.data.UpIncr mid(%d) incr(%v) err(%v)", incr, err)
		return
	}
	tyRank := make(map[string]int) //return type map to user.
	for k, v := range incr.Rank {
		ke := int16(k)
		if t, ok := s.p.TypeMapCache[ke]; ok {
			tyRank[t.Name] = v
		}
	}
	sortK := make([]int, 0, len(incr.TopAIDList))
	aids := make([]int64, 0, len(incr.TopAIDList))
	for k, v := range incr.TopAIDList {
		aids = append(aids, v)
		sortK = append(sortK, k)
	}
	avm, _ := s.p.BatchArchives(c, mid, aids, ip)
	if len(avm) == 0 {
		return
	}
	sort.Ints(sortK)
	arcs := make([]*data.ArcInc, 0, len(avm))
	for _, k := range sortK {
		if aid, ok := incr.TopAIDList[k]; ok {
			if av, ok := avm[aid]; ok {
				al := &data.ArcInc{}
				al.AID = av.Archive.Aid
				al.PTime = av.Archive.PTime
				al.Title = av.Archive.Title
				al.DayTime = daytime.Unix()
				if _, ok := incr.TopIncrList[k]; ok {
					al.Incr = incr.TopIncrList[k]
				}
				arcs = append(arcs, al)
			}
		}
	}
	vi := &data.ViewerIncr{}
	vi.Arcs = arcs
	vi.TotalIncr = incr.Incr
	if len(tyRank) == 0 {
		vi.TyRank = nil
	} else {
		vi.TyRank = tyRank
	}
	res[datekey] = vi
	// insert cache. NOTE: sync reason?
	s.data.AddCache(func() {
		s.pCacheMiss.Incr("viewer_incr_cache")
		s.data.AddViewerIncrCache(context.Background(), mid, tyStr, dt, vi)
	})
	return
}
