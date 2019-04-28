package v2

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	bp "go-common/app/interface/live/app-interface/api/http/v2"
	"go-common/app/interface/live/app-interface/dao"
	"go-common/app/service/live/room/api/liverpc/v1"
	accountM "go-common/app/service/main/account/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

type lastHourCache struct {
	Modules []*bp.MHourRank
	CTime   int64
}

const (
	_hourRankType = 5
)

var LastHourItemCache atomic.Value

// 获取首页上小时排行榜
func (s *IndexService) getLastHourTop3(ctx context.Context) (resp []*bp.MHourRank, err error) {
	moduleInfoMaps := s.GetAllModuleInfoMapFromCache(ctx)
	resp = []*bp.MHourRank{}
	module, ok := moduleInfoMaps[_hourRankType]
	if !ok || 0 == len(module) {
		return
	}
	cacheResp, err := getLastHour3FromCache(ctx)

	if 0 != len(cacheResp) && 0 != len(cacheResp[0].List) {
		resp = cacheResp
		return
	}
	// load
	resp, err = s.loadLastHourData(ctx)
	return
}

func setLastHour3Cache(c context.Context, modules []*bp.MHourRank) {
	cache := &lastHourCache{
		Modules: modules,
		CTime:   time.Now().Unix(),
	}
	LastHourItemCache.Store(cache)
	return
}

func getLastHour3FromCache(c context.Context) (resp []*bp.MHourRank, err error) {
	resp = make([]*bp.MHourRank, 0)
	now := time.Now().Unix()

	resp = []*bp.MHourRank{}
	cache, ok := LastHourItemCache.Load().(*lastHourCache)
	if !ok || nil == cache {
		return
	}
	cTime := cache.CTime
	if now > (cTime + 60) {
		return
	}
	resp = cache.Modules
	return
}

// loadLastHour 定时存cache
func (s *IndexService) loadLastHour() {
	for {
		time.Sleep(time.Second * 20)
		s.loadLastHourData(context.TODO())
	}
}

// loadLastHourData 获取上小时榜数据
func (s *IndexService) loadLastHourData(ctx context.Context) (resp []*bp.MHourRank, err error) {
	moduleInfoMaps := s.GetAllModuleInfoMapFromCache(ctx)
	resp = []*bp.MHourRank{}
	module, ok := moduleInfoMaps[_hourRankType]
	if !ok || 0 == len(module) {
		return
	}
	// extraInfo format
	var nowHourName string
	lastHourName := time.Now().Add(-time.Hour).Format("15")
	timeNum, _ := strconv.Atoi(lastHourName)
	lastHourName += ":00"
	if timeNum > 9 {
		nowHourName = strconv.Itoa(timeNum + 1)
	} else {
		nowHourName = time.Now().Format("15")
	}
	nowHourName += ":00"
	subTitle := lastHourName + "-" + nowHourName + " 总榜排名"

	list := []*bp.HourRankItem{}
	roomReq := &v1.RoomGetStatusInfoByUidsReq{}
	roomResp := &v1.RoomGetStatusInfoByUidsResp{}
	userInfo := map[int64]*accountM.Card{}
	wg := &errgroup.Group{}
	// liveRpc call rankdb
	uids, err := s.rankdbDao.GetLastHourTop3(ctx)
	if 0 == len(uids) || nil != err {
		goto formatReturn
	}

	wg, _ = errgroup.WithContext(ctx)
	wg.Go(func() error {
		// liveRpc call room
		roomReq.Uids = uids
		roomResp, err = dao.RoomApi.V1Room.GetStatusInfoByUids(ctx, roomReq)
		return err
	})
	wg.Go(func() error {
		// call account for UserInfo
		userInfo, err = s.rankdbDao.GetUserInfoData(ctx, uids)
		return err
	})

	if err = wg.Wait(); nil != err {
		goto formatReturn
	}

	if 0 != roomResp.Code || 0 == len(roomResp.Data) {
		log.Error("[app-interface][rankDbItem] liveRpc call room return error, code:%d, msg:%s", roomResp.Code, roomResp.Data)
		goto formatReturn
	}

	if 0 == len(userInfo) {
		log.Error("[app-interface][rankDbItem]  call account return empty")
		goto formatReturn
	}

	for k, v := range uids {
		detail := &bp.HourRankItem{}
		if nil != roomResp.Data[v] {
			detail = &bp.HourRankItem{
				Roomid:           roomResp.Data[v].RoomId,
				LiveStatus:       roomResp.Data[v].LiveStatus,
				AreaV2ParentId:   roomResp.Data[v].AreaV2ParentId,
				AreaV2Id:         roomResp.Data[v].AreaV2Id,
				AreaV2ParentName: roomResp.Data[v].AreaV2ParentName,
				AreaV2Name:       roomResp.Data[v].AreaV2Name,
				Uname:            userInfo[v].Name,
				Face:             userInfo[v].Face,
			}
		}
		if nil != userInfo[v] {
			detail.Uname = userInfo[v].Name
			detail.Face = userInfo[v].Face
		}
		detail.Rank = int64(k + 1)
		detail.Uid = v

		list = append(list, detail)
	}

	// format return
formatReturn:
	for _, v := range module {
		if v.Type == _hourRankType {
			item := &bp.MHourRank{
				ModuleInfo: v,
				ExtraInfo:  &bp.HourRankExtra{SubTitle: subTitle},
				List:       list,
			}
			resp = append(resp, item)
			break
		}
	}
	moduleCache := resp
	setLastHour3Cache(ctx, moduleCache)
	return
}
