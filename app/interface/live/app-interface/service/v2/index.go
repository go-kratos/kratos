package v2

import (
	"context"
	"sync/atomic"
	"time"

	"go-common/app/interface/live/app-interface/dao/account"

	"go-common/app/interface/live/app-interface/dao/av"
	"go-common/app/interface/live/app-interface/dao/fans_medal"
	"go-common/app/interface/live/app-interface/dao/live_data"
	"go-common/app/interface/live/app-interface/dao/rankdb"
	"go-common/app/interface/live/app-interface/dao/relation"
	"go-common/app/interface/live/app-interface/dao/room_ex"
	"go-common/app/interface/live/app-interface/dao/user_ext"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/conf"
	"go-common/app/interface/live/app-interface/dao"
	liveuserDao "go-common/app/interface/live/app-interface/dao/live_user"
	roomDao "go-common/app/interface/live/app-interface/dao/room"
	roomexDao "go-common/app/interface/live/app-interface/dao/room_ex"
	xuserDao "go-common/app/interface/live/app-interface/dao/xuser"
	liveUserV1 "go-common/app/service/live/live_user/api/liverpc/v1"
	recommendV1 "go-common/app/service/live/recommend/api/grpc/v1"
	roomV2 "go-common/app/service/live/room/api/liverpc/v2"
	xrf "go-common/app/service/live/xroom-feed/api"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_bannerType           = 1
	_entranceType         = 2
	_yunyingRecFormType   = 3
	_yunyingRecSquareType = 4
	_rankType             = 5
	_recFormType          = 6
	_recSquareType        = 7
	_feedType             = 8
	_parentAreaFormType   = 9
	_parentAreaSquareType = 10
	_activityType         = 11
	_myAreaTagType        = 12
	_myAreaTagListType    = 13
	_seaPatrolType        = 14
)

// IndexService struct
type IndexService struct {
	conf       *conf.Config
	commonType []int64
	// dao
	dao          *dao.Dao
	roomDao      *roomDao.Dao
	liveuserDao  *liveuserDao.Dao
	rankdbDao    *rankdb.Dao
	livedataDao  *live_data.Dao
	relationDao  *relation.Dao
	roomexDao    *room_ex.Dao
	userextDao   *user_ext.Dao
	avDao        *av.Dao
	fansMedalDao *fans_medal.Dao
	accountDao   *account.Dao
	xuserDao     *xuserDao.Dao
	// cache
	// all base module cache
	//AllMInfoMap     map[int64][]*v2pb.ModuleInfo
	AllMInfoMap atomic.Value
	//areaEntranceMap map[int64][]*v2pb.PicItem
	areaEntranceListMap atomic.Value
	commonRoomList      atomic.Value
	recommendConn       recommendV1.RecommendClient
	xrfClient           *xrf.Client
}

// NewIndexService init
func NewIndexService(c *conf.Config) (s *IndexService) {
	s = &IndexService{
		conf:        c,
		dao:         dao.New(c),
		roomDao:     roomDao.New(c),
		liveuserDao: liveuserDao.New(c),
		rankdbDao:   rankdb.New(c),
		roomexDao:   roomexDao.New(c),
		accountDao:  account.New(c),
		xuserDao:    xuserDao.New(c),
		commonType: []int64{
			_yunyingRecFormType,
			_yunyingRecSquareType,
			_recFormType,
			_recSquareType,
			_parentAreaFormType,
			_parentAreaSquareType,
		},
	}
	// init cache data
	s.loadAllModuleInfoMap()
	s.loadAreaEntranceCache()
	s.loadCommonListMap()
	s.loadLastHourData(context.TODO())
	go s.allModuleInfoProc()
	go s.areaEntranceProc()
	go s.allcommonListProc()
	go s.loadLastHour()
	conn, err := recommendV1.NewClient(conf.Conf.Warden)
	if err != nil {
		panic(err)
	}
	s.recommendConn = conn
	xrfc, err := xrf.NewClient(conf.Conf.Warden)
	if err != nil {
		panic(err)
	}
	s.xrfClient = xrfc
	return s
}

// Index 相关服务

// GetAllList implementation
// 首页大接口
// `midware:"guest,verify"`
func (s *IndexService) GetAllList(ctx context.Context, req *v2pb.GetAllListReq) (resp *v2pb.GetAllListResp, err error) {
	resp = &v2pb.GetAllListResp{
		Interval:       10,
		IsSkyHorseGray: 0,
		Banner:         []*v2pb.MBanner{},
		MyTag:          []*v2pb.MMyTag{},
		AreaEntrance:   []*v2pb.MAreaEntrance{},
		SeaPatrol:      []*v2pb.MSeaPatrol{},
		MyIdol:         []*v2pb.MMyIdol{},
		RoomList:       []*v2pb.MRoomBlock{},
		HourRank:       []*v2pb.MHourRank{},
		ActivityCard:   []*v2pb.MActivityCard{},
	}
	moduleInfoMap := s.GetAllModuleInfoMapFromCache(ctx)
	if moduleInfoMap == nil {
		log.Error("[GetAllList]module info list is nil, moduleIds:%+v", moduleInfoMap)
		return
	}

	// 初始化各模块返回信息
	respCommonRoomList := make([]*v2pb.MRoomBlock, 0)
	respMultiRoomList := make([]*v2pb.MRoomBlock, 0)
	myAreaMap := make(map[int64]bool)
	respMyIdol := &v2pb.MMyIdol{}
	respSkyHorseRoomList := make([]*v2pb.CommonRoomItem, 0)
	respLiveRecRoomList := make([]*v2pb.CommonRoomItem, 0)

	// 大多使用header里的mid解析, 框架已封装请求的header
	midInterface, isUIDSet := metadata.Value(ctx, metadata.Mid).(int64)
	mid := int64(0)
	if isUIDSet {
		mid = midInterface
	}

	buvid := ""
	// 主站封好的，可从device里获取到sid、buvid、buvid3、build、channel、device、mobi_app、platform
	device, ok := metadata.Value(ctx, metadata.Device).(*blademaster.Device)
	if ok {
		buvid = device.Buvid
	}

	// 第一波并发获取数据，无依赖
	func(device *blademaster.Device) {
		wg1 := errgroup.Group{}

		// banner
		if s.isModuleExist(_bannerType) {
			wg1.Go(func() (err error) {
				resp.Banner = s.getIndexBanner(ctx, req.Platform, req.Device, req.Build)
				return
			})
		}
		// 常用标签列表
		if s.isModuleExist(_myAreaTagType) {
			wg1.Go(func() (err error) {
				resp.MyTag, _ = s.GetIndexV2TagList(ctx, &liveUserV1.UserSettingGetTagReq{})
				return
			})
		}
		// 分区入口
		if s.isModuleExist(_entranceType) {
			wg1.Go(func() (err error) {
				resp.AreaEntrance = s.getAreaEntrance(ctx)
				return
			})
		}
		// 我的关注
		if s.isModuleExist(_feedType) {
			wg1.Go(func() (err error) {
				resp.MyIdol = s.LiveAnchorHomePage(ctx, req.RelationPage, req.Build, req.Platform, req.Quality)
				if len(resp.MyIdol) > 0 {
					respMyIdol = resp.MyIdol[0]
				}
				return
			})
		}
		// 通用房间列表（肯定是有的，此处不做判断），推荐、运营推荐分区、常用分区（特殊:要在第二波拿）、一级分区
		wg1.Go(func() (err error) {
			respCommonRoomList = s.getCommonRoomListForIndex(ctx, req.Build, req.Platform, req.Quality)
			return
		})
		// 活动卡片
		if s.isModuleExist(_activityType) {
			wg1.Go(func() (err error) {
				resp.ActivityCard = s.getActivityCard(ctx)
				return
			})
		}
		// 小时榜
		if s.isModuleExist(_rankType) {
			wg1.Go(func() (err error) {
				resp.HourRank, _ = s.getLastHourTop3(ctx)
				return
			})
		}
		// 大航海
		mobiApp := device.RawMobiApp
		if s.isModuleExist(_seaPatrolType) && mobiApp != "iphone_b" && mobiApp != "android_b" {
			wg1.Go(func() (err error) {
				resp.SeaPatrol, _ = s.GetIndexV2SeaPatrol(ctx, &liveUserV1.NoteGetReq{})
				return
			})
		}

		err1 := wg1.Wait()
		if err1 != nil {
			log.Error("[GetAllList]wg1 wait error: %+v", err1)
		}
	}(device)

	// 第二波获取数据 (依赖第一波)
	func() {
		wg2 := errgroup.Group{}
		// 天马个性化推荐 无法缓存 依赖关注
		if s.ifHitSkyHorse(mid, req.Device) {
			wg2.Go(func() (err error) {
				//目前只对第一个关注模块去重
				respSkyHorseRoomList, err = s.getSkyHorseRoomListForIndex(ctx, respMyIdol, mid, buvid, req.Build, req.Platform, req.RecPage, req.Quality)
				return
			})
		}
		// 直播个性化推荐 无法缓存 依赖关注
		if s.ifHitLiveRec(mid, req.Device) {
			wg2.Go(func() (err error) {
				respLiveRecRoomList, err = s.getLiveRecRoomList(ctx, respMyIdol, mid, req.Build, req.Platform, req.RecPage, req.Quality)
				return
			})
		}
		// 常用标签房间列表
		wg2.Go(func() (err error) {
			respMultiRoomList, myAreaMap = s.getMultiRoomList(ctx, resp.MyTag, req.Platform, req.Build, req.Quality)
			resp.RoomList = append(resp.RoomList, respMultiRoomList...)
			return
		})

		// 对保底逻辑的一些处理，对关注去重, 数量限制, 获取投放系统数据
		wg2.Go(func() (err error) {
			respCommonRoomList = s.handleCommonRoomList(ctx, respMyIdol, respCommonRoomList, req.Quality, req.Build, req.Platform, req.Device)
			return
		})

		err2 := wg2.Wait()
		if err2 != nil {
			log.Error("[GetAllList]wg2 wait error: %+v", err2)
		}
	}()

	// 推荐直播最终处理
	handleRecResult(resp, respCommonRoomList, respSkyHorseRoomList, respLiveRecRoomList, myAreaMap)
	if resp.IsSkyHorseGray == 0 && s.ifHitLiveRec(mid, req.Device) {
		log.Info("live rec hit miss, mid:%d, liveRec:%+v", mid, respLiveRecRoomList)
	}
	return
}

// Change implementation
// 换一换接口
// `midware:"guest,verify"`
func (s *IndexService) Change(ctx context.Context, req *v2pb.ChangeReq) (resp *v2pb.ChangeResp, err error) {
	resp = &v2pb.ChangeResp{}
	mid, isUIDSet := metadata.Value(ctx, metadata.Mid).(int64)
	var uid int64
	if isUIDSet {
		uid = mid
	}

	duplicates, _ := xstr.SplitInts(req.AttentionRoomId)
	duplicatesMap := make(map[int64]bool)

	for _, roomID := range duplicates {
		duplicatesMap[roomID] = true
	}
	buvid := ""
	// 主站封好的，可从device里获取到sid、buvid、buvid3、build、channel、device、mobi_app、platform
	device, ok := metadata.Value(ctx, metadata.Device).(*blademaster.Device)
	if ok {
		buvid = device.Buvid
	}

	moduleInfo, err := s.roomDao.GetAllModuleInfo(ctx, req.ModuleId)
	if err != nil || moduleInfo[0] == nil {
		log.Error("[Change]GetModuleInfoById error:%+v", err)
		return
	}

	// 给moduleInfo赋值
	resp.ModuleInfo = &v2pb.ModuleInfo{
		Id:    moduleInfo[0].Id,
		Link:  moduleInfo[0].Link,
		Pic:   moduleInfo[0].Pic,
		Title: moduleInfo[0].Title,
		Type:  moduleInfo[0].Type,
		Sort:  moduleInfo[0].Sort,
		Count: moduleInfo[0].Count,
	}
	resp.List = make([]*v2pb.CommonRoomItem, 0)

	isDefault := true
	if s.ifHitSkyHorse(uid, req.Device) {
		skyHorseList, err := s.getSkyHorseRoomList(ctx, mid, buvid, req.Build, req.Platform, duplicates, req.Page, req.Quality)
		if err == nil && len(skyHorseList) > 0 {
			isDefault = false
			resp.IsSkyHorseGray = 1
			resp.List = skyHorseList
		}
	}

	if s.ifHitLiveRec(mid, req.Device) {
		respLiveRoomList, err := s.getLiveRecRoomListForChange(ctx, mid, req.Build, req.Platform, duplicates, req.Page, req.Quality)
		if err == nil && len(respLiveRoomList) > 0 {
			isDefault = false
			resp.IsSkyHorseGray = 1
			resp.List = respLiveRoomList
		} else {
			log.Info("live rec hit miss, from:change, mid:%d, err:%+v, liveRec:%+v", mid, err, respLiveRoomList)
		}
	}

	if isDefault {
		resp.IsSkyHorseGray = 0
		respCommonRoomList, errTemp := s.getCommonRoomListByID(ctx, req.ModuleId, req.Build, req.Platform, req.Quality, req.Device, duplicates)
		if errTemp != nil {
			log.Error("[Change]GetModuleInfoById error:%+v", errTemp)
			err = errTemp
			return
		}

		resp.List = respCommonRoomList
	}

	return
}

// 指定模块是否存在
func (s *IndexService) isModuleExist(iType int64) (res bool) {
	res = false
	mInfoMap := s.GetAllModuleInfoMapFromCache(context.TODO())
	if _, ok := mInfoMap[iType]; ok {
		res = true
	}
	return
}

// 推荐模块最终处理：天马、对关注去重
func handleRecResult(resp *v2pb.GetAllListResp, respCommonRoomList []*v2pb.MRoomBlock, respSkyHorseRoomList []*v2pb.CommonRoomItem, respLiveRecRoomList []*v2pb.CommonRoomItem, myAreaMap map[int64]bool) {
	afterHandleRoomList := make([]*v2pb.MRoomBlock, 0)
	for _, roomBlock := range respCommonRoomList {
		if roomBlock.ModuleInfo.Type == _recFormType || roomBlock.ModuleInfo.Type == _recSquareType {
			if len(respSkyHorseRoomList) > 0 {
				resp.IsSkyHorseGray = 1
				roomBlock.List = respSkyHorseRoomList
			} else if len(respLiveRecRoomList) > 0 {
				resp.IsSkyHorseGray = 1
				roomBlock.List = respLiveRecRoomList
			}

			afterHandleRoomList = append(afterHandleRoomList, roomBlock)
			continue
		}
		// 常用分区对运营推荐分区去重
		if roomBlock.ModuleInfo.Type == _yunyingRecFormType || roomBlock.ModuleInfo.Type == _yunyingRecSquareType {
			for _, item := range roomBlock.List {
				if _, ok := myAreaMap[item.AreaV2Id]; !ok {
					afterHandleRoomList = append(afterHandleRoomList, roomBlock)
					break
				}
			}
			continue
		}
		afterHandleRoomList = append(afterHandleRoomList, roomBlock)
	}
	resp.RoomList = append(resp.RoomList, afterHandleRoomList...)
}

// get AllModuleInfoMap
func (s *IndexService) getAllModuleInfoMap(ctx context.Context) (allMInfoCacheMap map[int64][]*v2pb.ModuleInfo) {
	var allMInfoData []*roomV2.AppIndexGetBaseMInfoListResp_ModuleInfo
	var err error
	var retry int64
	for i := 1; i <= 3; i++ {
		// 最多重试3次
		allMInfoData, err = s.roomDao.GetAllModuleInfo(ctx, 0)
		if err != nil || len(allMInfoData) <= 0 {
			retry++
			log.Error("[loadAllModuleInfoMap] GetAllModuleInfo error(%+v) retry_times(%d)", err, retry)
			continue
		}
		break
	}

	if len(allMInfoData) > 0 && allMInfoData[1] != nil {
		allMInfoCacheMap = make(map[int64][]*v2pb.ModuleInfo)
		for _, m := range allMInfoData {
			allMInfoCacheMap[m.Type] = append(allMInfoCacheMap[m.Type], &v2pb.ModuleInfo{
				Id:    m.Id,
				Link:  m.Link,
				Pic:   m.Pic,
				Title: m.Title,
				Type:  m.Type,
				Sort:  m.Sort,
				Count: m.Count,
			})
		}
	}
	return
}

// cache load
func (s *IndexService) loadAllModuleInfoMap() {
	allMInfoCacheMap := s.getAllModuleInfoMap(context.TODO())
	if len(allMInfoCacheMap) > 0 {
		s.AllMInfoMap.Store(allMInfoCacheMap)
		log.Info("[loadAllModuleInfoMap]load data success!")
	}
}

// ticker
func (s *IndexService) allModuleInfoProc() {
	for {
		time.Sleep(time.Second * 5)
		s.loadAllModuleInfoMap()
	}
}

// GetAllModuleInfoMapFromCache get all module info fromcache
func (s *IndexService) GetAllModuleInfoMapFromCache(ctx context.Context) (res map[int64][]*v2pb.ModuleInfo) {
	// load
	i := s.AllMInfoMap.Load()
	// assert
	res, ok := i.(map[int64][]*v2pb.ModuleInfo)
	if ok {
		return
	}
	// 回源&log
	res = s.getAllModuleInfoMap(ctx)
	log.Warn("[GetAllModuleInfoMap]memory cache miss!! i:%+v; res:%+v", i, res)
	return
}

func (s *IndexService) getCommonListFromCache(sIds []int64) (commonList map[int64]*roomV2.AppIndexGetRoomListByIdsResp_RoomList) {
	commonListCache, ok := s.commonRoomList.Load().(map[int64]*roomV2.AppIndexGetRoomListByIdsResp_RoomList)
	if ok {
		commonList = commonListCache
		return
	}

	roomListMap, err := s.roomDao.GetListByIds(context.TODO(), sIds)
	if err != nil {
		log.Error("[getCommonListFromCache]GetListByIds error, error:%+v", err)
		return
	}
	commonList = roomListMap
	log.Warn("[getCommonListFromCache]memory cache miss!! res:%+v", commonList)
	return
}

func (s *IndexService) loadCommonListMap() {
	sIds := s.getIdsFromModuleMap(context.TODO(), s.commonType)

	roomListMap, err := s.roomDao.GetListByIds(context.TODO(), sIds)
	if err != nil {
		log.Error("[loadCommonListMap]GetListByIds error, error:%+v", err)
		return
	}

	copyRoomListMap := make(map[int64]*roomV2.AppIndexGetRoomListByIdsResp_RoomList)

	for moduleId, roomList := range roomListMap {
		for _, item := range roomList.List {
			if _, ok := copyRoomListMap[moduleId]; !ok {
				copyRoomListMap[moduleId] = &roomV2.AppIndexGetRoomListByIdsResp_RoomList{
					List: make([]*roomV2.AppIndexGetRoomListByIdsResp_RoomInfo, 0),
				}
			}
			copyRoomListMap[moduleId].List = append(copyRoomListMap[moduleId].List, item)
		}
	}
	s.commonRoomList.Store(copyRoomListMap)
}

func (s *IndexService) allcommonListProc() {
	for {
		time.Sleep(time.Second * 2)
		s.loadCommonListMap()
	}
}

// 根据type从模块信息map拿到模块ids列表
func (s *IndexService) getIdsFromModuleMap(ctx context.Context, iTypes []int64) (sIds []int64) {
	mMap := s.GetAllModuleInfoMapFromCache(ctx)
	sIds = make([]int64, 0)
	for _, iType := range iTypes {
		typeList, ok := mMap[iType]
		if !ok {
			continue
		}
		for _, item := range typeList {
			if item != nil {
				sIds = append(sIds, item.Id)
			}
		}
	}

	return
}
