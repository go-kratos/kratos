package service

import (
	"context"
	"go-common/app/job/live/gift/internal/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceAddGift(t *testing.T) {
	convey.Convey("AddGift", t, func(c convey.C) {
		var (
			ctx = context.Background()
			m   = &model.AddFreeGift{
				UID:      1,
				GiftID:   1,
				GiftNum:  1,
				ExpireAt: 0,
				Source:   "test",
			}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bagId, err := s.AddGift(ctx, m)
			c.Convey("Then err should be nil.bagId should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(bagId, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceGetBagID(t *testing.T) {
	convey.Convey("GetBagID", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(1)
			giftID   = int64(1)
			expireAt = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			id, err := s.GetBagID(ctx, uid, giftID, expireAt)
			c.Convey("Then err should be nil.id should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateFreeGiftCache(t *testing.T) {
	convey.Convey("UpdateFreeGiftCache", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(0)
			giftID   = int64(0)
			expireAt = int64(0)
			num      = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			s.UpdateFreeGiftCache(ctx, uid, giftID, expireAt, num)
			c.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
