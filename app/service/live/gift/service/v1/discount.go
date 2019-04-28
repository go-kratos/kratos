package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-common/app/service/live/gift/dao"
	"go-common/app/service/live/gift/model"
	v1liveuser "go-common/app/service/live/live_user/api/liverpc/v1"
	"go-common/library/log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	v1pb "go-common/app/service/live/gift/api/grpc/v1"
)

// DiscountCache DiscountCache
type DiscountCache struct {
	discountMap  map[int64]map[int64]map[int64]int64
	discountPlan map[int64]*model.DiscountInfo
	cacheTime    int64
}

var (
	dCache         = new(DiscountCache)
	discountAtomic = new(atomic.Value)
)

// GetHighestGuardLevel 获取最高等级大航海
func (s *GiftService) GetHighestGuardLevel(ctx context.Context, uid int64) (lv int64, err error) {
	if uid == 0 {
		return
	}
	res, err := dao.LiveUserApi.V1Guard.GetByUid(ctx, &v1liveuser.GuardGetByUidReq{
		Uid:        uid,
		IncExpire:  0,
		IsLimitOne: 0,
	})
	if err != nil || res.Code != 0 {
		log.Error("call V1Guard.GetByUid error,err:(%v)", err)
		return
	}
	if res.Code != 0 {
		err = errors.New("resp code err")
		log.Error("call V1Guard.GetByUid error,code:(%v)", res.Code)
		return
	}
	data := res.Data
	if len(data) == 0 {
		return
	}
	lv = data[0].PrivilegeType
	return
}

// GetDiscountList 获取折扣道具列表
func (s *GiftService) GetDiscountList(ctx context.Context, platform string, userType, roomID, areaParentID, areaID int64) (config []*v1pb.DiscountGiftListResp_GiftInfo, err error) {
	config = make([]*v1pb.DiscountGiftListResp_GiftInfo, 0)
	err = s.LoadDiscountCache(ctx, false)
	if err != nil {
		return
	}
	platform = strings.ToLower(platform)
	platformCheck := make([]int64, 0)
	pf, ok := platformMap[platform]
	if !ok {
		platformCheck = []int64{platformMap["default"]}
	}
	if pf != platformMap["default"] {
		platformCheck = append(platformCheck, pf)
	}
	// map顺序是随机的
	sceneCheck := []int64{sceneRoom, sceneArea, sceneAreaParent, sceneAll, sceneDefault}
	sceneMap := map[int64]int64{
		sceneRoom:       roomID,
		sceneArea:       areaID,
		sceneAreaParent: areaParentID,
		sceneAll:        0,
		sceneDefault:    0,
	}
	utStr := strconv.FormatInt(userType, 10)
	dCache = discountAtomic.Load().(*DiscountCache)
	for _, pf := range platformCheck {
		for _, sk := range sceneCheck {
			sv := sceneMap[sk]
			planID, ok := dCache.discountMap[pf][sk][sv]
			if !ok {
				continue
			}
			discountInfo, ok := dCache.discountPlan[planID]
			if ok {
				for giftID, info := range discountInfo.List {
					var (
						dp int64
						cp int64
						cm string
					)
					price, suc := info["discount_price_"+utStr]
					if suc {
						dp = int64(price.(int))
					} else {
						dp = int64(0)
					}
					pos, suc := info["corner_position_"+utStr]
					if suc {
						cp = int64(pos.(int))
					} else {
						cp = int64(0)
					}
					mark, suc := info["corner_mark_"+utStr]
					if suc {
						cm = mark.(string)
					} else {
						cm = ""
					}
					tmp := &v1pb.DiscountGiftListResp_GiftInfo{
						GiftId:         giftID,
						Price:          s.GetGiftPrice(ctx, giftID),
						DiscountPrice:  dp,
						CornerMark:     cm,
						CornerPosition: cp,
						CornerColor:    "#699553",
					}
					config = append(config, tmp)
				}
			}
		}
	}
	return
}

// GetGiftPrice 获取礼物价格
func (s *GiftService) GetGiftPrice(ctx context.Context, giftID int64) (price int64) {
	giftConfig, err := s.GetAllConfig(ctx)
	if err != nil {
		return
	}
	for _, v := range giftConfig {
		if v.Id == giftID {
			return v.Price
		}
	}

	return
}

// LoadDiscountCache LoadDiscountCache
func (s *GiftService) LoadDiscountCache(ctx context.Context, force bool) (err error) {
	curTime := time.Now().Unix()
	var ok bool
	dCache, ok = discountAtomic.Load().(*DiscountCache)
	if !ok {
		dCache = new(DiscountCache)
	}
	if force || (curTime-dCache.cacheTime) > 30 {
		err = s.SyncDiscountCache(ctx)
	}
	return
}

// SyncDiscountCache SyncDiscountCache
func (s *GiftService) SyncDiscountCache(ctx context.Context) (err error) {
	dPlan, err := s.GetDiscountPlans(ctx)
	if err != nil {
		return
	}
	//platform:scene_key:scene_value:plan_id
	dMap := make(map[int64]map[int64]map[int64]int64)
	for _, v := range dPlan {
		platformList := s.parsePlatform2Slice(v.Platform)
		for _, sv := range v.SceneValue {
			for _, platform := range platformList {
				if _, ok := dMap[platform]; !ok {
					dMap[platform] = make(map[int64]map[int64]int64)
				}
				if _, ok := dMap[platform][v.SceneKey]; !ok {
					dMap[platform][v.SceneKey] = make(map[int64]int64)
				}
				dMap[platform][v.SceneKey][sv] = v.Id
			}
		}
	}
	dCache.discountPlan = dPlan
	dCache.discountMap = dMap
	dCache.cacheTime = time.Now().Unix()
	discountAtomic.Store(dCache)
	log.Info("SyncDiscountCache ,%+v", dCache)

	return
}

func (s *GiftService) parsePlatform2Slice(platform int64) (pf []int64) {
	if platform == 0 {
		pf = []int64{0}
		return pf
	}
	for _, v := range platformMap {
		if v > 0 && (v&platform) == v {
			pf = append(pf, v)
		}
	}
	return
}

// GetDiscountPlans 获取折扣计划
func (s *GiftService) GetDiscountPlans(ctx context.Context) (dp map[int64]*model.DiscountInfo, err error) {
	plans, err := s.dao.GetDiscountPlan(ctx, time.Now())
	if err != nil {
		return
	}
	disIds := make([]int64, 0)
	for _, v := range plans {
		disIds = append(disIds, v.Id)
	}
	if len(disIds) == 0 {
		return
	}
	details, err := s.dao.GetByDiscountIds(ctx, disIds)
	if err != nil {
		return
	}

	list := make(map[int64]map[int64]map[string]interface{})
	for _, v := range details {
		utStr := strconv.FormatInt(v.UserType, 10)
		list[v.DiscountId] = map[int64]map[string]interface{}{
			v.GiftId: {
				"discount_price_" + utStr:  v.DiscountPrice,
				"corner_mark_" + utStr:     v.CornerMark,
				"corner_position_" + utStr: v.CornerPosition,
				"discount_id":              v.DiscountId,
			},
		}
	}
	dp = make(map[int64]*model.DiscountInfo)
	for _, v := range plans {
		tmp := strings.Split(v.SceneValue, ",")
		sv := make([]int64, 0)
		for _, svStr := range tmp {
			svInt, _ := strconv.ParseInt(svStr, 10, 64)
			sv = append(sv, svInt)
		}
		dp[v.Id] = &model.DiscountInfo{
			Id:         v.Id,
			SceneKey:   v.SceneKey,
			SceneValue: sv,
			Platform:   v.Platform,
			List:       list[v.Id],
		}
	}
	return
}

// MyPrintf 我的打印
func (s *GiftService) MyPrintf(v interface{}, msg string) {
	dump, _ := json.MarshalIndent(v, "", "	")
	fmt.Printf("**********%s**********\n%s\n", msg, dump)
}

// FormatPrint 我的打印
func (s *GiftService) FormatPrint(v interface{}, msg string) {
	fmt.Printf("**********%s**********\n%v\n", msg, v)
}
