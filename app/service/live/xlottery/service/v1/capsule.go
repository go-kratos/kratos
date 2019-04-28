package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"go-common/library/ecode"
	"go-common/library/queue/databus/report"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/live/xlottery/model"

	"go-common/app/service/live/gift/api/liverpc/v1"
	v12 "go-common/app/service/live/rc/api/liverpc/v1"
	v13 "go-common/app/service/live/userext/api/liverpc/v1"
	v1pb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/app/service/live/xlottery/conf"
	"go-common/app/service/live/xlottery/dao"
	actmdl "go-common/app/service/main/account/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"
)

// CapsuleService struct
type CapsuleService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao        *dao.Dao
	httpClient *bm.Client
	runCache   *fanout.Fanout
}

//NewCapsuleService init
func NewCapsuleService(c *conf.Config) (s *CapsuleService) {
	s = &CapsuleService{
		conf:       c,
		dao:        dao.New(c),
		httpClient: bm.NewClient(c.HTTPClient),
		runCache:   fanout.New("capsule"),
	}
	report.InitUser(conf.Conf.UserReport)
	s.tickerReloadCapsuleConf(context.TODO())
	return s
}

// GetCoinList implementation
// 获取扭蛋币列表
func (s *CapsuleService) GetCoinList(ctx context.Context, req *v1pb.GetCoinListReq) (resp *v1pb.GetCoinListResp, err error) {
	resp = &v1pb.GetCoinListResp{}
	var CoinRet []*model.Coin
	if CoinRet, err = s.dao.GetAllCoin(ctx); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	page := req.Page
	pageSize := req.PageSize
	totalNum := len(CoinRet)
	totalPage := math.Ceil(float64(totalNum) / float64(pageSize))
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > int64(totalNum) {
		end = int64(totalNum)
	}
	List := make([]*v1pb.GetCoinListResp_List, 0)
	coinList := CoinRet[start:end]
	for _, coinInfo := range coinList {
		var allCoinConfig []*model.CoinConfig
		if allCoinConfig, err = s.dao.GetCoinConfig(ctx, coinInfo.Id); err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
		giftConfigs := make([]string, 0)
		areaConfigs := make([][2]int64, 0)
		for _, config := range allCoinConfig {
			if config.Type == 2 {
				giftConfigs = append(giftConfigs, strconv.FormatInt(config.GiftId, 10))
			} else if config.Type == 1 {
				areaConfigs = append(areaConfigs, [2]int64{config.AreaV2ParentId, config.AreaV2Id})
			}
		}
		AreaMap := make(map[int64]*v1pb.GetCoinListResp_List_AreaIds)
		for _, areaConfig := range areaConfigs {
			parentAreaId := areaConfig[0]
			if _, ok := AreaMap[parentAreaId]; !ok {
				AreaMap[parentAreaId] = &v1pb.GetCoinListResp_List_AreaIds{}
				AreaMap[parentAreaId].List = make([]int64, 0)
			}
			AreaMap[parentAreaId].ParentId = parentAreaId
			if parentAreaId == areaConfig[1] {
				AreaMap[parentAreaId].IsAll = 1
				AreaMap[parentAreaId].List = make([]int64, 0)
				continue
			}
			AreaMap[parentAreaId].List = append(AreaMap[parentAreaId].List, areaConfig[1])
		}
		ix := 0
		Areas := make([]*v1pb.GetCoinListResp_List_AreaIds, len(AreaMap))
		for _, area := range AreaMap {
			Areas[ix] = area
			ix++
		}
		Singlelist := &v1pb.GetCoinListResp_List{}
		Singlelist.Id = coinInfo.Id
		Singlelist.Title = coinInfo.Title
		Singlelist.GiftType = coinInfo.GiftType
		Singlelist.ChangeNum = coinInfo.ChangeNum
		Singlelist.StartTime = coinInfo.StartTime
		Singlelist.EndTime = coinInfo.EndTime
		Singlelist.Status = coinInfo.Status
		Singlelist.GiftConfig = strings.Join(giftConfigs, ",")
		Singlelist.AreaIds = Areas
		List = append(List, Singlelist)
	}
	resp = &v1pb.GetCoinListResp{
		Total:     int64(totalNum),
		TotalPage: int64(totalPage),
		List:      List,
	}
	if err != nil {
		return
	}
	return
}

// UpdateCoinConfig implementation
// 更新扭蛋币配置
func (s *CapsuleService) UpdateCoinConfig(ctx context.Context, req *v1pb.UpdateCoinConfigReq) (resp *v1pb.UpdateCoinConfigResp, err error) {
	resp = &v1pb.UpdateCoinConfigResp{}
	var status bool
	if req.Id == dao.NormalCoinId {
		if req.GiftType != dao.CapsuleGiftTypeAll {
			err = ecode.Error(ecode.XLotteryCapsuleCoinNotChange, "普通扭蛋币不可更改道具类型")
			return
		}
	}
	status, err = s.dao.UpdateCoin(ctx, req)
	if err != nil || !status {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	if req.Id != dao.NormalCoinId {
		if status, err = s.dao.UpdateCoinConfig(ctx, req.Id, req.AreaIds, req.GiftIds); err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
	}
	resp.Status = status
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | UpdateCoinConfig] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// UpdateCoinStatus implementation
// 更新扭蛋币状态
func (s *CapsuleService) UpdateCoinStatus(ctx context.Context, req *v1pb.UpdateCoinStatusReq) (resp *v1pb.UpdateCoinStatusResp, err error) {
	resp = &v1pb.UpdateCoinStatusResp{}
	var status bool
	if status, err = s.dao.UpdateCoinStatus(ctx, req.Id, req.Status); err != nil {
		return
	}
	resp.Status = status
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | UpdateCoinStatus] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// DeleteCoin implementation
// 删除扭蛋币
func (s *CapsuleService) DeleteCoin(ctx context.Context, req *v1pb.DeleteCoinReq) (resp *v1pb.DeleteCoinResp, err error) {
	resp = &v1pb.DeleteCoinResp{}
	var status bool
	if status, err = s.dao.DeleteCoin(ctx, req.Id); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	resp.Status = status
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | DeleteCoin] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// GetPoolList implementation
func (s *CapsuleService) GetPoolList(ctx context.Context, req *v1pb.GetPoolListReq) (resp *v1pb.GetPoolListResp, err error) {
	resp = &v1pb.GetPoolListResp{}
	var poolRet []*model.Pool
	poolList := make([]*v1pb.GetPoolListResp_List, 0)
	if poolRet, err = s.dao.GetAllPool(ctx); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	coinMap, err := s.dao.GetCoinMap(ctx)
	if err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	page := req.Page
	pageSize := req.PageSize
	totalNum := len(poolRet)
	totalPage := math.Ceil(float64(totalNum) / float64(pageSize))
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > int64(totalNum) {
		end = int64(totalNum)
	}
	for _, single := range poolRet[start:end] {
		singlePool := &v1pb.GetPoolListResp_List{}
		singlePool.Id = single.Id
		singlePool.CoinId = single.CoinId
		singlePool.Title = single.Title
		if _, ok := coinMap[single.CoinId]; ok {
			singlePool.CoinTitle = coinMap[single.CoinId].Title
		}
		singlePool.StartTime = single.StartTime
		singlePool.EndTime = single.EndTime
		singlePool.Status = single.Status
		singlePool.Rule = single.Description
		poolList = append(poolList, singlePool)
	}
	resp = &v1pb.GetPoolListResp{
		Total:     int64(totalNum),
		TotalPage: int64(totalPage),
		List:      poolList,
	}
	return
}

// UpdatePool implementation
func (s *CapsuleService) UpdatePool(ctx context.Context, req *v1pb.UpdatePoolReq) (resp *v1pb.UpdatePoolResp, err error) {
	resp = &v1pb.UpdatePoolResp{}
	var status bool
	if req.Id == 0 {
		if status, err = s.dao.CreatePool(ctx, req); err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
	} else {
		var poolInfo *model.Pool
		poolInfo, err = s.dao.GetPoolById(ctx, req.Id)
		if err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
		if (poolInfo.StartTime != req.StartTime || poolInfo.EndTime != req.EndTime) && poolInfo.IsBottom == 1 {
			ecode.Error(ecode.XLotteryCapsulePoolNotChange, "系统繁忙，保底奖池的上下线时间不可修改")
			return
		}
		if poolInfo.CoinId != req.CoinId || poolInfo.Title != req.Title || poolInfo.Description != req.Rule || poolInfo.StartTime != req.StartTime || poolInfo.EndTime != req.EndTime {
			if status, err = s.dao.UpdatePool(ctx, req); err != nil {
				err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
				return
			}
		} else {
			status = true
		}

	}
	resp.Status = status
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | UpdatePool] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// DeletePool implementation
func (s *CapsuleService) DeletePool(ctx context.Context, req *v1pb.DeletePoolReq) (resp *v1pb.DeletePoolResp, err error) {
	resp = &v1pb.DeletePoolResp{}
	var rows int64
	id := req.Id
	if rows, err = s.dao.DelPool(ctx, id); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	resp.Status = rows > 0
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[servive.capsule | DeletePool] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// UpdatePoolStatus implementation
func (s *CapsuleService) UpdatePoolStatus(ctx context.Context, req *v1pb.UpdatePoolStatusReq) (resp *v1pb.UpdatePoolStatusResp, err error) {
	resp = &v1pb.UpdatePoolStatusResp{}
	poolInfo, err := s.dao.GetPoolById(ctx, req.Id)
	if err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	if poolInfo == nil {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return
	}

	if req.Status == 1 {
		var poolPrizeList []*model.PoolPrize
		poolPrizeList, err = s.dao.GetPoolPrize(ctx, req.Id)
		if err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return resp, err
		}
		var totalChance int64
		for _, prize := range poolPrizeList {
			totalChance += prize.Chance
		}
		if totalChance != 10000 {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，概率设置错误")
			return resp, err
		}
	} else {
		if poolInfo.IsBottom == 1 {
			err = ecode.Error(ecode.XLotteryCapsulePoolNotOffline, "保底奖池不可下线")
			return
		}
	}
	var status bool
	if status, err = s.dao.UpdatePoolStatus(ctx, req.Id, req.Status); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	resp.Status = status
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | UpdatePoolStatus] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// GetPoolPrize implementation
func (s *CapsuleService) GetPoolPrize(ctx context.Context, req *v1pb.GetPoolPrizeReq) (resp *v1pb.GetPoolPrizeResp, err error) {
	resp = &v1pb.GetPoolPrizeResp{}
	var poolPrizeList []*model.PoolPrize
	if poolPrizeList, err = s.dao.GetPoolPrize(ctx, req.PoolId); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	titleMap := make(map[int64]string)
	whitePoolIds := make([]int64, 0)
	list := make([]*v1pb.GetPoolPrizeResp_List, len(poolPrizeList))
	for _, prize := range poolPrizeList {
		if prize.Type == dao.CapsulePrizeTitleType {
			titleMap[prize.ObjectId] = ""
		}
		if prize.ProType == dao.ProTypeWhite {
			whitePoolIds = append(whitePoolIds, prize.Id)
		}
	}
	whiteMap := make(map[int64][]int64)
	if len(whitePoolIds) > 0 {
		whiteMap, err = s.dao.GetWhiteUserMap(ctx, whitePoolIds)
		if err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
	}
	tlen := len(titleMap)
	if tlen > 0 {
		titleIds := make([]int64, tlen)
		ix := 0
		for tltleId := range titleMap {
			titleIds[ix] = tltleId
			ix++
		}
		TitleData, err1 := dao.RcApi.V1UserTitle.GetTitleByIds(ctx, &v12.UserTitleGetTitleByIdsReq{Ids: titleIds})
		if err1 != nil {
			log.Error("[servive.capsule | GetPoolPrize]  GetTitleByIds error")
		}
		if TitleData != nil && TitleData.Data != nil {
			titleMap = TitleData.Data
		}
	}
	for ix, prize := range poolPrizeList {
		prizeItem := &v1pb.GetPoolPrizeResp_List{}
		prizeItem.Id = prize.Id
		prizeItem.PoolId = prize.PoolId
		prizeItem.Type = prize.Type
		prizeItem.Num = prize.Num
		prizeItem.ObjectId = prize.ObjectId
		prizeItem.Expire = prize.Expire
		prizeItem.WebUrl = prize.WebUrl
		prizeItem.MobileUrl = prize.MobileUrl
		prizeItem.Description = prize.Description
		prizeItem.JumpUrl = prize.JumpUrl
		prizeItem.ProType = prize.ProType
		var chance string
		if prize.ProType == dao.ProTypeNormal {
			chance = strconv.FormatInt(prize.Chance, 10)
			clen := len(chance)
			if clen >= 5 {
				chance = "1"
			} else if clen == 4 {
				chance = "0." + strings.Trim(chance, "0")
			} else if clen < 4 && clen >= 1 {
				if prize.Chance > 0 {
					chance = "0." + strings.Repeat("0", 4-clen) + strings.Trim(chance, "0")
				} else {
					chance = "0"
				}
			} else {
				chance = "0"
			}
		} else if prize.ProType == dao.ProTypeFixDay {
			chance = "每天" + strconv.FormatInt(prize.LimitNum, 10) + "个"
		} else if prize.ProType == dao.ProTypeFixWeek {
			chance = "每周" + strconv.FormatInt(prize.LimitNum, 10) + "个"
		} else if prize.ProType == dao.ProTypeWhite {
			chance = "白名单"
		}

		prizeItem.Chance = chance
		prizeItem.Loop = prize.LoopNum
		prizeItem.Limit = prize.LimitNum
		prizeItem.Weight = prize.Weight
		if prize.ProType == dao.ProTypeWhite {
			if _, ok := whiteMap[prize.Id]; ok {
				prizeItem.WhiteUids = whiteMap[prize.Id]
			}
		}
		if prize.Type == dao.CapsulePrizeTitleType {
			prizeItem.Name = "头衔"
			if _, ok := titleMap[prize.ObjectId]; ok {
				prizeItem.Name = titleMap[prize.ObjectId] + "头衔"
			}
		} else {
			prizeItem.Name = ""
			if _, ok := dao.PrizeNameMap[prize.Type]; ok {
				prizeItem.Name = dao.PrizeNameMap[prize.Type]
			}
		}
		list[ix] = prizeItem
	}
	resp.List = list
	return
}

// UpdatePoolPrize implementation
func (s *CapsuleService) UpdatePoolPrize(ctx context.Context, req *v1pb.UpdatePoolPrizeReq) (resp *v1pb.UpdatePoolPrizeResp, err error) {
	resp = &v1pb.UpdatePoolPrizeResp{}
	var status bool
	var prizeId int64
	if req.Type != dao.CapsulePrizeTitleType {
		req.ObjectId = 0
	}
	if req.Id == 0 {
		if prizeId, err = s.dao.AddPoolPrize(ctx, req); err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
		status = prizeId > 0
	} else {
		if status, err = s.dao.UpdatePoolPrize(ctx, req); err != nil {
			err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
			return
		}
		prizeId = req.Id
	}
	if status {
		if req.ProType == dao.ProTypeWhite {
			status, err = s.dao.AddWhiteUser(ctx, prizeId, req.WhiteUids)
			if err != nil {
				err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
				return
			}
		}
	}
	resp.Status = status
	resp.PrizeId = prizeId

	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | UpdatePoolPrize] SetCapsuleChangeFlag error : %v", err1)
	}
	return
}

// DeletePoolPrize implementation
func (s *CapsuleService) DeletePoolPrize(ctx context.Context, req *v1pb.DeletePoolPrizeReq) (resp *v1pb.DeletePoolPrizeResp, err error) {
	resp = &v1pb.DeletePoolPrizeResp{}
	var status bool
	if status, err = s.dao.DeletePoolPrize(ctx, req.Id); err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	resp.Status = status
	_, err1 := s.dao.SetCapsuleChangeFlag(ctx)
	if err1 != nil {
		log.Error("[service.capsule | DeletePoolPrize] SetCapsuleChangeFlag error : %v", err1)
		return
	}
	return
}

// GetDetail implementation
// * 获取扭蛋详情
// 需要登录态; 会清除红点
func (s *CapsuleService) GetDetail(ctx context.Context, req *v1pb.CapsuleGetDetailReq) (resp *v1pb.CapsuleGetDetailResp, err error) {
	resp = &v1pb.CapsuleGetDetailResp{}
	userInfo, err := s.dao.GetUserCapsuleInfo(ctx, req.Uid)
	if err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil || len(coinConfMap) == 0 {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	resp.Normal = new(v1pb.CapsuleGetDetailResp_CapsuleInfo)
	resp.Colorful = new(v1pb.CapsuleGetDetailResp_CapsuleInfo)
	for _, coinConf := range coinConfMap {
		// 普通扭蛋
		if coinConf.Id == dao.NormalCoinId {
			resp.Normal.Status = true
			resp.Normal.Coin = s.dao.GetCoin(userInfo[dao.NormalCoinId], coinConf)
			resp.Normal.Change = s.dao.GetChangeNum(ctx, req.Uid, dao.NormalCoinId)
			resp.Normal.Progress = s.dao.GetProgress(userInfo[dao.NormalCoinId], coinConf)
			if req.From != dao.GetCapsuleDetailFromRoom && coinConf.PoolConf != nil {
				resp.Normal.Rule = coinConf.PoolConf.Rule
				resp.Normal.Gift = s.getGift(ctx, coinConf.PoolConf)
			}
			if req.From == dao.GetCapsuleDetailFromWeb {
				resp.Normal.List = s.getOpenHistory(ctx, dao.NormalCoinId)
			}
		} else if coinConf.Id == dao.ColorfulCoinId {
			resp.Colorful.Status = s.dao.IsPoolOpen(coinConf, dao.ColorfulCoinId)
			if !resp.Colorful.Status {
				continue
			}
			resp.Colorful.Coin = s.dao.GetCoin(userInfo[dao.ColorfulCoinId], coinConf)
			resp.Colorful.Change = s.dao.GetChangeNum(ctx, req.Uid, dao.ColorfulCoinId)
			resp.Colorful.Progress = s.dao.GetProgress(userInfo[dao.ColorfulCoinId], coinConf)
			if req.From != dao.GetCapsuleDetailFromRoom && coinConf.PoolConf != nil {
				resp.Colorful.Rule = coinConf.PoolConf.Rule
				resp.Colorful.Gift = s.getGift(ctx, coinConf.PoolConf)
			}
			if req.From == dao.GetCapsuleDetailFromWeb {
				resp.Colorful.List = s.getOpenHistory(ctx, dao.ColorfulCoinId)
			}
		}

	}
	s.dao.ClearNoticeBoth(ctx, req.Uid)
	return
}

func (s *CapsuleService) getGift(ctx context.Context, poolConf *dao.CapsulePoolConf) (gifts []*v1pb.CapsuleGetDetailResp_Gift) {
	gifts = make([]*v1pb.CapsuleGetDetailResp_Gift, 0)
	if poolConf == nil {
		return
	}
	poolPrize := poolConf.PoolPrize
	if len(poolPrize) == 0 {
		return
	}
	giftMap := make(map[int64]struct{})
	for _, prize := range poolPrize {
		gType := prize.Type*1000000 + prize.ObjectId
		log.Info("gift %d %d", prize.Id, gType)
		if _, ok := giftMap[gType]; ok {
			continue
		}
		giftMap[gType] = struct{}{}
		gift := new(v1pb.CapsuleGetDetailResp_Gift)
		gift.Image = prize.MobileImage
		gift.MobileImage = prize.MobileImage
		gift.WebImage = prize.WebImage
		gift.Name = prize.Name
		gift.Usage = new(v1pb.Usage)
		gift.Usage.Text = prize.Description
		gift.Usage.Url = prize.JumpUrl
		gifts = append(gifts, gift)
	}
	return
}

func (s *CapsuleService) getOpenHistory(ctx context.Context, coinType int64) (list []*v1pb.CapsuleGetDetailResp_List) {
	list = make([]*v1pb.CapsuleGetDetailResp_List, 0)
	openData, err := s.dao.GetOpenHistory(ctx, coinType)
	if err != nil {
		return
	}
	length := len(openData)
	if length == 0 {
		return
	}

	list = make([]*v1pb.CapsuleGetDetailResp_List, length)
	uids := make([]int64, length)
	for ix, openInfo := range openData {
		uids[ix] = openInfo.Uid
	}
	userResult, _ := s.getUserInfo(ctx, uids, []string{"info"})
	for ix, openInfo := range openData {
		item := new(v1pb.CapsuleGetDetailResp_List)
		item.Num = openInfo.Num
		item.Name = strconv.FormatInt(openInfo.Uid, 10)
		if userInfo, ok := userResult[openInfo.Uid]; ok {
			item.Name = userInfo.Name
		}
		item.Gift = openInfo.Name
		item.Date = openInfo.Date
		list[ix] = item
	}
	return
}

// GetUserInfo
// 获取user信息
func (s *CapsuleService) getUserInfo(ctx context.Context, uids []int64, attr []string) (userResult map[int64]*actmdl.Card, err error) {
	lens := len(uids)
	if lens <= 0 {
		return
	}
	userResult, err = dao.AccountApi.Cards3(ctx, &actmdl.ArgMids{Mids: uids})
	if err != nil {
		log.Error("[service.capsule | getUserInfo] dao.AccountApi.Cards3 error : %v", err)
		return
	}
	return
}

// OpenCapsule implementation
// * web端开扭蛋
// 需要登录态; 会清除红点
func (s *CapsuleService) OpenCapsule(ctx context.Context, req *v1pb.CapsuleOpenCapsuleReq) (resp *v1pb.CapsuleOpenCapsuleResp, err error) {
	resp = &v1pb.CapsuleOpenCapsuleResp{}
	openCount := req.Count
	var coinId int64
	if _, ok := dao.CoinIdStringMap[req.Type]; !ok {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	coinId = dao.CoinIdStringMap[req.Type]
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil || coinConfMap == nil {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	if _, ok := coinConfMap[coinId]; !ok {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	isOpen := s.dao.IsPoolOpen(coinConfMap[coinId], coinId)
	if !isOpen {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	if openCount != 1 && openCount != 10 && openCount != 100 {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	if _, ok := coinConfMap[coinId]; !ok {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	if coinConfMap[coinId].PoolConf == nil || len(coinConfMap[coinId].PoolConf.PoolPrize) == 0 {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	action := "open"
	status, pInfo, err := s.dao.PayCoin(context.TODO(), req.Uid, coinConfMap[coinId], req.Count, action, req.Platform)
	if err != nil {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	if status == 1 {
		return resp, ecode.Error(ecode.XLotteryCapsuleCoinNotEnough, "扭蛋币不足")
	}
	userInfo, e := s.dao.GetUserCapsuleInfo(context.TODO(), req.Uid)
	if e == nil {
		resp.Info = &v1pb.CapsuleOpenCapsuleResp_Info{}
		resp.Info.Normal = &v1pb.CapsuleOpenCapsuleResp_CapsuleInfo{}
		resp.Info.Normal.Change = s.dao.GetChangeNum(ctx, req.Uid, dao.NormalCoinId)
		resp.Info.Normal.Coin = s.dao.GetCoin(userInfo[dao.NormalCoinId], coinConfMap[dao.NormalCoinId])
		resp.Info.Normal.Progress = s.dao.GetProgress(userInfo[dao.NormalCoinId], coinConfMap[dao.NormalCoinId])
		resp.Info.Colorful = &v1pb.CapsuleOpenCapsuleResp_CapsuleInfo{}
		resp.Info.Colorful.Change = s.dao.GetChangeNum(ctx, req.Uid, dao.ColorfulCoinId)
		resp.Info.Colorful.Coin = s.dao.GetCoin(userInfo[dao.ColorfulCoinId], coinConfMap[dao.ColorfulCoinId])
		resp.Info.Colorful.Progress = s.dao.GetProgress(userInfo[dao.ColorfulCoinId], coinConfMap[dao.ColorfulCoinId])
	}
	var awards []*dao.CapsulePoolPrize
	resp.IsEntity, awards, resp.Text = s.getAward(ctx, req.Uid, coinConfMap[coinId], req.Count, 1)
	s.dao.ClearNoticeBoth(ctx, req.Uid)
	f := func(ctx context.Context) {
		s.sendAward(ctx, req.Uid, awards)
		s.dao.LogAward(ctx, req.Uid, coinId, awards)
		s.dao.ReportCapsuleChange(ctx, coinId, req.Uid, openCount*(coinConfMap[coinId].ChangeNum), action, req.Platform, pInfo, userInfo, coinConfMap[coinId])
		awardBytes, _ := json.Marshal(awards)
		log.Info("[service.capsule | sendAward] uid: %d awards: %s", req.Uid, string(awardBytes))
	}
	if len(awards) > 0 {
		se := s.runCache.Do(context.Background(), func(ctx context.Context) {
			f(ctx)
		})
		if se != nil {
			log.Info("[service.capsule | sendAward] runCache is full")
			f(ctx)
		}
		resp.Awards = make([]*v1pb.CapsuleOpenCapsuleResp_Award, len(awards))
		for ix, award := range awards {
			resp.Awards[ix] = &v1pb.CapsuleOpenCapsuleResp_Award{}
			resp.Awards[ix].Num = award.Num
			resp.Awards[ix].Name = award.Name
			resp.Awards[ix].Img = award.MobileImage
			resp.Awards[ix].WebImage = award.WebImage
			resp.Awards[ix].MobileImage = award.MobileImage
			resp.Awards[ix].Usage = new(v1pb.Usage)
			resp.Awards[ix].Usage.Text = award.Description
			resp.Awards[ix].Usage.Url = award.JumpUrl
			resp.Awards[ix].Text = award.Name + " X " + strconv.FormatInt(award.Num, 10)
		}
	}

	return
}

func (s *CapsuleService) getAward(ctx context.Context, uid int64, coinConf *dao.CapsuleCoinConf, openCount int64, entryCount int) (bool, []*dao.CapsulePoolPrize, []string) {
	allAwards := make(map[int64]*dao.CapsulePoolPrize)
	var i, giftId int64
	var isGetFixAward, isEntity bool
	isGetFixAward = true
	entryMap := make(map[int64]bool)
	for i = 0; i < openCount; i++ {
		cnt := s.dao.IncrOpenCount(ctx, coinConf.Id)
		award := s.dao.OpenCapsule(ctx, uid, coinConf, i, cnt, isGetFixAward, entryMap)
		if award == nil {
			continue
		}
		if award.ProType == dao.ProTypeFixDay || award.ProType == dao.ProTypeFixWeek {
			entryMap[award.Id] = true
			if len(entryMap) >= entryCount {
				isGetFixAward = false
			}
		}
		aType := award.Type*1000000 + award.ObjectId
		if _, ok := allAwards[aType]; !ok {
			allAwards[aType] = award
		} else {
			num := allAwards[aType].Num
			if award.Type == dao.CapsulePrizeTitleType && award.Expire > allAwards[aType].Expire {
				allAwards[aType] = award
			}
			if allAwards[aType].Weight < award.Weight {
				allAwards[aType] = award
			}
			allAwards[aType].Num = num + award.Num
		}
		if s.dao.IsAwardEntry(award.Type) {
			isEntity = true
		}
	}
	filterAwards := s.filterAward(ctx, uid, coinConf, allAwards)
	awards := make([]*dao.CapsulePoolPrize, len(filterAwards))
	text := make([]string, len(filterAwards))
	ix := 0
	for _, award := range filterAwards {
		text[ix] = award.Name + " X " + strconv.FormatInt(award.Num, 10)
		awards[ix] = award
		giftId = award.Type
		ix++
	}
	if len(awards) == 1 && giftId == dao.CapsulePrizeGift1Type && awards[0].Num == 1 {
		text = []string{"你什么都没有抽中，获得一根辣条作为安慰。"}
	}
	return isEntity, awards, text
}

func (s *CapsuleService) filterWeekAward(ctx context.Context, uid int64, coinConf *dao.CapsuleCoinConf, allAwards map[int64]*dao.CapsulePoolPrize) map[int64]*dao.CapsulePoolPrize {
	// 获取保底奖品
	bottomPrize := s.dao.GetBottomPrize(ctx, coinConf)
	if bottomPrize == nil {
		return allAwards
	}
	bottomPrize.Num = 0
	titleChangeMap := make(map[int64][3]int64)
	for atype, award := range allAwards {
		// 特殊逻辑
		if award.Type == dao.CapsulePrizeTitleType {
			if award.Expire != dao.CapsulePrizeExpireForever {
				titleChangeMap[award.ObjectId] = [3]int64{1, award.Expire, atype}
			}
			if award.Num > 1 {
				bottomPrize.Num += award.Num - 1
			}
			award.Num = 1
		}
		allAwards[atype] = award
	}
	rcTimeout := time.Duration(conf.GetTimeout("getAllTitle", 50)) * time.Millisecond
	rcCtx, cancelRc := context.WithTimeout(ctx, rcTimeout)
	defer cancelRc()
	if len(titleChangeMap) > 0 {
		allTitleData, e := dao.RcApi.V1UserTitle.GetAll(rcCtx, &v12.UserTitleGetAllReq{Uid: uid})
		if e != nil {
			log.Info("[service.capsule | getAward] GetAll err uid(%d) err(%v)", uid, e)
			for titleId, cconf := range titleChangeMap {
				bottomPrize.Num += cconf[0]
				titleChangeMap[titleId] = [3]int64{0, cconf[1], cconf[2]}
			}
		}
		for titleId, cconf := range titleChangeMap {
			if cconf[0] == 0 {
				continue
			}
			if allTitleData == nil || len(allTitleData.Data) == 0 {
				continue
			}
			if _, ok := allTitleData.Data[titleId]; !ok {
				continue
			}
			expireTime := s.dao.GetExpireTime(cconf[1])
			if expireTime.Format("2006-01-02 15:04:05") > allTitleData.Data[titleId].ExpireTime {
				continue
			}
			bottomPrize.Num += cconf[0]
			titleChangeMap[titleId] = [3]int64{0, cconf[1], cconf[2]}
		}

	}
	for _, cconf := range titleChangeMap {
		if cconf[0] == 0 {
			delete(allAwards, cconf[2])
		}
	}
	if bottomPrize.Num == 0 {
		return allAwards
	}
	aType := bottomPrize.Type*1000000 + bottomPrize.ObjectId
	if _, ok := allAwards[aType]; ok {
		allAwards[aType].Num += bottomPrize.Num
	} else {
		allAwards[aType] = bottomPrize
	}
	return allAwards
}

func (s *CapsuleService) filterLplAward(ctx context.Context, uid int64, coinConf *dao.CapsuleCoinConf, allAwards map[int64]*dao.CapsulePoolPrize) map[int64]*dao.CapsulePoolPrize {
	// 获取保底奖品
	bottomPrize := s.dao.GetBottomPrize(ctx, coinConf)
	if bottomPrize == nil {
		return allAwards
	}
	bottonNum := bottomPrize.Num
	bottomPrize.Num = 0
	var colorChange int64
	for atype, award := range allAwards {
		// 特殊逻辑
		if award.Type == dao.CapsulePrizeDanmuColor {
			colorChange = atype
			if award.Num > 1 {
				bottomPrize.Num += (award.Num - 1) * bottonNum
			}
			award.Num = 1
		}
		allAwards[atype] = award
	}
	danmuTimeout := time.Duration(conf.GetTimeout("getAllColor", 50)) * time.Millisecond
	userExtCtx, cancelExt := context.WithTimeout(ctx, danmuTimeout)
	defer cancelExt()
	if colorChange > 0 {
		allColorData, e := dao.UserExtApi.V1DanmuConf.GetAllPrivilege(userExtCtx, &v13.DanmuConfGetAllPrivilegeReq{Uid: uid})
		if e != nil {
			log.Info("[service.capsule | getAward] GetAllPrivilege err uid(%d) err(%v)", uid, e)
			bottomPrize.Num += bottonNum
			delete(allAwards, colorChange)
		} else {
			isSend := true
			if allColorData != nil && len(allColorData.Data) > 0 {
				expireTime := s.dao.GetExpireTime(dao.CapsulePrizeExpire3Day).Format("2006-01-02 15:04:05")
				for _, colorData := range allColorData.Data {
					if colorData.Type != "0" || colorData.Value != "16764789" || colorData.Roomid != "0" {
						continue
					}
					if expireTime <= colorData.ExpiredAt {
						isSend = false
					}
				}
			}
			if !isSend {
				bottomPrize.Num += bottonNum
				delete(allAwards, colorChange)
			}
		}
	}

	if bottomPrize.Num == 0 {
		return allAwards
	}
	aType := bottomPrize.Type*1000000 + bottomPrize.ObjectId
	if _, ok := allAwards[aType]; ok {
		allAwards[aType].Num += bottomPrize.Num
	} else {
		allAwards[aType] = bottomPrize
	}
	return allAwards
}

func (s *CapsuleService) filterBlessAward(ctx context.Context, uid int64, coinConf *dao.CapsuleCoinConf, allAwards map[int64]*dao.CapsulePoolPrize) map[int64]*dao.CapsulePoolPrize {
	// 获取保底奖品
	bottomPrize := s.dao.GetBottomPrize(ctx, coinConf)
	if bottomPrize == nil {
		return allAwards
	}
	bottonNum := bottomPrize.Num
	bottomPrize.Num = 0
	var awardTypes = make([]string, 0)
	var awardTypeMap = make(map[string][2]int64)
	for atype, award := range allAwards {
		// 特殊逻辑
		if award.Type >= dao.CapsulePrizeCoupon1 {
			if award.Num > 1 {
				bottomPrize.Num += (award.Num - 1) * bottonNum
			}
			award.Num = 1
			etype := "CouponRetry" + strconv.FormatInt(award.Id, 10)
			awardTypes = append(awardTypes, etype)
			awardTypeMap[etype] = [2]int64{atype, award.Type}
		}
		allAwards[atype] = award
	}
	if len(awardTypes) == 0 {
		return allAwards
	}
	extraData, err := s.dao.GetExtraDataByType(ctx, uid, awardTypes)
	if err != nil {
		for _, typeConf := range awardTypeMap {
			bottomPrize.Num += bottonNum
			delete(allAwards, typeConf[0])
		}
	} else {
		for _, extra := range extraData {
			if _, ok := awardTypeMap[extra.Type]; !ok {
				continue
			}
			atype := awardTypeMap[extra.Type][0]
			bottomPrize.Num += bottonNum
			delete(allAwards, atype)
			delete(awardTypeMap, extra.Type)
		}
		awardTypes = make([]string, 0)
		for etype := range awardTypeMap {
			awardTypes = append(awardTypes, etype)
		}
		if len(awardTypes) > 0 {
			extraData = make([]*model.ExtraData, len(awardTypes))
			for ix, etype := range awardTypes {
				extraData[ix] = &model.ExtraData{
					Uid:       uid,
					Type:      etype,
					ItemValue: 0,
					ItemExtra: strconv.FormatInt(awardTypeMap[etype][1], 10),
				}
			}
			status, err := s.dao.AddExtraData(ctx, extraData)
			if err != nil || !status {
				for _, etype := range awardTypes {
					atype := awardTypeMap[etype][0]
					bottomPrize.Num += bottonNum
					delete(allAwards, atype)
				}
			}
		}
	}
	if bottomPrize.Num == 0 {
		return allAwards
	}
	aType := bottomPrize.Type*1000000 + bottomPrize.ObjectId
	if _, ok := allAwards[aType]; ok {
		allAwards[aType].Num += bottomPrize.Num
	} else {
		allAwards[aType] = bottomPrize
	}
	return allAwards
}

func (s *CapsuleService) filterAward(ctx context.Context, uid int64, coinConf *dao.CapsuleCoinConf, allAwards map[int64]*dao.CapsulePoolPrize) map[int64]*dao.CapsulePoolPrize {
	if len(allAwards) == 0 {
		return allAwards
	}
	var filterCoinId int64
	for _, award := range allAwards {
		// 特殊逻辑
		if coinConf.Id == dao.WeekCoinId && award.Type == dao.CapsulePrizeTitleType {
			if award.Expire != dao.CapsulePrizeExpireForever || award.Num > 0 {
				filterCoinId = dao.WeekCoinId
			}
		}
		if coinConf.Id == dao.LplCoinId && award.Type == dao.CapsulePrizeDanmuColor {
			filterCoinId = dao.LplCoinId
		}
		if coinConf.Id == dao.BlessCoinId && award.Type >= dao.CapsulePrizeCoupon1 {
			filterCoinId = dao.BlessCoinId
		}
	}
	if filterCoinId == 0 {
		return allAwards
	}
	switch filterCoinId {
	case dao.WeekCoinId:
		allAwards = s.filterWeekAward(ctx, uid, coinConf, allAwards)
	case dao.LplCoinId:
		allAwards = s.filterLplAward(ctx, uid, coinConf, allAwards)
	case dao.BlessCoinId:
		allAwards = s.filterBlessAward(ctx, uid, coinConf, allAwards)
	}
	return allAwards
}

func (s *CapsuleService) sendAward(ctx context.Context, uid int64, awards []*dao.CapsulePoolPrize) {
	for _, award := range awards {
		if s.dao.IsAwardEntry(award.Type) {
			s.sendProductAward(ctx, uid, award)
		} else if s.dao.IsAwardCoupon(award.Type) {
			s.sendCouponAward(ctx, uid, award)
		} else {
			switch award.Type {
			case dao.CapsulePrizeGift1Type:
				ts := s.dao.GetExpireTime(award.Expire)
				_, err := dao.GiftApi.V1Gift.AddFreeGift(ctx, &v1.GiftAddFreeGiftReq{Uid: uid, Giftid: 1, Num: award.Num, Expireat: ts.Unix()})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddFreeGift uid: %d num: %d error : %v", uid, award.Num, err)
				}
				log.Info("[service.capsule | sendAward] AddFreeGift uid: %d num: %d ts : %d", uid, award.Num, ts.Unix())
			case dao.CapsulePrizeTitleType:
				ts := s.dao.GetExpireTime(award.Expire)
				expireTime := "0"
				if ts.Unix() > 0 {
					expireTime = ts.Format("2006-01-02 15:04:05")
				}
				var i int64
				for i = 0; i < award.Num; i++ {
					_, err := dao.RcApi.V1UserTitle.Add(ctx, &v12.UserTitleAddReq{Uid: uid, TitleId: award.ObjectId, ExpireTime: expireTime})
					if err != nil {
						log.Error("[service.capsule | sendAward] AddTitle uid: %d title: %d error : %v", uid, award.ObjectId, err)
					}
				}
				log.Info("[service.capsule | sendAward] AddTitle uid: %d title: %d ts : %s", uid, award.ObjectId, expireTime)
			case dao.CapsulePrizeStuff1Type:
				_, err := dao.RcApi.V1TitleStuff.Add(ctx, &v12.TitleStuffAddReq{Uid: uid, Sid: 1, Num: award.Num})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddStuff uid: %d stuff: %d error : %v", uid, 1, err)
				}
			case dao.CapsulePrizeStuff2Type:
				_, err := dao.RcApi.V1TitleStuff.Add(ctx, &v12.TitleStuffAddReq{Uid: uid, Sid: 2, Num: award.Num})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddStuff uid: %d stuff: %d error : %v", uid, 2, err)
				}
			case dao.CapsulePrizeStuff3Type:
				_, err := dao.RcApi.V1TitleStuff.Add(ctx, &v12.TitleStuffAddReq{Uid: uid, Sid: 3, Num: award.Num})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddStuff uid: %d stuff: %d error : %v", uid, 3, err)
				}
			case dao.CapsulePrizeGuard3Type:
				var i int64
				for i = 0; i < award.Num; i++ {
					isAdd := s.dao.AddAward(ctx, uid, time.Now().Add(15*86400*time.Second).Format("2006-01-02 15:04:05"), 3, award.Name, 1, "大航海舰长体验", 0)
					if !isAdd {
						log.Error("[service.capsule | sendAward] AddAward uid: %d name: %s error", uid, award.Name)
					}
				}
			case dao.CapsulePrizeGuard2Type:
				var i int64
				for i = 0; i < award.Num; i++ {
					isAdd := s.dao.AddAward(ctx, uid, time.Now().Add(15*86400*time.Second).Format("2006-01-02 15:04:05"), 3, award.Name, 1, "大航海提督体验", 0)
					if !isAdd {
						log.Error("[service.capsule | sendAward] AddAward uid: %d name: %s error", uid, award.Name)
					}
				}
			case dao.CapsulePrizeGuard1Type:
				var i int64
				for i = 0; i < award.Num; i++ {
					isAdd := s.dao.AddAward(ctx, uid, time.Now().Add(15*86400*time.Second).Format("2006-01-02 15:04:05"), 3, award.Name, 1, "大航海总督体验", 0)
					if !isAdd {
						log.Error("[service.capsule | sendAward] AddAward uid: %d name: %s error", uid, award.Name)
					}
				}
			case dao.CapsulePrizeScoreAdd:
				ts := s.dao.GetExpireTime(dao.CapsulePrizeExpire1Week)
				_, err := dao.GiftApi.V1Gift.AddFreeGift(ctx, &v1.GiftAddFreeGiftReq{Uid: uid, Giftid: 30081, Num: award.Num, Expireat: ts.Unix()})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddFreeGift uid: %d num: %d error : %v", uid, award.Num, err)
				}
				log.Info("[service.capsule | sendAward] AddFreeGift uid: %d num: %d  type: %d ts : %d", uid, award.Num, dao.CapsulePrizeScoreAdd, ts.Unix())
			case dao.CapsulePrizeSmallStar:
				ts := s.dao.GetExpireTime(dao.CapsulePrizeExpire1Week)
				_, err := dao.GiftApi.V1Gift.AddFreeGift(ctx, &v1.GiftAddFreeGiftReq{Uid: uid, Giftid: 30085, Num: award.Num, Expireat: ts.Unix()})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddFreeGift uid: %d num: %d type: %d  error : %v", uid, award.Num, dao.CapsulePrizeSmallStar, err)
				}
				log.Info("[service.capsule | sendAward] AddFreeGift uid: %d num: %d  type: %d ts : %d", uid, award.Num, dao.CapsulePrizeSmallStar, ts.Unix())
			case dao.CapsulePrizeLplScore:
				coinConfMap, err := s.dao.GetCapsuleConf(ctx)
				if err != nil || len(coinConfMap) == 0 {
					log.Error("[service.capsule | sendAward] GetCapsuleConf uid: %d num: %d type: %d  error : %v", uid, award.Num, dao.CapsulePrizeWeekScore, err)
					continue
				}
				if _, ok := coinConfMap[dao.LplCoinId]; !ok {
					log.Error("[service.capsule | sendAward] WeekCoinConf uid: %d num: %d type: %d  error : %v", uid, award.Num, dao.CapsulePrizeWeekScore, err)
					continue
				}
				pInfo, _ := s.dao.GetUserInfo(ctx, uid, dao.LplCoinId)
				s.dao.UpdateCapsule(ctx, uid, dao.LplCoinId, coinConfMap[dao.LplCoinId].ChangeNum*award.Num, "lottery", "", coinConfMap[dao.LplCoinId])
				userInfo, _ := s.dao.GetUserInfo(ctx, uid, dao.LplCoinId)
				s.dao.ReportCapsuleChange(ctx, dao.LplCoinId, uid, coinConfMap[dao.LplCoinId].ChangeNum*award.Num, "lottery", "", pInfo, userInfo, coinConfMap[dao.LplCoinId])
			case dao.CapsulePrizeWeekScore:
				coinConfMap, err := s.dao.GetCapsuleConf(ctx)
				if err != nil || len(coinConfMap) == 0 {
					log.Error("[service.capsule | sendAward] GetCapsuleConf uid: %d num: %d type: %d  error : %v", uid, award.Num, dao.CapsulePrizeWeekScore, err)
					continue
				}
				if _, ok := coinConfMap[dao.WeekCoinId]; !ok {
					log.Error("[service.capsule | sendAward] WeekCoinConf uid: %d num: %d type: %d  error : %v", uid, award.Num, dao.CapsulePrizeWeekScore, err)
					continue
				}
				pInfo, _ := s.dao.GetUserInfo(ctx, uid, dao.WeekCoinId)
				s.dao.UpdateCapsule(ctx, uid, dao.WeekCoinId, coinConfMap[dao.WeekCoinId].ChangeNum*award.Num, "lottery", "", coinConfMap[dao.WeekCoinId])
				userInfo, _ := s.dao.GetUserInfo(ctx, uid, dao.WeekCoinId)
				s.dao.ReportCapsuleChange(ctx, dao.WeekCoinId, uid, coinConfMap[dao.WeekCoinId].ChangeNum*award.Num, "lottery", "", pInfo, userInfo, coinConfMap[dao.WeekCoinId])
			case dao.CapsulePrizeDanmuColor:
				ts := s.dao.GetExpireTime(dao.CapsulePrizeExpire3Day)
				_, err := dao.UserExtApi.V1DanmuConf.AddByType(ctx, &v13.DanmuConfAddByTypeReq{Uid: uid, Type: "color", Roomid: 0, Expire: ts.Unix(), Value: 16764789})
				if err != nil {
					log.Error("[service.capsule | sendAward] AddByType uid: %d error : %v", uid, err)
				}
				log.Info("[service.capsule | sendAward] AddByType uid: %d ts : %d", uid, ts.Unix())
			}

		}

	}
}
func (s *CapsuleService) sendProductAward(ctx context.Context, uid int64, award *dao.CapsulePoolPrize) {
	var i int64
	for i = 0; i < award.Num; i++ {
		isAdd := s.dao.AddAward(ctx, uid, time.Now().Add(15*86400*time.Second).Format("2006-01-02 15:04:05"), 2, award.Name, 1, "扭蛋活动", 0)
		if !isAdd {
			log.Error("[service.capsule | sendAward] AddAward uid: %d name: %s error", uid, award.Name)
		}
	}
}

func (s *CapsuleService) sendCouponAward(ctx context.Context, uid int64, award *dao.CapsulePoolPrize) {
	if s.conf.CouponConf == nil || s.conf.CouponConf.Url == "" || len(s.conf.CouponConf.Coupon) == 0 {
		log.Error("[service.capsule | sendAward] CouponConf is empty")
		return
	}
	atype := strconv.FormatInt(award.Type, 10)
	if _, ok := s.conf.CouponConf.Coupon[atype]; !ok {
		log.Error("[service.capsule | sendAward] CouponConf.coupon is empty")
		return
	}
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}
	endPoint := s.conf.CouponConf.Url
	postJson := make(map[string]interface{})
	postJson["mid"] = uid
	postJson["couponId"] = s.conf.CouponConf.Coupon[atype]
	bytesData, err := json.Marshal(postJson)
	if err != nil {
		log.Error("[service.capsule | sendAward] json.Marshal(%v) error(%v)", postJson, err)
		return
	}
	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(bytesData))
	if err != nil {
		log.Error("[service.capsule | sendAward] http.NewRequest(%v) url(%v) error(%v)", postJson, endPoint, err)
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	log.Info("coupon vip mid(%d) couponID(%s)", uid, s.conf.CouponConf.Coupon[atype])
	if err = s.httpClient.Do(ctx, req, &res); err != nil {
		log.Error("[service.capsule | sendAward] s.client.Do error(%v)", err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("coupon vip url(%v) res code(%d)", endPoint, res.Code)
		return
	}
	log.Info("[service.capsule | sendAward] s.client.Do endpoint (%v) req (%v)", endPoint, postJson)
	etype := "CouponRetry" + strconv.FormatInt(award.Id, 10)
	s.dao.UpdateExtraValue(ctx, uid, etype, 1)
}

//定时重置Capusule
func (s *CapsuleService) tickerReloadCapsuleConf(ctx context.Context) {
	changeFlag, _ := s.dao.GetCapsuleChangeFlag(ctx)
	s.dao.RelaodCapsuleConfig(ctx, changeFlag)
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			redisChangeFlag, _ := s.dao.GetCapsuleChangeFlag(ctx)
			capsuleCacheTime, capsuleChangeFlag := s.dao.GetCapsuleChangeInfo(ctx)
			log.Info("[service.capsule | reload] redis %d cache %d", redisChangeFlag, capsuleChangeFlag)
			if redisChangeFlag != capsuleChangeFlag || time.Now().Unix()-capsuleCacheTime > 60 {
				s.dao.RelaodCapsuleConfig(ctx, redisChangeFlag)
			}
		}
	}()
}

// GetPrizeType implementation
func (s *CapsuleService) GetPrizeType(ctx context.Context, req *v1pb.GetPrizeTypeReq) (resp *v1pb.GetPrizeTypeResp, err error) {
	resp = &v1pb.GetPrizeTypeResp{}
	resp.List = make([]*v1pb.GetPrizeTypeResp_List, len(dao.PrizeNameMap))
	ix := 0
	for ptype, pname := range dao.PrizeNameMap {
		resp.List[ix] = &v1pb.GetPrizeTypeResp_List{Type: ptype, Name: pname}
		ix++
	}
	sort.Slice(resp.List, func(i, j int) bool {
		return resp.List[i].Type < resp.List[j].Type
	})
	return
}

// GetPrizeExpire implementation
func (s *CapsuleService) GetPrizeExpire(ctx context.Context, req *v1pb.GetPrizeExpireReq) (resp *v1pb.GetPrizeExpireResp, err error) {
	resp = &v1pb.GetPrizeExpireResp{}
	resp.List = make([]*v1pb.GetPrizeExpireResp_List, len(dao.PrizeExpireMap))
	ix := 0
	for expire, pname := range dao.PrizeExpireMap {
		resp.List[ix] = &v1pb.GetPrizeExpireResp_List{Expire: expire, Name: pname}
		ix++
	}
	return
}

// GetCapsuleInfo implementation
func (s *CapsuleService) GetCapsuleInfo(ctx context.Context, req *v1pb.CapsuleGetCapsuleInfoReq) (resp *v1pb.CapsuleGetCapsuleInfoResp, err error) {
	resp = &v1pb.CapsuleGetCapsuleInfoResp{}
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil || len(coinConfMap) == 0 {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	var coinId int64
	if _, ok := dao.CoinIdStringMap[req.Type]; !ok {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	coinId = dao.CoinIdStringMap[req.Type]
	if _, ok := coinConfMap[coinId]; !ok {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	userInfo := make(map[int64]int64)
	if req.Uid > 0 {
		userInfo, err = s.dao.GetUserInfo(ctx, req.Uid, coinId)
	} else {
		userInfo[coinId] = 0
	}
	log.Info("[service.capsule | GetCapsuleInfo] uid(%d) %d error(%v)", req.Uid, userInfo[coinId], err)
	if err != nil {
		err = ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		return
	}
	coinConf := coinConfMap[coinId]
	resp.Coin = s.dao.GetCoin(userInfo[coinId], coinConf)
	resp.Rule = coinConf.PoolConf.Rule
	if coinConf.PoolConf != nil && len(coinConf.PoolConf.PoolPrize) > 0 {
		giftMap := make(map[int64]struct{})
		prizeList := make([]*v1pb.CapsuleGetCapsuleInfoResp_GiftList, len(coinConf.PoolConf.PoolPrize))
		prizeFilter := make([]*v1pb.CapsuleGetCapsuleInfoResp_GiftFilter, 0)
		for ix, prize := range coinConf.PoolConf.PoolPrize {
			var expire string
			if _, ok := dao.PrizeExpireMap[prize.Expire]; ok {
				expire = dao.PrizeExpireMap[prize.Expire]
			}
			prizeList[ix] = &v1pb.CapsuleGetCapsuleInfoResp_GiftList{
				Id:        prize.Id,
				Name:      prize.Name,
				Num:       prize.Num,
				Weight:    prize.Weight,
				MobileUrl: prize.MobileImage,
				WebUrl:    prize.WebImage,
				Usage: &v1pb.Usage{
					Url:  prize.JumpUrl,
					Text: prize.Description,
				},
				Type:   prize.Type,
				Expire: expire,
			}
			gType := prize.Type*1000000 + prize.ObjectId
			if _, ok := giftMap[gType]; ok {
				continue
			}
			fprize := &v1pb.CapsuleGetCapsuleInfoResp_GiftFilter{
				Id:        prize.Id,
				Name:      prize.Name,
				MobileUrl: prize.MobileImage,
				WebUrl:    prize.WebImage,
				Usage: &v1pb.Usage{
					Url:  prize.JumpUrl,
					Text: prize.Description,
				},
			}
			giftMap[gType] = struct{}{}
			prizeFilter = append(prizeFilter, fprize)
		}
		resp.GiftList = prizeList
		resp.GiftFilter = prizeFilter

	}
	return
}

// OpenCapsuleByType implementation
func (s *CapsuleService) OpenCapsuleByType(ctx context.Context, req *v1pb.CapsuleOpenCapsuleByTypeReq) (resp *v1pb.CapsuleOpenCapsuleByTypeResp, err error) {
	resp = &v1pb.CapsuleOpenCapsuleByTypeResp{}
	openCount := req.Count
	var coinId int64
	if _, ok := dao.CoinIdStringMap[req.Type]; !ok {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	coinId = dao.CoinIdStringMap[req.Type]
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil || coinConfMap == nil {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	if _, ok := coinConfMap[coinId]; !ok {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	if coinId < 3 {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	isOpen := s.dao.IsPoolOpen(coinConfMap[coinId], coinId)
	if !isOpen {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	if openCount != 1 && openCount != 10 && openCount != 100 {
		return resp, ecode.Error(ecode.InvalidParam, "参数错误")
	}
	if _, ok := coinConfMap[coinId]; !ok {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	if coinConfMap[coinId].PoolConf == nil || len(coinConfMap[coinId].PoolConf.PoolPrize) == 0 {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	action := "open"
	status, pInfo, err := s.dao.PayCapsule(context.TODO(), req.Uid, coinConfMap[coinId], req.Count, action, req.Platform)
	if err != nil {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	if status == 1 {
		return resp, ecode.Error(ecode.XLotteryCapsuleOperationFrequent, "你的操作太频繁啦")
	} else if status == 2 {
		return resp, ecode.Error(ecode.XLotteryCapsuleCoinNotEnough, "抽奖券不足")
	}
	userInfo, e := s.dao.GetUserInfo(context.TODO(), req.Uid, coinId)
	if e == nil {
		resp.Info = &v1pb.CapsuleOpenCapsuleByTypeResp_CapsuleInfo{}
		resp.Info.Coin = s.dao.GetCoin(userInfo[coinId], coinConfMap[coinId])
	}
	entryCount := 1
	if req.Count == 100 && coinId == dao.BlessCoinId {
		entryCount = 3
	}
	var awards []*dao.CapsulePoolPrize
	resp.IsEntity, awards, _ = s.getAward(ctx, req.Uid, coinConfMap[coinId], req.Count, entryCount)

	s.dao.ClearNoticeBoth(ctx, req.Uid)
	f := func(ctx context.Context) {
		s.sendAward(ctx, req.Uid, awards)
		s.dao.LogAward(ctx, req.Uid, coinId, awards)
		s.dao.ReportCapsuleChange(ctx, coinId, req.Uid, openCount*(coinConfMap[coinId].ChangeNum), action, req.Platform, pInfo, userInfo, coinConfMap[coinId])
		awardBytes, _ := json.Marshal(awards)
		log.Info("[service.capsule | sendAward] uid: %d awards: %s", req.Uid, string(awardBytes))
	}
	if len(awards) > 0 {
		se := s.runCache.Do(context.Background(), func(ctx context.Context) {
			f(ctx)
		})
		if se != nil {
			log.Info("[service.capsule | sendAward] runCache is full")
			f(ctx)
		}

		resp.Awards = make([]*v1pb.CapsuleOpenCapsuleByTypeResp_Award, len(awards))
		for ix, award := range awards {
			var expire string
			if _, ok := dao.PrizeExpireMap[award.Expire]; ok {
				expire = dao.PrizeExpireMap[award.Expire]
			}
			resp.Awards[ix] = &v1pb.CapsuleOpenCapsuleByTypeResp_Award{}
			resp.Awards[ix].Id = award.Id
			resp.Awards[ix].Num = award.Num
			resp.Awards[ix].Name = award.Name
			resp.Awards[ix].WebUrl = award.WebImage
			resp.Awards[ix].MobileUrl = award.MobileImage
			resp.Awards[ix].Usage = new(v1pb.Usage)
			resp.Awards[ix].Usage.Text = award.Description
			resp.Awards[ix].Usage.Url = award.JumpUrl
			resp.Awards[ix].Text = award.Name + " X " + strconv.FormatInt(award.Num, 10)
			resp.Awards[ix].Weight = award.Weight
			resp.Awards[ix].Type = award.Type
			resp.Awards[ix].Expire = expire
		}
		sort.Slice(resp.Awards, func(i, j int) bool {
			return resp.Awards[i].Weight > resp.Awards[j].Weight
		})
		resp.Text = make([]string, len(awards))
		for ix, award := range resp.Awards {
			resp.Text[ix] = award.Text
		}
	}
	return
}

// GetCouponList implementation
func (s *CapsuleService) GetCouponList(ctx context.Context, req *v1pb.CapsuleGetCouponListReq) (resp *v1pb.CapsuleGetCouponListResp, err error) {
	resp = &v1pb.CapsuleGetCouponListResp{}
	coinConfMap, err := s.dao.GetCapsuleConf(ctx)
	if err != nil {
		return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
	}
	awardMap := make(map[string]string)
	awardIds := make([]string, 0)
	for _, coinConf := range coinConfMap {
		if coinConf.PoolConf == nil || len(coinConf.PoolConf.PoolPrize) == 0 {
			continue
		}
		for _, award := range coinConf.PoolConf.PoolPrize {
			if s.dao.IsAwardCoupon(award.Type) {
				etype := "CouponRetry" + strconv.FormatInt(award.Id, 10)
				awardIds = append(awardIds, etype)
				awardMap[etype] = award.Name
			}
		}

	}
	couponConf := make(map[string]string)
	if s.conf.CouponConf != nil && len(s.conf.CouponConf.Coupon) > 0 {
		couponConf = s.conf.CouponConf.Coupon
	}
	resp.List = make([]*v1pb.CapsuleGetCouponListResp_List, 0)
	if len(awardIds) > 0 {
		extraData, err := s.dao.GetExtraDataByType(ctx, req.Uid, awardIds)
		if err != nil {
			return resp, ecode.Error(ecode.XLotteryCapsuleSystemErr, "系统繁忙，请稍后再试")
		}
		if len(extraData) > 0 {
			resp.List = make([]*v1pb.CapsuleGetCouponListResp_List, len(extraData))
			for ix, extra := range extraData {
				resp.List[ix] = &v1pb.CapsuleGetCouponListResp_List{}
				resp.List[ix].Uid = req.Uid
				resp.List[ix].Status = extra.ItemValue
				resp.List[ix].AwardTime = extra.Ctime
				if resp.List[ix].Status == 0 {
					resp.List[ix].RetryTime = extra.Mtime
				}
				resp.List[ix].AwardName = awardMap[extra.Type]
				resp.List[ix].AwardCode = ""
				if _, ok := couponConf[extra.ItemExtra]; ok {
					resp.List[ix].AwardCode = couponConf[extra.ItemExtra]
				}
			}
		}

	}
	return
}
