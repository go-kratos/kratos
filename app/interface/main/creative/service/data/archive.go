package data

import (
	"context"
	"sort"
	"sync"
	"time"

	"go-common/app/interface/main/creative/model/data"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

// UpPlaySourceAnalysis get play analysis.
func (s *Service) UpPlaySourceAnalysis(c context.Context, mid int64) (res *data.PlaySource, err error) {
	if res, err = s.data.UpPlaySourceAnalysis(c, mid); err != nil {
		log.Error("s.data.UpPlaySourceAnalysis err(%v)", err)
	}
	return
}

// reverse reverses a slice of ints in place.
func reverse(s []int64) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// UpArcPlayAnalysis get play list.
func (s *Service) UpArcPlayAnalysis(c context.Context, mid int64, cp int, ip string) (res *data.ArchivePlayList, err error) {
	var (
		originAIDs []int64
	)
	res = &data.ArchivePlayList{} //初始化返回值
	dt := "0030"                  //最近30个稿件
	originAIDs, err = s.data.UpArcQuery(c, mid, dt, cp)
	if err != nil {
		log.Error("s.data.UpArcQuery mid(%d)|err(%v)", err, mid)
		return
	}
	log.Info("s.data.UpArcQuery originAIDs(%+v)|len(%d)|mid(%d)", originAIDs, len(originAIDs), mid)
	reverse(originAIDs) //反转获取最近aid list.
	aids := make([]int64, 0, 30)
	aidsMap := make(map[int]int64)
	sortK := make([]int, 0, 30)
	for k, aid := range originAIDs { //取最近30个aid
		aidsMap[k] = aid
		sortK = append(sortK, k)
		if len(sortK) == 30 {
			break
		}
	}
	sort.Slice(sortK, func(i, j int) bool { //30个aid索引从小到大排列
		return sortK[i] < sortK[j]
	})
	for _, ak := range sortK {
		if aid, ok := aidsMap[ak]; ok {
			aids = append(aids, aid)
		}
	}
	count := len(aids)
	log.Info("s.data.UpArcQuery aids(%+v)|len(%d)|mid(%d)", aids, count, mid)
	if count == 0 {
		return
	}
	aps := make([]*data.ArchivePlay, 0, count)
	apsMap := make(map[int64]*data.ArchivePlay, count)
	var (
		wg  sync.WaitGroup
		l   sync.RWMutex
		avm map[int64]*api.Arc
	)
	start := time.Now()
	for _, aid := range aids { //并发获取aid对应的播放信息
		wg.Add(1)
		go func(aid int64) {
			defer wg.Done()
			start = time.Now()
			var ap *data.ArchivePlay
			if ap, err = s.data.UpArcPlayAnalysis(context.Background(), aid); err != nil {
				log.Error("s.data.UpArcPlayAnalysis err(%v)", err)
			}
			elapsed := time.Since(start)
			log.Info("s.data.UpArcPlayAnalysis ap(%+v) elapsed(%v)", ap, elapsed)
			if ap != nil {
				l.Lock()
				apsMap[aid] = ap
				l.Unlock()
			}
		}(aid)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if avm, err = s.arc.Archives(c, aids, ip); err != nil {
			log.Error("s.arc.Archives mid(%d)|aids(%+v)|ip(%s)|err(%v)", mid, aids, ip, err)
			err = nil
		}
	}()
	wg.Wait()
	elapsed := time.Since(start)
	log.Info("s.data.UpArcPlayAnalysis aids(%+v)|apsMap(%+v)|len(apsMap|%d)|len(avm|%d)|elapsed(%v)", apsMap, aids, len(avm), len(apsMap), elapsed)
	for _, aid := range aids {
		if ap, ok := apsMap[aid]; ok {
			if av, ok := avm[aid]; ok && av != nil && ap != nil {
				ap.Title = av.Title
			}
			aps = append(aps, ap)
		}
	}
	sort.Slice(aps, func(i, j int) bool { //30个aid按创建时间从小到大排列
		return aps[i].CTime < aps[j].CTime
	})
	res.ArcPlayList = aps
	return
}

// ThirtyDayArchive for Play/Dm/Reply/Fav/Share/Elec/Coin for 30 days.
func (s *Service) ThirtyDayArchive(c context.Context, mid int64, ty int8) (res []*data.ThirtyDay, err error) {
	tyStr, _ := data.IncrTy(ty)
	if res, err = s.data.ThirtyDayArchiveCache(c, mid, tyStr); err != nil && len(res) != 0 {
		log.Info("creatorDataArchive mid(%d) type(%d) cache hit", mid, ty)
		return
	}
	if res, err = s.data.ThirtyDayArchive(c, mid, ty); len(res) != 0 {
		timeSlice := make([]int64, 0, 30)
		resMap := make(map[int64]int64)
		for i := 30; i > 0; i-- { //从昨天往前推算30天生成30天对应的时间戳
			var tm time.Time
			dt := time.Now().AddDate(0, 0, -1-i).Add(-12 * time.Hour).Format("20060102")
			tm, err = time.Parse("20060102", dt)
			if err != nil {
				log.Error("time.Parse error(%v)", err)
				return
			}
			t := tm.Unix()
			timeSlice = append(timeSlice, t)
		}
		for _, v := range res { //映射原始时间戳集合
			resMap[v.DateKey] = v.DateKey
		}
		if len(res) < 30 { //如果当前的数据不够30条则补齐30条
			for _, v := range timeSlice {
				if _, ok := resMap[v]; !ok { //当前时间戳不在原始时间戳集合存在则使用生成的时间戳
					td := &data.ThirtyDay{
						DateKey: v,
					}
					res = append(res, td)
				}
				if len(res) == 30 { //添加到30条为止
					break
				}
			}
			log.Info("s.data.ThirtyDayArchive append mid(%d) type(%d) return data(%+v)", mid, ty, res)
		}
		tds := make([]*data.ThirtyDay, 0, 30)
		skeys := make([]int64, 0, 30)
		tdMap := make(map[int64]int64)
		for _, v := range res {
			if v != nil {
				tdMap[v.DateKey] = v.TotalIncr   //设置时间戳和动态数据的键值对
				skeys = append(skeys, v.DateKey) //获取时间戳集合
			}
		}
		sort.Slice(skeys, func(i, j int) bool { //对时间戳做降序排列
			return skeys[i] > skeys[j]
		})
		for _, k := range skeys { //按照排序好的时间戳组装最终的返回数据
			if v, ok := tdMap[k]; ok {
				td := &data.ThirtyDay{}
				td.DateKey = k
				td.TotalIncr = v
				tds = append(tds, td)
				if len(tds) == 30 { //只取30条
					break
				}
			}
		}
		log.Info("s.data.ThirtyDayArchive mid(%d) type(%d) cache miss res(%+v)", mid, ty, tds)
		res = tds
		s.data.AddCache(func() {
			s.data.AddThirtyDayArchiveCache(context.Background(), mid, tyStr, tds)
		})
	}
	return
}
