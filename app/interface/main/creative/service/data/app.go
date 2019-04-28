package data

import (
	"context"
	"sort"
	"sync"
	"time"

	"go-common/library/ecode"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/data"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"math"
)

var (
	zeroSummary = map[string]int64{
		"total": 0,
		"inc":   0,
		"play":  0,
		"dm":    0,
		"reply": 0,
		"coin":  0,
		"inter": 0,
		"vv":    0,
		"da":    0,
		"re":    0,
		"co":    0,
		"fv":    0,
		"sh":    0,
		"lk":    0,
	}
)

// UpFansAnalysisForApp get app fans analysis.
func (s *Service) UpFansAnalysisForApp(c context.Context, mid int64, ty int, ip string) (res *data.AppFan, err error) {
	var (
		g, ctx     = errgroup.WithContext(c)
		origin     *data.AppFan
		typeList   map[string]int64
		tagList    []*data.Rank
		viewerArea map[string]int64
		viewerBase *data.ViewerBase
	)
	res = &data.AppFan{}
	if origin, err = s.data.UpFansAnalysisForApp(c, mid, ty); err != nil {
		log.Error("s.data.UpFansAnalysisForApp mid(%d)|ip(%s)|err(%v)", mid, ip, err)
		return
	}
	if origin == nil {
		return
	}
	if origin.Summary == nil {
		origin.Summary = zeroSummary
	}
	log.Info("s.data.UpFansAnalysisForApp origin mid(%d)|origin(%+v)", mid, origin)
	g.Go(func() (err error) {
		total, ok := origin.Summary["total"]
		inc, oki := origin.Summary["inc"]
		if !ok || !oki {
			return
		}
		pfl, err := s.acc.ProfileWithStat(ctx, mid)
		if err != nil {
			err = nil
			log.Error("s.acc.ProfileWithStat mid(%d)|err(%v)", mid, err)
			return
		}
		if pfl == nil {
			log.Error("s.acc.ProfileWithStat mid(%d)|Follower(%+v) err(%v)", mid, pfl, err)
			return
		}
		origin.Summary["total"] = pfl.Follower
		origin.Summary["inc"] = inc + (pfl.Follower - total)
		return
	})
	g.Go(func() (err error) {
		if tys, tgs, err := s.fanTrend(ctx, mid); err == nil {
			typeList = tys
			tagList = tgs
		}
		return nil
	})
	g.Go(func() (err error) {
		if va, err := s.fanViewerArea(ctx, mid); err == nil {
			viewerArea = va
		}
		return nil
	})
	g.Go(func() (err error) {
		if vb, err := s.fanViewerBase(ctx, mid); err == nil {
			viewerBase = vb
		}
		return nil
	})
	g.Wait()
	res.Summary = origin.Summary
	res.TypeList = typeList
	res.TagList = tagList
	res.ViewerArea = viewerArea
	res.ViewerBase = viewerBase
	log.Info("UpFansAnalysisForApp mid(%d)|res(%+v)", mid, res)
	return
}

//FanRankApp get app fans rank list.
func (s *Service) FanRankApp(c context.Context, mid int64, ty int, ip string) (res map[string][]*data.RankInfo, err error) {
	var (
		origin *data.AppFan
		rkList map[string][]*data.RankInfo
	)
	if origin, err = s.data.UpFansAnalysisForApp(c, mid, ty); err != nil {
		log.Error("s.data.UpFansAnalysisForApp mid(%d)|ip(%s)|err(%v)", mid, ip, err)
		return
	}
	if origin == nil {
		return
	}
	log.Info("s.UpFansRankApp origin mid(%d)|origin(%+v)", mid, origin)
	if rkList, err = s.getTopList(c, mid, origin.RankMap, ip); err != nil {
		log.Error("s.getTopList err(%v)", err)
	}
	if len(rkList) == 0 { //排行列表容错，必须返回对应的key
		rkList = make(map[string][]*data.RankInfo)
	}
	for _, key := range rankKeys {
		if v, ok := rkList[key]; ok {
			rkList[key] = v
		} else {
			rkList[key] = nil
		}
	}
	res = rkList
	log.Info("s.UpFansRankApp mid(%d)|res(%+v)|len[rkList](%d)", mid, res, len(rkList))
	return
}

func (s *Service) fanTrend(c context.Context, mid int64) (tys map[string]int64, tags []*data.Rank, err error) {
	var (
		dt    = getDate()
		trend map[string]*data.ViewerTrend
	)
	if trend, err = s.viewerTrend(c, mid, dt); err != nil {
		log.Error("fanTrend s.viewerTrend mid(%d)|err(%v)", mid, err)
		return
	}
	if tr, ok := trend["fan"]; ok && tr != nil {
		tys = tr.Ty
		skeys := []int{2, 1, 4, 3, 6, 5, 8, 7, 10, 9}
		for _, k := range skeys {
			if v, ok := tr.Tag[k]; ok {
				tg := &data.Rank{
					Rank: k,
					Name: v,
				}
				tags = append(tags, tg)
			}
		}
	}
	return
}

func (s *Service) fanViewerArea(c context.Context, mid int64) (res map[string]int64, err error) {
	var (
		origin map[string]map[string]int64
	)
	if origin, err = s.data.ViewerArea(c, mid, getDate()); err != nil {
		log.Error("fanViewerBase s.data.ViewerArea mid(%d)|err(%v)", mid, err)
	}
	if len(origin) == 0 {
		return
	}
	if v, ok := origin["fan"]; ok {
		res = v
	}
	return
}

func (s *Service) fanViewerBase(c context.Context, mid int64) (res *data.ViewerBase, err error) {
	var (
		origin map[string]*data.ViewerBase
	)
	if origin, err = s.data.ViewerBase(c, mid, getDate()); err != nil {
		log.Error("fanViewerBase s.data.ViewerArea mid(%d)|err(%v)", mid, err)
	}
	if len(origin) == 0 {
		return
	}
	if v, ok := origin["fan"]; ok {
		res = v
	}
	return
}

//OverView for app data overview.
func (s *Service) OverView(c context.Context, mid int64, ty int8, ip string) (res *data.AppOverView, err error) {
	var (
		g, ctx       = errgroup.WithContext(c)
		stat         *data.Stat
		allIncr      []*data.ThirtyDay
		singleArcInc []*data.ArcInc
	)
	g.Go(func() (err error) {
		if stat, err = s.NewStat(ctx, mid, ip); err != nil {
			err = nil
			log.Error("OverView s.ThirtyDayArchive mid(%d)|err(%v)", mid, err)
		}
		return nil
	})
	g.Go(func() (err error) {
		if allIncr, err = s.ThirtyDayArchive(ctx, mid, ty); err != nil {
			err = nil
			log.Error("OverView s.ThirtyDayArchive mid(%d)|err(%v)", mid, err)
		}
		if len(allIncr) >= 7 {
			allIncr = allIncr[0:7]
		}
		return nil
	})
	g.Go(func() (err error) {
		if prop, err := s.AppUpIncr(c, mid, ty, ip); err == nil && len(prop) > 0 {
			if prop[len(prop)-1] != nil { //获取最后一天的数据
				singleArcInc = prop[len(prop)-1].Arcs
			}
		}
		return nil
	})
	g.Wait()
	res = &data.AppOverView{
		Stat:         stat,
		AllArcIncr:   allIncr,
		SingleArcInc: singleArcInc,
	}
	return
}

// ArchiveAnalyze get single archive data.
func (s *Service) ArchiveAnalyze(c context.Context, aid, mid int64, ip string) (stat *data.ArchiveData, err error) {
	// check aid valid and owner
	isWhite := false
	for _, m := range s.c.Whitelist.DataMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	a, re := s.arc.Archive(c, aid, ip)
	if re != nil {
		err = re
		return
	}
	if a == nil {
		err = ecode.AccessDenied
		return
	}
	if !isWhite && a.Author.Mid != mid {
		err = ecode.AccessDenied
		return
	}
	stat, err = s.data.ArchiveStat(c, aid)
	if err != nil {
		log.Error("s.data.ArchiveStat aid(%d)|mid(%d)|error(%v)", aid, mid, err)
		return
	}
	if stat == nil {
		return
	}
	if stat.ArchiveStat.Play >= 100 {
		stat.ArchiveAreas, err = s.data.ArchiveArea(c, aid)
	}
	log.Info("s.ArchiveAnalyze aid(%d)|mid(%d)|stat(%+v)", aid, mid, stat)
	return
}

// VideoRetention get video quit data.
func (s *Service) VideoRetention(c context.Context, cid, mid int64, ip string) (res *data.VideoQuit, err error) {
	var (
		qs      []int64
		arc     *api.Arc
		video   *archive.Video
		isWhite bool
	)
	for _, m := range s.c.Whitelist.DataMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	if video, err = s.arc.VideoByCid(c, int64(cid), ip); err != nil {
		err = ecode.AccessDenied
		return
	}
	if video == nil {
		return
	}
	if arc, err = s.arc.Archive(c, video.Aid, ip); err == nil && !isWhite && arc.Author.Mid != mid {
		err = ecode.AccessDenied
		return
	}
	qs, err = s.data.VideoQuitPoints(c, cid)
	res = &data.VideoQuit{
		Point:    appVideoQuit(qs),
		Duration: sliceDuration(qs),
	}
	log.Info("s.VideoRetention cid(%d)|mid(%d)|quitPoints(%+v)|video(%+v)|res(%+v)", cid, mid, qs, video, res)
	return
}

// >7  interval=(n/7四舍五入)
func appVideoQuit(qps []int64) []int64 {
	cnt := len(qps) + 1
	if cnt <= 7 {
		return qps
	}
	nqps := make([]int64, 0)
	interval := int64(math.Floor(float64(cnt)/7.0 + 0.5))
	for i := 1; i < 8; i++ {
		k := interval * int64(i)
		if k > int64(cnt)-1 {
			break
		}
		nqps = append(nqps, qps[k-1])
	}
	return nqps
}

func sliceDuration(qps []int64) (ds []int64) {
	cnt := len(qps) + 1
	if cnt <= 7 {
		for i := 1; i < cnt; i++ {
			ds = append(ds, int64(i)*30)
		}
		return
	}
	interval := int64(math.Floor(float64(cnt)/7.0 + 0.5))
	for i := 1; i < 8; i++ {
		k := interval * int64(i)
		if k > int64(cnt)-1 {
			break
		}
		ds = append(ds, (k)*30)

	}
	return
}

// AppVideoQuitPoints get app video quit data.
func (s *Service) AppVideoQuitPoints(c context.Context, cid, mid int64, ip string) (res []int64, err error) {
	if res, err = s.VideoQuitPoints(c, cid, mid, ip); err != nil {
		return
	}
	res = appVideoQuit(res)
	return
}

// AppUpIncr for Play/Dm/Reply/Fav/Share/Elec/Coin incr.
func (s *Service) AppUpIncr(c context.Context, mid int64, ty int8, ip string) (res []*data.AppViewerIncr, err error) {
	incr, err := s.UpIncr(c, mid, ty, ip)
	if err != nil {
		return
	}
	if len(incr) == 0 {
		return
	}
	res = make([]*data.AppViewerIncr, 0, len(incr))
	sortK := make([]string, 0, len(incr))
	for k := range incr {
		sortK = append(sortK, k)
	}
	sort.Strings(sortK)
	for _, k := range sortK {
		v, ok := incr[k]
		if !ok {
			continue
		}
		if v != nil {
			av := &data.AppViewerIncr{}
			tm, _ := time.Parse("20060102", k)
			av.DateKey = tm.Unix()
			av.Arcs = v.Arcs
			av.TotalIncr = v.TotalIncr
			trs := make([]*data.Rank, 0, len(v.TyRank))
			for rk, rv := range v.TyRank {
				tr := &data.Rank{}
				tr.Name = rk
				tr.Rank = rv
				trs = append(trs, tr)
			}
			av.TyRank = trs
			res = append(res, av)
		}
	}
	return
}

// AppStat get app archive static data.
func (s *Service) AppStat(c context.Context, mid int64) (sts *data.AppStatList, err error) {
	sts = &data.AppStatList{}
	sts.Show = 1
	if sts.Show == 0 {
		return
	}
	var (
		wg       sync.WaitGroup
		viewChan = make(chan *data.AppStat, 6)
		fanChan  = make(chan *data.AppStat, 6)
		comChan  = make(chan *data.AppStat, 6)
		dmChan   = make(chan *data.AppStat, 6)
	)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			datekey := time.Now().AddDate(0, 0, -1-i).Add(-12 * time.Hour).Format("2006-01-02")
			dt := time.Now().AddDate(0, 0, -1-i).Add(-12 * time.Hour).Format("20060102")
			res, err := s.data.UpStat(c, mid, dt)
			if err != nil {
				log.Error("s.data.UpStat mid(%d)|err(%v)", mid, err)
				return
			}
			if res == nil {
				return
			}
			viewChan <- &data.AppStat{Date: datekey, Num: res.View}
			fanChan <- &data.AppStat{Date: datekey, Num: res.Fans}
			comChan <- &data.AppStat{Date: datekey, Num: res.Reply}
			dmChan <- &data.AppStat{Date: datekey, Num: res.Dm}
		}(i)
	}
	wg.Wait()
	close(viewChan)
	close(fanChan)
	close(comChan)
	close(dmChan)
	for v := range viewChan {
		sts.View = append(sts.View, v)
	}
	for v := range fanChan {
		sts.Fans = append(sts.Fans, v)
	}
	for v := range comChan {
		sts.Comment = append(sts.Comment, v)
	}
	for v := range dmChan {
		sts.Danmu = append(sts.Danmu, v)
	}
	sort.Slice(sts.View, func(i, j int) bool {
		return sts.View[i].Date > sts.View[j].Date
	})
	sort.Slice(sts.Fans, func(i, j int) bool {
		return sts.Fans[i].Date > sts.Fans[j].Date
	})
	sort.Slice(sts.Comment, func(i, j int) bool {
		return sts.Comment[i].Date > sts.Comment[j].Date
	})
	sort.Slice(sts.Danmu, func(i, j int) bool {
		return sts.Danmu[i].Date > sts.Danmu[j].Date
	})
	// set increment num
	for i := 0; i < len(sts.View)-1; i++ {
		if sts.View[i].Num = sts.View[i].Num - sts.View[i+1].Num; sts.View[i].Num < 0 {
			sts.View[i].Num = 0
		}
		if sts.Fans[i].Num = sts.Fans[i].Num - sts.Fans[i+1].Num; sts.Fans[i].Num < 0 {
			sts.Fans[i].Num = 0
		}
		if sts.Comment[i].Num = sts.Comment[i].Num - sts.Comment[i+1].Num; sts.Comment[i].Num < 0 {
			sts.Comment[i].Num = 0
		}
		if sts.Danmu[i].Num = sts.Danmu[i].Num - sts.Danmu[i+1].Num; sts.Danmu[i].Num < 0 {
			sts.Danmu[i].Num = 0
		}
	}
	// rm last element
	if len(sts.View) > 0 {
		sts.View = sts.View[:len(sts.View)-1]
		sts.Fans = sts.Fans[:len(sts.Fans)-1]
		sts.Comment = sts.Comment[:len(sts.Comment)-1]
		sts.Danmu = sts.Danmu[:len(sts.Danmu)-1]
	}
	return
}
