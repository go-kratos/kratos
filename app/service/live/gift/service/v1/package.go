package v1

import (
	"context"
	"errors"
	"go-common/app/service/live/gift/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// AddGift 发送免费道具
func (s *GiftService) AddGift(ctx context.Context, uid int64, giftID, giftNum, expireAt int64, source string) (bagId int64, err error) {
	if uid == 0 || giftID == 0 || giftNum == 0 {
		log.Error("add gift params error,uid:%d,giftID:%d,giftNum:%d", uid, giftID, giftNum)
		err = errors.New("params error")
		return
	}
	bagID, err := s.GetBagID(ctx, uid, giftID, expireAt)
	if err != nil {
		return
	}
	var (
		affectNum int64
		isUpdate  = false
		eg, _     = errgroup.WithContext(ctx)
	)
	if bagID != 0 {
		isUpdate = true
		affectNum, _ = s.dao.UpdateBagNum(ctx, uid, bagID, giftNum)
	} else {
		affectNum, _ = s.dao.AddBag(ctx, uid, giftID, giftNum, expireAt)
		bagID = affectNum
		eg.Go(
			func() error {
				s.dao.SetBagIDCache(ctx, uid, giftID, expireAt, bagID, 14400)
				return nil
			})

	}
	newNum := giftNum
	if affectNum > 0 {
		eg.Go(
			func() error {
				s.dao.ClearBagListCache(ctx, uid)
				return nil
			})
		if isUpdate {
			res, _ := s.dao.GetBagByID(ctx, uid, bagID)
			newNum = res.GiftNum
			//上报lancer TODO
			s.bagLogInfoc(uid, bagID, giftID, giftNum, newNum, source)
		}
	}
	// 更新免费礼物数量缓存
	eg.Go(
		func() error {
			s.UpdateFreeGiftCache(ctx, uid, giftID, expireAt, newNum)
			return nil
		})
	eg.Wait()
	return
}

// GetBagID 获取包裹id
func (s *GiftService) GetBagID(ctx context.Context, uid, giftID, expireAt int64) (id int64, err error) {
	id, err = s.dao.GetBagIDCache(ctx, uid, giftID, expireAt)
	if err != nil {
		return
	}
	if id == 0 {
		//queryDB
		var r *model.BagInfo
		r, err = s.dao.GetBag(ctx, uid, giftID, expireAt)
		if err != nil {
			return
		}
		id = r.ID
	}
	// 缓存或数据库本身有，再更新缓存
	if id != 0 {
		s.dao.SetBagIDCache(ctx, uid, giftID, expireAt, id, 14400)
	}
	return
}

// UpdateFreeGiftCache UpdateFreeGiftCache
func (s *GiftService) UpdateFreeGiftCache(ctx context.Context, uid, giftID, expireAt, num int64) {
	giftInfo := s.GetGiftInfoByID(ctx, giftID)
	if giftInfo.Id == 0 || giftInfo.Type != 3 {
		return
	}
	s.dao.SetBagNumCache(ctx, uid, giftID, expireAt, num, 14400)
}
