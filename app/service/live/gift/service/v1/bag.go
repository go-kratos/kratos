package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1activity "go-common/app/service/live/activity/api/liverpc/v1"
	v2fansMedal "go-common/app/service/live/fans_medal/api/liverpc/v2"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/dao"
	"go-common/app/service/live/gift/model"
	v3user "go-common/app/service/live/user/api/liverpc/v3"
	xuser "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/library/log"
	liveCtx "go-common/library/net/rpc/liverpc/context"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

var medalGiftConfig = map[int64]map[int64]int64{
	1: {1: 1}, 2: {1: 2}, 3: {1: 3}, 4: {1: 4}, 5: {1: 5}, 6: {1: 6}, 7: {1: 7},
	8: {1: 8}, 9: {1: 9}, 10: {1: 10}, 11: {1: 15}, 12: {1: 20}, 13: {1: 30}, 14: {1: 40},
	15: {1: 50}, 16: {1: 60}, 17: {1: 70}, 18: {1: 80}, 19: {1: 90}, 20: {1: 100},
}

var levelGiftConfig = map[int64]map[int64]int64{
	10: {1: 10}, 15: {1: 20}, 20: {1: 30}, 25: {1: 50}, 30: {1: 75}, 35: {1: 100}, 40: {1: 150}, 45: {1: 200}, 50: {1: 300},
}

// 包裹发送状态
const (
	needSend int64 = 1
	hasSent  int64 = 2
)

const (
	medalDailyBag int64 = 1
	levelWeekBag  int64 = 2
	vipMonthBag   int64 = 3
)

//GetUserInfo GetUserInfo
func (s *GiftService) GetUserInfo(c context.Context, uid int64) (res map[int64]*v3user.UserGetMultipleResp_UserInfo, err error) {
	res = make(map[int64]*v3user.UserGetMultipleResp_UserInfo)
	if uid == 0 {
		err = errors.New("nil uid")
		return
	}
	reply, err := dao.UserApi.V3User.GetMultiple(c, &v3user.UserGetMultipleReq{
		Uids:       []int64{uid},
		Attributes: []string{"exp"},
	})
	if err != nil {
		log.Error("getUserInfo,uid:%d, err:%v", uid, err)
		return
	}
	if reply.Code != 0 {
		log.Error("user_getUserInfo_error:%d,%d,%s,%v,uid", uid, reply.Code, reply.Msg, reply.Data)
		return
	}
	res = reply.Data
	return
}

func getBagLockKey(uid int64) string {
	return fmt.Sprintf("gift:daily_bag.%d", uid)
}

// GetMedalGift 获取勋章礼包
func (s *GiftService) GetMedalGift(ctx context.Context, uid int64) (bag *v1pb.DailyBagResp_BagList, err error) {
	medal, giftConfig, err := s.GetMedalGiftConfig(ctx, uid)
	if err != nil || giftConfig == nil {
		return
	}
	bgs := s.GetDailyGiftBag(ctx, uid, giftConfig)
	if bgs.Status == needSend {
		bag = &v1pb.DailyBagResp_BagList{
			Type:    medalDailyBag,
			BagName: "粉丝勋章礼包",
			Source: &v1pb.DailyBagResp_BagList_Source{
				MedalId:   medal.MedalId,
				MedalName: medal.MedalName,
				Level:     medal.Level,
			},
		}
		expireAt := getDeltaDayTime(1)
		for _, v := range bgs.Gift {
			bag.GiftList = append(bag.GiftList, &v1pb.DailyBagResp_BagList_GiftList{
				GiftId:   strconv.FormatInt(v.GiftID, 10),
				GiftNum:  v.GiftNum,
				ExpireAt: expireAt,
			})
		}
	}
	return
}

// GetMedalGiftConfig 获取勋章礼物配置
func (s *GiftService) GetMedalGiftConfig(ctx context.Context, uid int64) (medal *v2fansMedal.HighQpsLiveReceivedResp_Data, giftConfig map[int64]int64, err error) {
	medal = &v2fansMedal.HighQpsLiveReceivedResp_Data{}
	timeout := time.Duration(200) * time.Millisecond
	ctx = liveCtx.WithTimeout(ctx, timeout)
	resp, err := dao.FansMedalApi.V2HighQps.LiveReceived(ctx, &v2fansMedal.HighQpsLiveReceivedReq{Uid: uid})
	if err != nil {
		log.Error("get medal list error:%v", err)
		return
	}
	medalList := resp.Data
	if len(medalList) == 0 {
		return
	}
	var medalLevel int64
	for _, v := range medalList {
		if v.Level >= medalLevel {
			medalLevel = v.Level
			medal = v
		}
	}
	log.Info("daily_bag_info,uid:%d,medalLevel:%d", uid, medalLevel)

	var ok bool
	giftConfig, ok = medalGiftConfig[medalLevel]
	if !ok {
		return medal, nil, nil
	}
	return
}

// GetDailyGiftBag 获取勋章礼包
func (s *GiftService) GetDailyGiftBag(ctx context.Context, uid int64, giftConfig map[int64]int64) (data *model.BagGiftStatus) {
	data, err := s.dao.GetMedalDailyBagCache(ctx, uid)
	if err != nil {
		return
	}
	if data != nil {
		return
	}
	wg := sync.WaitGroup{}
	date := time.Now().Format("2006-01-02 00:00:00")
	res, err := s.dao.GetDayBagStatus(ctx, uid, date)
	if err != nil {
		return
	}
	curTime := time.Now().Unix()
	expireAt := getDeltaDayTime(1)
	if res.ID != 0 {
		dayInfo := res.DayInfo
		data = &model.BagGiftStatus{}
		err = json.Unmarshal([]byte(dayInfo), data)
		if err != nil {
			log.Error("unmarshal err :%v", err)
			return
		}
		data.Status = hasSent
	} else {
		//发送礼物
		giftCfg := s.FormatGift(giftConfig, expireAt)
		data = &model.BagGiftStatus{
			Status: needSend, Gift: giftCfg,
		}
		s.dao.AddDayBag(ctx, uid, date, data)
		for _, v := range giftCfg {
			wg.Add(1)
			go func(v *model.GiftInfo) {
				defer wg.Done()
				s.SendAddGiftMsg(ctx, uid, v.GiftID, v.GiftNum, expireAt, "daily_bag_"+date, uuid.NewV4().String())
			}(v)
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.dao.SetMedalDailyBagCache(ctx, uid, data, expireAt-curTime)
	}()
	wg.Wait()
	return
}

func getDeltaDayTime(day int) int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.AddDate(0, 0, day).Unix()
}

func getThisWeekTime() int64 {
	weekDay := int64(time.Now().Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	return getDeltaDayTime(1) - weekDay*86400
}

func getNextWeekTime() int64 {
	return getThisWeekTime() + 7*86400
}

// FormatGift FormatGift
func (s *GiftService) FormatGift(giftConfig map[int64]int64, expireAt int64) (gift []*model.GiftInfo) {
	gift = make([]*model.GiftInfo, 0)
	for k, v := range giftConfig {
		g := &model.GiftInfo{}
		g.GiftID = k
		g.GiftNum = v
		var liveTime = float64(expireAt - time.Now().Unix())
		var day float64 = 86400
		if liveTime < day {
			g.ExpireAt = "今天"
		} else {
			g.ExpireAt = fmt.Sprintf("%.0f天", math.Ceil(liveTime/day))
		}
		gift = append(gift, g)
	}
	return
}

// GetUserLevelInfo grpc获取用户经验信息
func (s *GiftService) GetUserLevelInfo(ctx context.Context, uid int64) (levelInfo *xuser.UserLevelInfo) {
	levelInfo = &xuser.UserLevelInfo{}
	uids := []int64{uid}
	resp, err := dao.XuserClient.UserExpClient.GetUserExp(ctx, &xuser.GetUserExpReq{Uids: uids})
	if err != nil {
		log.Error("call xuser get exp err,%v", err)
		return
	}
	return resp.Data[uid].UserLevel
}

// GetLevelGift 获取用户等级礼包
func (s *GiftService) GetLevelGift(ctx context.Context, uid int64) (bag *v1pb.DailyBagResp_BagList) {
	bag = &v1pb.DailyBagResp_BagList{}
	userInfo := s.GetUserLevelInfo(ctx, uid)
	level := userInfo.Level
	log.Info("daily_bag_info,uid:%d,userLevel:%d", uid, level)
	cfg := s.GetLevelGiftConfig(level)
	if len(cfg) == 0 {
		return
	}
	bgs := s.GetWeekLevelGiftBag(ctx, uid, level, cfg)
	if bgs.Status == needSend {
		bag = &v1pb.DailyBagResp_BagList{
			Type:    levelWeekBag,
			BagName: "用户等级礼包",
			Source: &v1pb.DailyBagResp_BagList_Source{
				UserLevel: level,
			},
		}
		expireAt := getNextWeekTime() + time.Now().Unix() - getDeltaDayTime(0)
		for _, v := range bgs.Gift {
			bag.GiftList = append(bag.GiftList, &v1pb.DailyBagResp_BagList_GiftList{
				GiftId:   strconv.FormatInt(v.GiftID, 10),
				GiftNum:  v.GiftNum,
				ExpireAt: expireAt,
			})
		}
	}
	return
}

// GetLevelGiftConfig 获取周等级礼包礼物配置
func (s *GiftService) GetLevelGiftConfig(level int64) (r map[int64]int64) {
	var flag int64 = 100
	for l, v := range levelGiftConfig {
		gap := level - l
		if gap < 0 {
			continue
		}
		if gap < flag {
			flag = gap
			r = v
		}
	}
	return
}

// GetWeekLevelGiftBag 获取周等级礼包
func (s *GiftService) GetWeekLevelGiftBag(ctx context.Context, uid, level int64, giftConfig map[int64]int64) (data *model.BagGiftStatus) {
	data, err := s.dao.GetWeekLevelBagCache(ctx, uid, level)
	if err != nil {
		return
	}
	if data != nil {
		return
	}
	wg := sync.WaitGroup{}
	_, week := time.Now().ISOWeek()
	res, err := s.dao.GetWeekBagStatus(ctx, uid, week, level)
	if err != nil {
		return
	}
	curTime := time.Now().Unix()
	expireAt := getNextWeekTime()
	if res.ID != 0 {
		weekInfo := res.WeekInfo
		data = &model.BagGiftStatus{}
		err = json.Unmarshal([]byte(weekInfo), data)
		if err != nil {
			log.Error("unmarshal err :%v", err)
			return
		}
		data.Status = hasSent
	} else {
		//发送礼物
		giftCfg := s.FormatGift(giftConfig, expireAt)
		data = &model.BagGiftStatus{
			Status: needSend, Gift: giftCfg,
		}
		s.dao.AddWeekBag(ctx, uid, week, level, data)
		source := "week_bag_" + strconv.Itoa(week)
		for _, v := range giftCfg {
			wg.Add(1)
			go func(v *model.GiftInfo) {
				defer wg.Done()
				s.SendAddGiftMsg(ctx, uid, v.GiftID, v.GiftNum, expireAt, source, uuid.NewV4().String())
			}(v)
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.dao.SetWeekLevelBagCache(ctx, uid, level, data, expireAt-curTime)
	}()
	wg.Wait()
	return
}

// GetVipMonthGift 获取月费老爷礼包
func (s *GiftService) GetVipMonthGift(ctx context.Context, uid int64) (bag *v1pb.DailyBagResp_BagList) {
	bag = &v1pb.DailyBagResp_BagList{}
	status, err := s.dao.GetVipStatusCache(ctx, uid)
	log.Info("daily_bag_info,uid:%d,vip_status:%d", uid, status)
	if err != nil {
		return
	}
	if status == 1 {
		bag = &v1pb.DailyBagResp_BagList{
			Type:    vipMonthBag,
			BagName: "年费老爷每月礼包",
			GiftList: []*v1pb.DailyBagResp_BagList_GiftList{
				{GiftId: "1", GiftNum: 99}, {GiftId: "3", GiftNum: 1}, {GiftId: "egg", GiftNum: 5},
			},
		}
		s.dao.ClearVipStatusCache(ctx, uid)
	}
	return
}

// GetUnionFansGift 获取友爱社礼包
func (s *GiftService) GetUnionFansGift(ctx context.Context, uid int64) (bag []*v1pb.DailyBagResp_BagList) {
	bag = make([]*v1pb.DailyBagResp_BagList, 0)
	unionResp := s.GetUnionConfig(ctx, uid)
	if len(unionResp) != 0 {
		for _, v := range unionResp {
			tmp := &v1pb.DailyBagResp_BagList{
				Type:    v.Type,
				BagName: v.GiftTypeName,
			}
			for _, u := range v.GiftList {
				tmp.GiftList = append(tmp.GiftList, &v1pb.DailyBagResp_BagList_GiftList{
					GiftId:  u.GiftId,
					GiftNum: u.GiftNum,
				})
			}
			bag = append(bag, tmp)
		}
	}
	return
}

// GetUnionConfig 获取友爱社礼物配置
func (s *GiftService) GetUnionConfig(ctx context.Context, uid int64) (res []*v1activity.UnionFansGetSendGiftResp_Data) {
	res = make([]*v1activity.UnionFansGetSendGiftResp_Data, 0)
	timeout := time.Duration(200) * time.Millisecond
	ctx = liveCtx.WithTimeout(ctx, timeout)
	reply, err := dao.ActivityApi.V1UnionFans.GetSendGift(ctx, &v1activity.UnionFansGetSendGiftReq{
		Uid: uid,
	})
	if err != nil {
		log.Error("GetSendGift err:%v", err)
		return
	}
	if reply.Code != 0 {
		log.Error("GetSendGifterror:%d,%s,%v", reply.Code, reply.Msg, reply.Data)
		return
	}
	res = reply.Data
	log.Info("daily_bag_info,uid:%d,unionConfig:%v", uid, res)
	return
}

// GetBagExpireStatus 包裹过期状态，0没有7日内过期的，1有7日内过期的，-1包裹为空,-100表示cache miss
func (s *GiftService) GetBagExpireStatus(ctx context.Context, uid int64) (status int64) {
	status, err := s.dao.GetBagStatusCache(ctx, uid)
	if err != nil {
		return
	}
	if status != -100 {
		return
	}
	status = -1
	bagList := s.GetBagList(ctx, uid)
	now := time.Now().Unix()
	for _, v := range bagList {
		giftInfo := s.GetGiftInfoByID(ctx, v.GiftID)
		if giftInfo.Id == 0 {
			continue
		}
		if (v.ExpireAt > 0 && v.ExpireAt < now) || v.GiftNum < 1 {
			continue
		}
		if v.ExpireAt > 0 && (v.ExpireAt <= now+7*86400) {
			status = 1
			break
		}
	}
	s.dao.SetBagStatusCache(ctx, uid, status, 30)
	return
}

// GetBagList GetBagList from cache and DB
func (s *GiftService) GetBagList(ctx context.Context, uid int64) (bagList []*model.BagGiftList) {
	bagList, err := s.dao.GetBagListCache(ctx, uid)
	if err != nil {
		return
	}
	if bagList != nil {
		return
	}
	// cache nil then get from db
	bagList, err = s.dao.GetBagList(ctx, uid)
	if err != nil {
		return
	}
	s.dao.SetBagListCache(ctx, uid, bagList, 3600)
	return
}
