package mcndao

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/admin/main/up/util"
	"go-common/app/admin/main/up/util/mathutil"
	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/dao/cache"
	"go-common/app/interface/main/mcn/dao/global"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	arcgrpc "go-common/app/service/main/archive/api"
	"go-common/library/log"

	"github.com/bluele/gcache"
)

// RankByTid .
// 存储rank的缓存
// 对每种排行榜的分类进行缓存
// 对排序进行实时计算
// 先去localcache中取，取不到的话，去db中取
type RankByTid struct {
	// [tid][datatype] rank list
	TidMap         map[int16]map[mcnmodel.DataType][]mcnmodel.RankDataInterface
	TidTypeListMap map[mcnmodel.DataType][]*mcnmodel.TidnameInfo
}

//RankFunc rank func
type RankFunc func(signId int64) (result *RankByTid, err error)

type tidnameUnique struct {
	tidInfoMap map[mcnmodel.DataType]map[int16]*mcnmodel.TidnameInfo
}

func newTidNameUnique() *tidnameUnique {
	return &tidnameUnique{tidInfoMap: make(map[mcnmodel.DataType]map[int16]*mcnmodel.TidnameInfo)}
}

func (t *tidnameUnique) addTid(tid int16, name string, datatype mcnmodel.DataType) {
	var dmap map[int16]*mcnmodel.TidnameInfo
	var ok bool
	if dmap, ok = t.tidInfoMap[datatype]; !ok {
		dmap = make(map[int16]*mcnmodel.TidnameInfo)
		t.tidInfoMap[datatype] = dmap
	}
	dmap[tid] = &mcnmodel.TidnameInfo{Tid: tid, Name: name}
}

func (t *tidnameUnique) getList(datatype mcnmodel.DataType) (typeList []*mcnmodel.TidnameInfo) {
	for _, v := range t.tidInfoMap[datatype] {
		typeList = append(typeList, v)
	}
	return
}

func (t *tidnameUnique) export(dmap map[mcnmodel.DataType][]*mcnmodel.TidnameInfo) {
	for dataType, tidMap := range t.tidInfoMap {
		var typeList []*mcnmodel.TidnameInfo
		for k, v := range tidMap {
			if k == 0 {
				continue
			}
			typeList = append(typeList, v)
		}
		dmap[dataType] = typeList
	}
}

// // 排序，根据increase数量做倒序
// type rankByFansIncreaseDesc []*mcnmodel.RankUpFansInfo

// func (s rankByFansIncreaseDesc) Len() int {
// 	return len(s)
// }
// func (s rankByFansIncreaseDesc) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s rankByFansIncreaseDesc) Less(i, j int) bool {
// 	return s[i].FansIncrease > s[j].FansIncrease
// }

type sortRankFunc func(p1, p2 mcnmodel.RankDataInterface) bool

type rankDataSorter struct {
	datas []mcnmodel.RankDataInterface
	by    sortRankFunc // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *rankDataSorter) Len() int {
	return len(s.datas)
}

// Swap is part of sort.Interface.
func (s *rankDataSorter) Swap(i, j int) {
	s.datas[i], s.datas[j] = s.datas[j], s.datas[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *rankDataSorter) Less(i, j int) bool {
	return s.by(s.datas[i], s.datas[j])
}

func byValueDesc(p1, p2 mcnmodel.RankDataInterface) bool {
	return p1.GetValue() > p2.GetValue()
}

//GetList get list
func (s *RankByTid) GetList(tid int16, dataType mcnmodel.DataType) (res []mcnmodel.RankDataInterface) {
	if s.TidMap == nil {
		return
	}
	tMap, ok := s.TidMap[tid]
	if !ok || tMap == nil {
		return
	}
	if res, ok = tMap[dataType]; !ok {
		res = nil
	}
	return
}

var nulltidlist = make([]*mcnmodel.TidnameInfo, 0)

//GetTypeList get type's list
func (s *RankByTid) GetTypeList(dataType mcnmodel.DataType) (res []*mcnmodel.TidnameInfo) {
	res = s.TidTypeListMap[dataType]
	if res == nil {
		res = nulltidlist
	}
	return
}

func (s *RankByTid) addRank(v mcnmodel.RankDataInterface) {
	var dmap, allMap map[mcnmodel.DataType][]mcnmodel.RankDataInterface
	if s.TidMap == nil {
		s.TidMap = make(map[int16]map[mcnmodel.DataType][]mcnmodel.RankDataInterface)
	}
	var ok bool
	if v.GetTid() != 0 {
		if dmap, ok = s.TidMap[v.GetTid()]; !ok {
			dmap = make(map[mcnmodel.DataType][]mcnmodel.RankDataInterface)
			s.TidMap[v.GetTid()] = dmap
		}
		dmap[v.GetDataType()] = append(dmap[v.GetDataType()], v)
	}

	if allMap, ok = s.TidMap[0]; !ok {
		allMap = make(map[mcnmodel.DataType][]mcnmodel.RankDataInterface)
		s.TidMap[0] = allMap
	}

	allMap[v.GetDataType()] = append(allMap[v.GetDataType()], v)
}

func (s *RankByTid) addTidMap(v *tidnameUnique) {
	if s.TidTypeListMap == nil {
		s.TidTypeListMap = make(map[mcnmodel.DataType][]*mcnmodel.TidnameInfo)
	}

	v.export(s.TidTypeListMap)
}

//Truncate truncate all the sorted list by max item number
func (s *RankByTid) Truncate(max int) {
	for _, v := range s.TidMap {
		for k2, v2 := range v {
			var length = len(v2)
			if length == 0 {
				continue
			}
			var m = mathutil.Min(max, length)
			v[k2] = v2[0:m]
		}
	}
}

// Sort sort by sorting function
func (s *RankByTid) Sort(sortFunc sortRankFunc) {
	for _, v := range s.TidMap {
		for k2, v2 := range v {
			var sorter = &rankDataSorter{
				datas: v2,
				by:    sortFunc,
			}
			sort.Sort(sorter)
			v[k2] = v2
		}
	}
}

type keyFunc func(int64) string
type loadRankFunc func(signID int64, date time.Time) (result *RankByTid, err error)

// GetRankUpFans get fans
func (d *Dao) GetRankUpFans(signID int64) (result *RankByTid, err error) {
	return d.getRankCache(signID, cacheKeyRankFans, d.loadRankUpFansCache)
}

// GetRankArchiveLikes get fans
func (d *Dao) GetRankArchiveLikes(signID int64) (result *RankByTid, err error) {
	return d.getRankCache(signID, cacheKeyRankArchiveLikes, d.loadRankArchiveLikesCache)
}

func (d *Dao) getRankCache(signID int64, keyCalc keyFunc, load loadRankFunc) (result *RankByTid, err error) {
	var key = keyCalc(signID)
	v, err := d.localcache.Get(key)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			// load cache
			v, err = load(signID, time.Now())
			if err != nil {
				log.Error("load cache error, signID=%d, err=%s", signID, err)
				return
			}

			d.localcache.SetWithExpire(key, v, time.Duration(conf.Conf.RankCache.ExpireTime))
		} else {
			log.Error("get from gcache err, signID=%d, err=%s", signID, err)
			return
		}
	}
	if v == nil {
		return
	}
	result, _ = v.(*RankByTid)
	return
}

func cacheKeyRankFans(signID int64) string {
	return fmt.Sprintf("rank_fans_%d", signID)
}

func cacheKeyRankArchiveLikes(signID int64) string {
	return fmt.Sprintf("rank_likes_%d", signID)
}

var dataTypes = []mcnmodel.DataType{
	mcnmodel.DataTypeAccumulate,
	mcnmodel.DataTypeDay,
	mcnmodel.DataTypeWeek,
	mcnmodel.DataTypeMonth,
	mcnmodel.DataTypeActiveFans,
}

// --------------------------------------- load rank up fans -------------------------------
func (d *Dao) loadRankUpFansCache(signID int64, date time.Time) (result *RankByTid, err error) {
	rawRanks, err := d.RawRankUpFans(signID, date)
	if err != nil {
		log.Error("fail to get raw rank up fans, signid=%d, err=%s", signID, err)
		return
	}
	result = new(RankByTid)
	if len(rawRanks) == 0 {
		log.Info("up fans rank data is empty, sign id=%d", signID)
		return
	}
	var midMap = make(map[int64]struct{})
	//var accumulateMap = make(map[int64]*mcnmodel.McnRankUpFan)
	// 获取mid列表
	for _, v := range rawRanks {
		midMap[v.UpMid] = struct{}{}
		//if v.DataType == mcnmodel.DataTypeAccumulate {
		//	accumulateMap[v.UpMid] = v
		//}
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
	// 组装信息
	for _, v := range rawRanks {
		var info mcnmodel.RankUpFansInfo
		info.Copy(v)
		if account, ok := accInfos[v.UpMid]; ok {
			info.Name = account.Name
			info.UpFaceLink = account.Face
			info.TidName = cache.GetTidName(int64(info.Tid))
			if info.TidName == "" {
				info.TidName = "其他"
			}
			tidUnique.addTid(info.Tid, info.TidName, info.DataType)
		}
		//if accumulateData, ok := accumulateMap[v.UpMid]; ok {
		//	info.FansAccumulate = accumulateData.Value1
		//}
		result.addRank(&info)
	}

	result.addTidMap(tidUnique)

	// 排序
	result.Sort(byValueDesc)
	// 截断到10个，截断需要在排序之后
	result.Truncate(10)
	return
}

//RawRankUpFans get from db
func (d *Dao) RawRankUpFans(signID int64, date time.Time) (result []*mcnmodel.McnRankUpFan, err error) {
	// 有X种类型，
	// 昨日、上周、上月、累计
	// 每种类型取最近日期的数据
	for _, typ := range dataTypes {
		var tmpResult []*mcnmodel.McnRankUpFan
		e := d.mcndb.Raw(`select * from mcn_rank_up_fans where data_type=? and sign_id=?
			and generate_date=(select generate_date from mcn_rank_up_fans where data_type=? and sign_id=? and generate_date <= ? order by generate_date desc limit 1)`,
			typ, signID, typ, signID, date).
			Find(&tmpResult).Error
		if e != nil {
			log.Error("fail to get rank, type=%d, sign id=%d, err=%s", typ, signID, e)
			continue
		}
		result = append(result, tmpResult...)
	}
	log.Info("get rank from db, sign id=%d, len=%d, date=%s", signID, len(result), date.Format(dateFmt))
	return
}

//ReloadRank reload rank from db
func (d *Dao) ReloadRank(signID int64) (err error) {
	// load cache
	v, err := d.loadRankUpFansCache(signID, time.Now())

	if err != nil {
		log.Error("load cache error, signID=%d, err=%s", signID, err)
		return
	}

	var key = cacheKeyRankFans(signID)
	d.localcache.SetWithExpire(key, v, time.Hour)
	log.Info("reload rank cache, sign id=%d", signID)
	return
}

// ----------------------------------------- load rank archive likes ------------------------------------
func (d *Dao) loadRankArchiveLikesCache(signID int64, date time.Time) (result *RankByTid, err error) {
	rawArchiveLike, err := d.RawRankArchiveLikes(signID, date)
	if err != nil {
		log.Error("fail to get raw rank up fans, signid=%d, err=%s", signID, err)
		return
	}

	result = new(RankByTid)
	if len(rawArchiveLike) == 0 {
		log.Info("archive likes rank data is empty, sign id=%d", signID)
		return
	}
	// 获取 aid列表
	var aids []int64
	//var accumulateMap = make(map[int64]*mcnmodel.McnRankArchiveLike)
	for _, v := range rawArchiveLike {
		aids = append(aids, v.ArchiveID)
		//if v.DataType == mcnmodel.DataTypeAccumulate {
		//	accumulateMap[v.ArchiveID] = v
		//}
	}

	aids = util.Unique(aids)
	// 获取archive信息

	arcsReply, err := global.GetArcGRPC().Arcs(context.Background(), &arcgrpc.ArcsRequest{Aids: aids})
	if err != nil {
		log.Error("fail to get archive info, sign_id=%d err=%s", signID, err)
		return
	}
	archiveDataMap := arcsReply.Arcs
	var tidUnique = newTidNameUnique()
	// 组装archive信息
	for _, v := range rawArchiveLike {
		var info = mcnmodel.RankArchiveLikeInfo{}
		info.CopyFromDB(v)
		var archive, ok = archiveDataMap[v.ArchiveID]
		if !ok {
			continue
		}
		info.CopyFromArchive(archive)

		tidUnique.addTid(info.Tid, info.TidName, info.DataType)

		//if accumulateData, ok := accumulateMap[v.ArchiveID]; ok {
		//	info.LikesAccumulate = accumulateData.LikeCount
		//}
		result.addRank(&info)
	}

	result.addTidMap(tidUnique)
	// 排序
	result.Sort(byValueDesc)
	result.Truncate(10)
	return
}

//RawRankArchiveLikes get from db
func (d *Dao) RawRankArchiveLikes(signID int64, date time.Time) (result []*mcnmodel.McnRankArchiveLike, err error) {
	// 有X种类型，
	// 昨日、上周、上月、累计
	// 每种类型取最近日期的数据
	for _, typ := range dataTypes {
		var tmpResult []*mcnmodel.McnRankArchiveLike
		e := d.mcndb.Raw(`select * from mcn_rank_archive_likes where data_type=? and sign_id=?
			and generate_date=(select generate_date from mcn_rank_archive_likes where data_type=? and sign_id=? and generate_date <= ? order by generate_date desc limit 1)`,
			typ, signID, typ, signID, date).
			Find(&tmpResult).Error
		if e != nil {
			log.Error("fail to get rank, type=%d, sign id=%d, err=%s", typ, signID, e)
			continue
		}
		result = append(result, tmpResult...)
	}
	log.Info("get rank from db, sign id=%d, len=%d, date=%s", signID, len(result), date.Format(dateFmt))
	return
}
