package mcndao

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/dao/cache"
	"go-common/app/interface/main/mcn/dao/global"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/log"

	"github.com/bluele/gcache"
)

// RecommendPoolCache .
type RecommendPoolCache struct {
	// map[tid] list
	UpTidMap    map[int16][]*mcnmodel.McnGetRecommendPoolInfo
	TidTypeList []*mcnmodel.TidnameInfo
}

func (r *RecommendPoolCache) add(v *mcnmodel.McnGetRecommendPoolInfo) {
	if v == nil {
		return
	}
	if r.UpTidMap == nil {
		r.UpTidMap = make(map[int16][]*mcnmodel.McnGetRecommendPoolInfo)
	}

	r.UpTidMap[v.ActiveTid] = append(r.UpTidMap[v.ActiveTid], v)
	if v.ActiveTid != 0 {
		// 加入到全部分类中
		r.UpTidMap[0] = append(r.UpTidMap[0], v)
	}
}

//RecommendSortFunc sort func
type RecommendSortFunc func(p1, p2 *mcnmodel.McnGetRecommendPoolInfo) bool

//RecommendDataSorter data sorter
type RecommendDataSorter struct {
	Datas []*mcnmodel.McnGetRecommendPoolInfo
	By    RecommendSortFunc // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *RecommendDataSorter) Len() int {
	return len(s.Datas)
}

// Swap is part of sort.Interface.
func (s *RecommendDataSorter) Swap(i, j int) {
	s.Datas[i], s.Datas[j] = s.Datas[j], s.Datas[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *RecommendDataSorter) Less(i, j int) bool {
	return s.By(s.Datas[i], s.Datas[j])
}

// RecommendSortByFansDesc .
func RecommendSortByFansDesc(p1, p2 *mcnmodel.McnGetRecommendPoolInfo) bool {
	return p1.FansCount > p2.FansCount
}

// RecommendSortByFansAsc .
func RecommendSortByFansAsc(p1, p2 *mcnmodel.McnGetRecommendPoolInfo) bool {
	return p1.FansCount < p2.FansCount
}

// RecommendSortByMonthFansDesc .
func RecommendSortByMonthFansDesc(p1, p2 *mcnmodel.McnGetRecommendPoolInfo) bool {
	return p1.FansCountIncreaseMonth > p2.FansCountIncreaseMonth
}

// RecommendSortByArchiveCountDesc .
func RecommendSortByArchiveCountDesc(p1, p2 *mcnmodel.McnGetRecommendPoolInfo) bool {
	return p1.ArchiveCount > p2.ArchiveCount
}

type loadRecommandFunc func() (res *RecommendPoolCache, err error)

func cacheKeyRecommend(int64) string {
	return "recommend"
}

//rawGetRecommendPool get recommend pool
func (d *Dao) rawGetRecommendPool() (res []*mcnmodel.McnGetRecommendPoolInfo, err error) {
	var dbresult []*mcnmodel.McnUpRecommendPool
	err = d.mcndb.Where("state=?", mcnmodel.MCNUPRecommendStateOn).Find(&dbresult).Error
	if err != nil {
		log.Error("fail to get recommend pool, err=%s", err)
		return
	}

	for _, v := range dbresult {
		var info = &mcnmodel.McnGetRecommendPoolInfo{}
		info.Copy(v)
		res = append(res, info)
	}
	return
}

func (d *Dao) loadRecommendPool() (res *RecommendPoolCache, err error) {
	recommendInfos, err := d.rawGetRecommendPool()
	if err != nil {
		log.Error("get recommend fail, err=%s", err)
		return
	}

	var midMap = make(map[int64]struct{})
	// 获取mid列表
	for _, v := range recommendInfos {
		midMap[v.UpMid] = struct{}{}
	}

	var mids []int64
	for k := range midMap {
		mids = append(mids, k)
	}

	// 获取账号信息，头像
	accInfos, err := global.GetInfos(context.Background(), mids)
	if err != nil || accInfos == nil {
		log.Warn("get infos fail, err=%s", err)
	}

	var tidUnique = newTidNameUnique()
	res = new(RecommendPoolCache)
	for _, v := range recommendInfos {
		if account, ok := accInfos[v.UpMid]; ok {
			v.UpName = account.Name
		}
		v.TidName = cache.GetTidName(int64(v.ActiveTid))
		// 这里DataTypeAccumulate只是用来记录一下
		tidUnique.addTid(v.ActiveTid, v.TidName, mcnmodel.DataTypeAccumulate)
		res.add(v)
	}
	res.TidTypeList = tidUnique.getList(mcnmodel.DataTypeAccumulate)

	// 进行默认排序，按照粉丝数，降序
	for k := range res.UpTidMap {
		sort.Sort(&RecommendDataSorter{Datas: res.UpTidMap[k], By: RecommendSortByFansDesc})
	}
	return
}

func (d *Dao) getRecommendCache(keyCalc keyFunc, load loadRecommandFunc) (result *RecommendPoolCache, err error) {
	var key = keyCalc(0)
	v, err := d.localcache.Get(key)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			// load cache
			v, err = load()
			if err != nil {
				log.Error("load cache error, key=%s. err=%s", key, err)
				return
			}

			d.localcache.SetWithExpire(key, v, time.Duration(conf.Conf.RankCache.RecommendPoolExpireTime))
		} else {
			log.Error("get from gcache err, key=%s, err=%s", key, err)
			return
		}
	}
	if v == nil {
		return
	}
	result, _ = v.(*RecommendPoolCache)
	return
}

//GetRecommendPool get recommend pool
func (d *Dao) GetRecommendPool() (res *RecommendPoolCache, err error) {
	res, err = d.getRecommendCache(cacheKeyRecommend, d.loadRecommendPool)
	if err != nil {
		log.Error("fail to get recommend pool, err=%s", err)
		return
	}
	return
}
