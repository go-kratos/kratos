package data

import (
	"context"
	"sort"

	"go-common/app/interface/main/creative/model/data"
	"go-common/app/interface/main/creative/model/medal"
	account "go-common/app/service/main/account/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	rankKeys = []string{data.PlayDuration, data.VideoAct, data.DynamicAct}
)

// UpFansAnalysisForWeb get web fans data.
func (s *Service) UpFansAnalysisForWeb(c context.Context, mid int64, ip string) (res *data.WebFan, err error) {
	var (
		origin  *data.WebFan
		rkList  map[string][]*data.RankInfo
		mdlRank []*medal.FansRank
	)
	if origin, err = s.data.UpFansAnalysisForWeb(c, mid, data.Thirty); err != nil {
		log.Error("s.data.UpFansAnalysisForWeb err(%v)", err)
		return
	}
	if origin == nil {
		return
	}
	group, ctx := errgroup.WithContext(c)
	group.Go(func() (err error) {
		if rkList, err = s.getTopList(ctx, mid, origin.RankMap, ip); err != nil {
			log.Error("s.getTopList mid(%d)|err(%v)", mid, err)
			err = nil
		}
		return
	})
	group.Go(func() (err error) {
		if mdlRank, err = s.medal.Rank(ctx, mid); err != nil {
			log.Error("s.medal.Rank mid(%d)|err(%v)", mid, err)
			err = nil
		}
		return
	})
	group.Wait()
	if len(rkList) == 0 {
		log.Info("s.getTopList is empty rkList(%+v) ", rkList)
		rkList = make(map[string][]*data.RankInfo)
	}
	for _, key := range rankKeys {
		if v, ok := rkList[key]; ok {
			rkList[key] = v
		} else {
			rkList[key] = nil
		}
	}
	log.Info("s.getTopList rkList(%+v)|len(%d)", rkList, len(rkList))
	res = &data.WebFan{
		Summary:  origin.Summary,
		Source:   origin.Source,
		RankList: rkList,
		RankMedal: map[string][]*medal.FansRank{
			"medal": mdlRank,
		},
	}
	return
}

// sort mid list and get uname/avatar info
func (s *Service) getTopList(c context.Context, mid int64, RankMap map[string]map[string]int32, ip string) (rkList map[string][]*data.RankInfo, err error) {
	rkList = make(map[string][]*data.RankInfo)
	for _, key := range rankKeys {
		vd, ok := RankMap[key] //map playduration videoact dynamicact top10 mids.
		if !ok {
			continue
		}
		sortK := make([]string, 0, 10)
		for key := range vd { //sort top mids by key, for example (dr2,dr1,dr3...)
			sortK = append(sortK, key)
		}
		sort.Slice(sortK, func(i, j int) bool { //for example (dr1,dr2,dr3...)
			return sortK[i] < sortK[j]
		})
		mids := make([]int64, 0, 10)
		for _, k := range sortK { //get sort mids slice.
			if mid, ok := vd[k]; ok {
				mids = append(mids, int64(mid))
			}
		}
		if len(mids) == 0 {
			continue
		}
		var (
			g, ctx    = errgroup.WithContext(c)
			users     map[int64]*account.Info
			followers map[int64]int
		)
		g.Go(func() error { //获取用户信息
			if users, err = s.acc.Infos(ctx, mids, ip); err != nil {
				log.Error("s.acc.Infos mid(%d)|mids(%v)|ip(%s)|err(%v)", mid, mids, ip, err)
			}
			return err
		})
		g.Go(func() error { //获取关注状态
			if followers, err = s.acc.Relations2(ctx, mid, mids, ip); err != nil {
				log.Error("s.acc.Relations mid(%d)|ip(%s)|err(%v)", mid, ip, err)
			}
			return err
		})
		if g.Wait() != nil {
			return
		}
		rkInfos := make([]*data.RankInfo, 0, 10)
		for _, m := range mids {
			ri := &data.RankInfo{}
			if v, ok := users[m]; ok { //get userinfo by sorted mids.
				ri.MID = v.Mid
				ri.Uname = v.Name
				ri.Photo = v.Face
			}
			if fl, ok := followers[m]; ok { //get relation by sorted mids.
				ri.Relation = fl
			}
			rkInfos = append(rkInfos, ri) //append top10 useinfo.
			if len(rkInfos) == 10 {
				break
			}
		}
		rkList[key] = rkInfos //map playduration videoact dynamicact top10 userinfo.
	}
	log.Info("getTopList mid(%d)|rkList(%+v) ", mid, rkList)
	return
}

//UpFansMedal for get fan medel count.
func (s *Service) UpFansMedal(c context.Context, mid int64) (fansMdl *data.UpFansMedal, err error) {
	if fansMdl, err = s.data.UpFansMedal(c, mid); err != nil {
		log.Error("s.data.UpFansMedal mid(%d)|err(%v)", mid, err)
	}
	return
}
