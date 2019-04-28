package dao

import (
	"context"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaodailyBagKey(t *testing.T) {
	convey.Convey("dailyBagKey", t, func(c convey.C) {
		var (
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := dailyBagKey(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetDailyBagCache(t *testing.T) {
	convey.Convey("GetDailyBagCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetDailyBagCache(ctx, uid)
			c.Convey("Then err should be nil.res should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetDailyBagCache(t *testing.T) {
	convey.Convey("SetDailyBagCache", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			uid    = int64(9527)
			data   = []*v1pb.DailyBagResp_BagList{}
			expire = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetDailyBagCache(ctx, uid, data, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaodailyMedalBagKey(t *testing.T) {
	convey.Convey("dailyMedalBagKey", t, func(c convey.C) {
		var (
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := dailyMedalBagKey(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetMedalDailyBagCache(t *testing.T) {
	convey.Convey("GetMedalDailyBagCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(-1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetMedalDailyBagCache(ctx, uid)
			c.Convey("Then err should be nil.res should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetMedalDailyBagCache(t *testing.T) {
	convey.Convey("SetMedalDailyBagCache", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			uid    = int64(9527)
			data   = &model.BagGiftStatus{}
			expire = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetMedalDailyBagCache(ctx, uid, data, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoweekLevelBagKey(t *testing.T) {
	convey.Convey("weekLevelBagKey", t, func(c convey.C) {
		var (
			uid   = int64(0)
			level = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := weekLevelBagKey(uid, level)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetWeekLevelBagCache(t *testing.T) {
	convey.Convey("GetWeekLevelBagCache", t, func(c convey.C) {
		var (
			ctx   = context.Background()
			uid   = int64(-1)
			level = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetWeekLevelBagCache(ctx, uid, level)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetWeekLevelBagCache(t *testing.T) {
	convey.Convey("SetWeekLevelBagCache", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			uid    = int64(9527)
			level  = int64(0)
			data   = &model.BagGiftStatus{}
			expire = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetWeekLevelBagCache(ctx, uid, level, data, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoLock(t *testing.T) {
	convey.Convey("Lock", t, func(c convey.C) {
		var (
			ctx        = context.Background()
			key        = "test lock"
			ttl        = int(1)
			retry      = int(0)
			retryDelay = int(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			gotLock, lockValue, err := d.Lock(ctx, key, ttl, retry, retryDelay)
			c.Convey("Then err should be nil.gotLock,lockValue should not be nil.", func(c convey.C) {
				c.So(lockValue, convey.ShouldNotBeNil)
				c.So(err, convey.ShouldBeNil)
				c.So(gotLock, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUnLock(t *testing.T) {
	convey.Convey("UnLock", t, func(c convey.C) {
		var (
			ctx       = context.Background()
			key       = "test unlock"
			lockValue = "1"
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.UnLock(ctx, key, lockValue)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoForceUnLock(t *testing.T) {
	convey.Convey("ForceUnLock", t, func(c convey.C) {
		var (
			ctx = context.Background()
			key = ""
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.ForceUnLock(ctx, key)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaolockKey(t *testing.T) {
	convey.Convey("lockKey", t, func(c convey.C) {
		var (
			key = ""
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := lockKey(key)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaorandomString(t *testing.T) {
	convey.Convey("randomString", t, func(c convey.C) {
		var (
			l = int(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := randomString(l)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaobagIDCache(t *testing.T) {
	convey.Convey("bagIDCache", t, func(c convey.C) {
		var (
			uid      = int64(0)
			giftID   = int64(0)
			expireAt = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := bagIDCache(uid, giftID, expireAt)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBagIDCache(t *testing.T) {
	convey.Convey("GetBagIDCache", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(0)
			giftID   = int64(0)
			expireAt = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bagID, err := d.GetBagIDCache(ctx, uid, giftID, expireAt)
			c.Convey("Then err should be nil.bagID should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(bagID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetBagIDCache(t *testing.T) {
	convey.Convey("SetBagIDCache", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(9527)
			giftID   = int64(1)
			expireAt = int64(1)
			bagID    = int64(1)
			expire   = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetBagIDCache(ctx, uid, giftID, expireAt, bagID, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaobagListKey(t *testing.T) {
	convey.Convey("bagListKey", t, func(c convey.C) {
		var (
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := bagListKey(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBagListCache(t *testing.T) {
	convey.Convey("GetBagListCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetBagListCache(ctx, uid)
			c.Convey("Then err should be nil.res should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetBagListCache(t *testing.T) {
	convey.Convey("SetBagListCache", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			uid    = int64(0)
			data   = []*model.BagGiftList{}
			expire = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetBagListCache(ctx, uid, data, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoClearBagListCache(t *testing.T) {
	convey.Convey("ClearBagListCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.ClearBagListCache(ctx, uid)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaobagNumKey(t *testing.T) {
	convey.Convey("bagNumKey", t, func(c convey.C) {
		var (
			uid      = int64(0)
			giftID   = int64(0)
			expireAt = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := bagNumKey(uid, giftID, expireAt)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetBagNumCache(t *testing.T) {
	convey.Convey("SetBagNumCache", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(9527)
			giftID   = int64(1)
			expireAt = int64(1)
			giftNum  = int64(1)
			expire   = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetBagNumCache(ctx, uid, giftID, expireAt, giftNum, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaovipMonthBag(t *testing.T) {
	convey.Convey("vipMonthBag", t, func(c convey.C) {
		var (
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := vipMonthBag(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetVipStatusCache(t *testing.T) {
	convey.Convey("GetVipStatusCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			status, err := d.GetVipStatusCache(ctx, uid)
			c.Convey("Then err should be nil.status should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoClearVipStatusCache(t *testing.T) {
	convey.Convey("ClearVipStatusCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.ClearVipStatusCache(ctx, uid)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaogiftBagStatus(t *testing.T) {
	convey.Convey("giftBagStatus", t, func(c convey.C) {
		var (
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := giftBagStatus(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBagStatusCache(t *testing.T) {
	convey.Convey("GetBagStatusCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			status, err := d.GetBagStatusCache(ctx, uid)
			c.Convey("Then err should be nil.status should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetBagStatusCache(t *testing.T) {
	convey.Convey("SetBagStatusCache", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			uid    = int64(9527)
			status = int64(1)
			expire = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := d.SetBagStatusCache(ctx, uid, status, expire)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}
