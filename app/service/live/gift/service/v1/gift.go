package v1

import (
	"context"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/conf"
	"go-common/app/service/live/gift/dao"
	v1room "go-common/app/service/live/room/api/liverpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"
	"sync"
	"time"
)

// GiftService struct
type GiftService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	dao     *dao.Dao
	addGift *databus.Databus
}

//NewGiftService init
func NewGiftService(c *conf.Config) (s *GiftService) {
	s = &GiftService{
		conf:    c,
		dao:     dao.New(c),
		addGift: databus.New(c.Databus.AddGift),
	}
	go s.infocproc()
	return s
}

// RoomGiftList implementation
func (s *GiftService) RoomGiftList(ctx context.Context, req *v1pb.RoomGiftListReq) (resp *v1pb.RoomGiftListResp, err error) {
	resp, err = s.GetOnlinePlanGiftList(ctx, req.RoomId, req.AreaV2ParentId, req.AreaV2Id, req.Platform, req.Build, req.MobiApp)
	if err != nil {
		return
	}
	resp.ShowCountMap = 0
	resp.OldList = s.getOldList()
	return
}

// GiftConfig implementation
func (s *GiftService) GiftConfig(ctx context.Context, req *v1pb.GiftConfigReq) (resp *v1pb.GiftConfigResp, err error) {
	resp = &v1pb.GiftConfigResp{
		Data: make([]*v1pb.GiftConfigResp_Config, 0),
	}
	all, err := s.GetAllConfig(ctx)
	if err != nil || len(all) == 0 {
		log.Error("all gift num:%d", len(all))
		return
	}
	resp.Data = all
	enable := s.conf.Gift.EnableFilterGift
	bubbleID := s.conf.Gift.BubbleId
	if enable {
		for _, giftInfo := range all {
			giftID := giftInfo.Id
			//泡泡糖特殊逻辑
			if giftID == bubbleID {
				giftInfo.Type = 1
			}
			//kfc道具特殊逻辑
			if giftID == 30082 {
				giftInfo.Type = 2
			}
			if giftInfo.BagGift == 0 && !s.IsValidGift(giftID) && !s.IsSpecialGift(giftID) && !s.isPrivilegeGift(giftID) {
				continue
			}

			if req.Platform == "android" && req.Build != 0 && req.Build < 5270000 {
				if giftID >= 30013 && giftID <= 30038 {
					giftInfo.FullScHorizontalSvga = ""
					giftInfo.FullScVerticalSvga = ""
				}
				if giftID == 25 {
					giftInfo.FullScHorizontalSvga = "http://i0.hdslb.com/bfs/live/9f88c8260ee23b86685a8cae3efa6d9be46717a5.svga"
					giftInfo.FullScVerticalSvga = "http://i0.hdslb.com/bfs/live/ec028520a8a951553cc42dc7d841edb38daba427.svga"
				}
			}
			resp.Data = append(resp.Data, giftInfo)
		}
	}
	return
}

// TickerReloadGift 定时加载数据库数据到全局配置
func (s *GiftService) TickerReloadGift() {
	ctx := context.Background()
	s.LoadGiftCache(ctx, true)
	s.LoadDiscountCache(ctx, true)
	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for range ticker.C {
			s.LoadGiftCache(ctx, true)
			s.LoadDiscountCache(ctx, true)
		}
	}()
}

// DiscountGiftList implementation
func (s *GiftService) DiscountGiftList(ctx context.Context, req *v1pb.DiscountGiftListReq) (resp *v1pb.DiscountGiftListResp, err error) {
	resp = &v1pb.DiscountGiftListResp{}
	if req.Uid == 0 {
		err = ecode.Error(-1, "请先登录哦")
		return
	}
	roomID := req.Roomid
	areaParentID := req.AreaV2ParentId
	areaID := req.AreaV2Id
	if req.Roomid == 0 {
		var r *v1room.RoomGetStatusInfoByUidsResp
		r, err = dao.RoomApi.V1Room.GetStatusInfoByUids(ctx, &v1room.RoomGetStatusInfoByUidsReq{
			Uids:       []int64{req.Ruid},
			ShowHidden: 1,
		})
		if err != nil {
			log.Error("call V1Room.GetStatusInfoByUids error,params:uid:%v", req.Ruid)
			err = ecode.Error(-1, "内部错误")
			return
		}
		if r.Code != 0 {
			log.Error("call V1Room.GetStatusInfoByUids error,params:code:(%v),msg:(%v)", r.Code, r.Msg)
			err = ecode.Error(-1, "内部错误")
			return
		}
		if len(r.Data) == 0 {
			err = ecode.Error(-1, "折扣道具获取失败")
			return
		}
		roomID = r.Data[req.Ruid].RoomId
		areaParentID = r.Data[req.Ruid].AreaV2ParentId
		areaID = r.Data[req.Ruid].AreaV2Id
	}

	userType, err := s.GetHighestGuardLevel(ctx, req.Uid)
	if err != nil {
		err = ecode.Error(-1, "折扣道具获取失败")
		return
	}

	list, err := s.GetDiscountList(ctx, req.Platform, userType, roomID, areaParentID, areaID)
	if err != nil {
		return
	}
	resp.DiscountList = list
	return
}

// DailyBag implementation
func (s *GiftService) DailyBag(ctx context.Context, req *v1pb.DailyBagReq) (resp *v1pb.DailyBagResp, err error) {
	resp = &v1pb.DailyBagResp{}
	if req.Uid == 0 {
		err = ecode.Error(400, "请先登录哦")
		return
	}
	uid := req.Uid
	res, err := s.dao.GetDailyBagCache(ctx, uid)
	if err != nil {
		err = ecode.Error(-1, "系统错误")
		return
	}
	if res != nil {
		resp = &v1pb.DailyBagResp{
			BagStatus:       2,
			BagExpireStatus: s.GetBagExpireStatus(ctx, uid),
			BagToast: &v1pb.DailyBagResp_BagToast{
				ToastStatus:  0,
				ToastMessage: "",
			},
			BagList: res,
		}
		return
	}

	gotLock, _, err := s.dao.Lock(ctx, getBagLockKey(uid), 5000, 0, 0)
	if !gotLock || err != nil {
		err = ecode.Error(400, "包裹领取失败，请重试")
		return
	}

	eg, _ := errgroup.WithContext(ctx)
	//同时获取各种礼包返回
	resp.BagList = make([]*v1pb.DailyBagResp_BagList, 0)
	var mu sync.Mutex
	eg.Go(func() error {
		bag, _ := s.GetMedalGift(ctx, uid)
		if bag.Type > 0 {
			mu.Lock()
			resp.BagList = append(resp.BagList, bag)
			mu.Unlock()
		}
		return nil
	})
	eg.Go(func() error {
		bag := s.GetLevelGift(ctx, uid)
		if bag.Type > 0 {
			mu.Lock()
			resp.BagList = append(resp.BagList, bag)
			mu.Unlock()
		}
		return nil
	})
	eg.Go(func() error {
		bag := s.GetVipMonthGift(ctx, uid)
		if bag.Type > 0 {
			mu.Lock()
			resp.BagList = append(resp.BagList, bag)
			mu.Unlock()
		}
		return nil
	})
	eg.Go(func() error {
		unionGift := s.GetUnionFansGift(ctx, uid)
		mu.Lock()
		resp.BagList = append(resp.BagList, unionGift...)
		mu.Unlock()
		return nil
	})
	eg.Wait()

	// 同时进行后续清理
	eg.Go(func() error {
		return s.dao.ForceUnLock(ctx, getBagLockKey(uid))
	})
	var expireStatus int64 = -1
	eg.Go(func() error {
		expireStatus = s.GetBagExpireStatus(ctx, uid)
		return nil
	})
	eg.Go(func() error {
		return s.dao.SetDailyBagCache(ctx, uid, resp.BagList, 3600)
	})
	eg.Wait()

	resp.BagStatus = 2
	if len(resp.BagList) > 0 {
		resp.BagStatus = 1
	}
	resp.BagExpireStatus = expireStatus
	resp.BagToast = &v1pb.DailyBagResp_BagToast{
		ToastStatus:  0,
		ToastMessage: "",
	}

	return
}
