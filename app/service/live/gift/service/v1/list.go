package v1

import (
	"context"
	"encoding/json"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/model"
	"go-common/app/service/live/resource/sdk"
	"go-common/library/log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	//"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var privilegeGift = []int64{30047}
var defaultPlan = "20005,20008,20013,3,20009,20012,20003,25,7,8,20011,20014,20007,20002,20010,20004"
var giftCache = new(GiftCache)
var giftAtomic = new(atomic.Value)

// GiftCache GiftCache
type GiftCache struct {
	cacheTime int64
	giftList  map[int64]*model.GiftOnline
	plan      map[string]map[int64]map[int64]*model.GiftPlan
	validGift string
}

var platformMap = map[string]int64{
	"default": 0,
	"pc":      1,
	"ios":     2,
	"android": 4,
	"ipad":    8,
}

const (
	sceneAll        int64 = 0
	sceneAreaParent int64 = 1
	sceneArea       int64 = 2
	sceneRoom       int64 = 3
	sceneDefault    int64 = 99

	//coinTypeSilver = 0
	coinTypeGold = 1
	coinTypeBag  = 2
)

// GetOnlinePlanGiftList 获取在线礼物计划
func (s *GiftService) GetOnlinePlanGiftList(ctx context.Context, roomID, areaParentID, areaID int64, platform string, build int64, mobiApp string) (resp *v1pb.RoomGiftListResp, err error) {
	resp = &v1pb.RoomGiftListResp{}
	if err = s.LoadGiftCache(ctx, false); err != nil {
		return
	}
	platform = strings.ToLower(platform)
	if _, ok := platformMap[platform]; !ok {
		platform = "pc"
	}
	platformCheck := []string{platform, "default"}
	// map顺序是随机的
	sceneCheck := []int64{sceneRoom, sceneArea, sceneAreaParent, sceneAll, sceneDefault}
	sceneMap := map[int64]int64{
		sceneRoom:       roomID,
		sceneArea:       areaID,
		sceneAreaParent: areaParentID,
		sceneAll:        0,
		sceneDefault:    0,
	}
	giftCache = giftAtomic.Load().(*GiftCache)
	var plan *model.GiftPlan
	for _, pf := range platformCheck {
		for _, sceneKey := range sceneCheck {
			sceneVal := sceneMap[sceneKey]
			if plan == nil {
				plan = giftCache.plan[pf][sceneKey][sceneVal]
			}
		}
	}
	list := s.ConvertGiftList(plan.List, plan.Id, platform, build, mobiApp)
	list = s.AddPrivilegeGift(list, platform, build, mobiApp)
	silverList := s.ConvertGiftList(plan.SilverList, plan.Id, platform, build, mobiApp)
	resp = &v1pb.RoomGiftListResp{
		List:       list,
		SilverList: silverList,
	}
	return
}

// LoadGiftCache LoadGiftCache
func (s *GiftService) LoadGiftCache(ctx context.Context, needReload bool) (err error) {
	curTime := time.Now().Unix()
	var ok bool
	giftCache, ok = giftAtomic.Load().(*GiftCache)
	if !ok {
		giftCache = new(GiftCache)
	}
	if needReload || (curTime-giftCache.cacheTime) > 30 {
		err = s.SyncLocalCache(ctx)
	}
	return
}

// SyncLocalCache SyncLocalCache
func (s *GiftService) SyncLocalCache(ctx context.Context) (err error) {
	var (
		allGifts []*model.GiftOnline
		allPlans []*model.GiftPlan
	)
	eg := errgroup.Group{}
	eg.Go(func() error {
		allGifts, err = s.dao.GetAllGift(ctx)
		return err
	})
	eg.Go(func() error {
		allPlans, err = s.dao.GetOnlinePlan(ctx)
		return err
	})
	if err = eg.Wait(); err != nil {
		return
	}

	gifts := make(map[int64]*model.GiftOnline)
	for _, gift := range allGifts {
		giftID := gift.GiftId
		if s.isPrivilegeGift(giftID) {
			gift.PrivilegeRequired = 1
		}
		gifts[giftID] = gift
	}

	validGiftStr := defaultPlan
	plan := make(map[string]map[int64]map[int64]*model.GiftPlan)
	s.addDefaultPlan(plan)
	for _, v := range allPlans {
		sceneKey := v.SceneKey
		sceneVal := v.SceneValue
		platformList := s.parsePlatform(v.Platform)
		for _, platform := range platformList {
			if _, ok := plan[platform][sceneKey][sceneVal]; ok {
				continue
			}
			validGiftStr += "," + v.List + "," + v.SilverList
			if _, ok := plan[platform]; !ok {
				plan[platform] = make(map[int64]map[int64]*model.GiftPlan)
			}
			if _, ok := plan[platform][sceneKey]; !ok {
				plan[platform][sceneKey] = make(map[int64]*model.GiftPlan)
			}
			plan[platform][sceneKey][sceneVal] = v
		}
	}

	giftCache.giftList = gifts
	giftCache.plan = plan
	giftCache.validGift = validGiftStr
	giftCache.cacheTime = time.Now().Unix()
	giftAtomic.Store(giftCache)
	log.Info("SyncLocalCache Gift info,%+v", giftCache)
	return
}

func (s *GiftService) isPrivilegeGift(giftID int64) bool {
	for _, ID := range privilegeGift {
		if ID == giftID {
			return true
		}
	}
	return false
}

func (s *GiftService) parsePlatform(platform int64) (pf []string) {
	if platform == 0 {
		pf = []string{"default"}
		return pf
	}
	for k, v := range platformMap {
		if v > 0 && (v&platform) == v {
			pf = append(pf, k)
		}
	}
	return
}

func (s *GiftService) addDefaultPlan(plan map[string]map[int64]map[int64]*model.GiftPlan) {
	plan["default"] = map[int64]map[int64]*model.GiftPlan{
		sceneDefault: {
			0: &model.GiftPlan{
				Id:         0,
				List:       defaultPlan,
				SilverList: "",
			},
		},
	}
}

// ConvertGiftList ConvertGiftList
func (s *GiftService) ConvertGiftList(list string, planID int64, platform string, build int64, mobiApp string) (data []*v1pb.RoomGiftListResp_List) {
	l := strings.Split(list, ",")
	var p int64 = 1
	for _, v := range l {
		giftID, _ := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		ret := s.CheckGiftVersion(giftID, platform, build, mobiApp)
		if giftID == 0 || !ret {
			continue
		}
		l := s.NewGiftList(giftID, p, planID)
		p++
		data = append(data, l)
	}
	return
}

// CheckGiftVersion CheckGiftVersion
func (s *GiftService) CheckGiftVersion(giftID int64, platform string, build int64, mobiApp string) bool {
	giftCache = giftAtomic.Load().(*GiftCache)
	giftConfig := giftCache.giftList[giftID]
	if giftConfig == nil {
		return false
	}
	var err error
	verLimit, err := titansSdk.Get("gift_ver_limit")
	if err != nil {
		log.Error("titans get error,%v", err)
		return true
	}
	var limit map[int64]map[string]int64
	if err = json.Unmarshal([]byte(verLimit), &limit); err != nil {
		log.Error("json.Unmarshal error,%v", verLimit)
		return true
	}
	minBuild := limit[giftID][mobiApp]
	if build > 0 && minBuild > 0 && minBuild > build {
		return false
	}
	return true
}

// AddPrivilegeGift 添加大航海道具
func (s *GiftService) AddPrivilegeGift(list []*v1pb.RoomGiftListResp_List, platform string, build int64, mobiApp string) []*v1pb.RoomGiftListResp_List {
	//特殊逻辑: web端展示不支持大航海专属道具
	if platform == "pc" {
		return list
	}
	if len(list) >= 32 {
		return list
	}
	var lastPosition int64
	if len(list) == 0 {
		lastPosition = 0
	} else {
		last := list[len(list)-1]
		lastPosition = last.Position
	}
	for _, v := range privilegeGift {
		if ok := s.CheckGiftVersion(v, platform, build, mobiApp); !ok {
			continue
		}
		lastPosition++
		list = append(list, s.NewGiftList(v, lastPosition, 0))
		if len(list) >= 32 {
			return list
		}
	}
	return list
}

// NewGiftList NewGiftList
func (s *GiftService) NewGiftList(id, position, planID int64) (list *v1pb.RoomGiftListResp_List) {
	list = &v1pb.RoomGiftListResp_List{
		Id:       id,
		Position: position,
		PlanId:   planID,
	}
	return
}

func (s *GiftService) getOldList() (old []*v1pb.RoomGiftListResp_OldList) {
	old = make([]*v1pb.RoomGiftListResp_OldList, 0)
	l := &v1pb.RoomGiftListResp_OldList{
		Id:    1,
		Name:  "辣条",
		Price: 100,
		Type:  0,
		CoinType: map[string]string{
			"silver": "silver",
		},
		Img:      "https://s1.hdslb.com/bfs/live/da6656add2b14a93ed9eb55de55d0fd19f0fc7f6.png",
		GiftUrl:  "https://s1.hdslb.com/bfs/static/blive/live-assets/mobile/gift/mobilegift/2/1.gif",
		CountSet: "1,10,99,520",
		ComboNum: 0,
		SuperNum: 0,
		CountMap: map[int64]string{
			1:   "",
			10:  "",
			99:  "",
			520: "",
		},
	}
	old = append(old, l)
	return
}

// GetAllConfig 获取所有礼物配置
func (s *GiftService) GetAllConfig(ctx context.Context) (all []*v1pb.GiftConfigResp_Config, err error) {
	if err = s.LoadGiftCache(ctx, false); err != nil {
		return
	}
	giftCache = giftAtomic.Load().(*GiftCache)
	gift := giftCache.giftList
	for _, v := range gift {
		all = append(all, s.ConvertDB2Config(v))
	}
	return
}

// ConvertDB2Config ConvertDB2Config
func (s *GiftService) ConvertDB2Config(ol *model.GiftOnline) (res *v1pb.GiftConfigResp_Config) {
	coinType := "silver"
	bagGift := int64(0)
	countMap := []*v1pb.GiftConfigResp_CountMap{
		{Num: 1, Text: ""}, {Num: 10, Text: ""}, {Num: 99, Text: ""}, {Num: 520, Text: ""},
	}
	switch ol.CoinType {
	case coinTypeGold:
		coinType = "gold"
		countMap = []*v1pb.GiftConfigResp_CountMap{
			{Num: 1, Text: ""},
		}
	case coinTypeBag:
		bagGift = 1
	}
	cm := map[string]string{
		"BROADCAST": "广播",
		"ACTIVITY":  "活动",
	}
	cornerMark, ok := cm[ol.CornerMark]
	if !ok {
		cornerMark = ol.CornerMark
	}
	stayTime := s.GetStayTime(ol)

	res = &v1pb.GiftConfigResp_Config{
		Id:                   ol.GiftId,
		Name:                 ol.Name,
		Price:                ol.Price,
		Type:                 ol.Type,
		CoinType:             coinType,
		BagGift:              bagGift,
		Effect:               ol.Effect,
		CornerMark:           cornerMark,
		Broadcast:            ol.Broadcast,
		Draw:                 ol.Draw,
		StayTime:             stayTime,
		AnimationFrameNum:    ol.AnimationFrameNum,
		Desc:                 ol.Desc,
		Rule:                 ol.Rule,
		Rights:               ol.Rights,
		PrivilegeRequired:    ol.PrivilegeRequired,
		CountMap:             countMap,
		ImgBasic:             strings.Replace(ol.AssetImgBasic, "i0.hdslb.com", "s1.hdslb.com", 1),
		ImgDynamic:           ol.AssetImgDynamic,
		FrameAnimation:       ol.AssetFrameAnimation,
		Gif:                  ol.AssetGif,
		Webp:                 ol.AssetWebp,
		FullScWeb:            ol.AssetFullScWeb,
		FullScHorizontal:     ol.AssetFullScHorizontal,
		FullScVertical:       ol.AssetFullScVertical,
		FullScHorizontalSvga: ol.AssetFullScHorizontalSvga,
		FullScVerticalSvga:   ol.AssetFullScVerticalSvga,
		BulletHead:           ol.AssetBulletHead,
		BulletTail:           ol.AssetBulletTail,
	}
	return
}

// GetStayTime 获取礼物stayTime
func (s *GiftService) GetStayTime(ol *model.GiftOnline) (time int64) {
	giftIDMap := map[int64]int64{
		25: 30,
	}
	typeMap := map[int64]int64{
		0: 3,
		1: 2,
		2: 20,
		3: 6,
	}
	var ok bool
	if time, ok = giftIDMap[ol.GiftId]; !ok {
		if time, ok = typeMap[ol.Type]; !ok {
			time = 3
		}
	}
	return
}

// IsValidGift 判断道具合法性
func (s *GiftService) IsValidGift(giftID int64) bool {
	giftCache = giftAtomic.Load().(*GiftCache)
	giftIDStrArr := strings.Split(giftCache.validGift, ",")
	for _, giftIDStr := range giftIDStrArr {
		validGiftID, _ := strconv.ParseInt(strings.TrimSpace(giftIDStr), 10, 64)
		if giftID == validGiftID {
			return true
		}
	}
	return false
}

// IsSpecialGift 是否特殊道具
func (s *GiftService) IsSpecialGift(giftID int64) bool {
	sp := []int64{1, 3, 4, 6, 10, 25}
	for _, v := range sp {
		if giftID == v {
			return true
		}
	}
	return false
}

// GetGiftInfoByID 通过道具id获取道具信息
func (s *GiftService) GetGiftInfoByID(ctx context.Context, giftID int64) (gift *v1pb.GiftConfigResp_Config) {
	gift = &v1pb.GiftConfigResp_Config{}
	all, err := s.GetAllConfig(ctx)
	if err != nil {
		return
	}
	for _, v := range all {
		if v.Id == giftID {
			gift = v
			return
		}
	}
	return
}
