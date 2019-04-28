package v1

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV1GetUserInfo(t *testing.T) {
	convey.Convey("GetUserInfo", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(88895029)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := s.GetUserInfo(ctx, uid)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1getBagLockKey(t *testing.T) {
	convey.Convey("getBagLockKey", t, func(c convey.C) {
		var (
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := getBagLockKey(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetMedalGift(t *testing.T) {
	convey.Convey("GetMedalGift", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(88895029)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bag, err := s.GetMedalGift(ctx, uid)
			c.Convey("Then err should be nil.bag should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(bag, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetMedalGiftConfig(t *testing.T) {
	convey.Convey("GetMedalGiftConfig", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(88895029)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			medal, giftConfig, err := s.GetMedalGiftConfig(ctx, uid)
			c.Convey("Then err should be nil.medal,giftConfig should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(giftConfig, convey.ShouldNotBeNil)
				c.So(medal, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetDailyGiftBag(t *testing.T) {
	convey.Convey("GetDailyGiftBag", t, func(c convey.C) {
		var (
			ctx        = context.Background()
			uid        = int64(0)
			giftConfig map[int64]int64
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			data := s.GetDailyGiftBag(ctx, uid, giftConfig)
			c.Convey("Then data should not be nil.", func(c convey.C) {
				c.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1getDeltaDayTime(t *testing.T) {
	convey.Convey("getDeltaDayTime", t, func(c convey.C) {
		var (
			day = int(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := getDeltaDayTime(day)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1getThisWeekTime(t *testing.T) {
	convey.Convey("getThisWeekTime", t, func(c convey.C) {
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := getThisWeekTime()
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1getNextWeekTime(t *testing.T) {
	convey.Convey("getNextWeekTime", t, func(c convey.C) {
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := getNextWeekTime()
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1FormatGift(t *testing.T) {
	convey.Convey("FormatGift", t, func(c convey.C) {
		var (
			giftConfig map[int64]int64
			expireAt   = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			gift := s.FormatGift(giftConfig, expireAt)
			c.Convey("Then gift should not be nil.", func(c convey.C) {
				c.So(gift, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetLevelGift(t *testing.T) {
	convey.Convey("GetLevelGift", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bag := s.GetLevelGift(ctx, uid)
			c.Convey("Then bag should not be nil.", func(c convey.C) {
				c.So(bag, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetLevelGiftConfig(t *testing.T) {
	convey.Convey("GetLevelGiftConfig", t, func(c convey.C) {
		c.Convey("When everything gose positive", func(c convey.C) {
			r := s.GetLevelGiftConfig(9)
			c.So(r, convey.ShouldBeNil)

			r = s.GetLevelGiftConfig(20)
			c.So(r, convey.ShouldResemble, map[int64]int64{1: 30})

			r = s.GetLevelGiftConfig(26)
			c.So(r, convey.ShouldResemble, map[int64]int64{1: 50})

			r = s.GetLevelGiftConfig(60)
			c.So(r, convey.ShouldResemble, map[int64]int64{1: 300})

		})
	})
}

func TestV1GetWeekLevelGiftBag(t *testing.T) {
	convey.Convey("GetWeekLevelGiftBag", t, func(c convey.C) {
		var (
			ctx        = context.Background()
			uid        = int64(0)
			level      = int64(0)
			giftConfig map[int64]int64
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			data := s.GetWeekLevelGiftBag(ctx, uid, level, giftConfig)
			c.Convey("Then data should not be nil.", func(c convey.C) {
				c.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetVipMonthGift(t *testing.T) {
	convey.Convey("GetVipMonthGift", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bag := s.GetVipMonthGift(ctx, uid)
			c.Convey("Then bag should not be nil.", func(c convey.C) {
				c.So(bag, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetUnionFansGift(t *testing.T) {
	convey.Convey("GetUnionFansGift", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bag := s.GetUnionFansGift(ctx, uid)
			c.Convey("Then bag should not be nil.", func(c convey.C) {
				c.So(bag, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetUnionConfig(t *testing.T) {
	convey.Convey("GetUnionConfig", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res := s.GetUnionConfig(ctx, uid)
			c.Convey("Then res should not be nil.", func(c convey.C) {
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetBagExpireStatus(t *testing.T) {
	convey.Convey("GetBagExpireStatus", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(-1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			status := s.GetBagExpireStatus(ctx, uid)
			c.Convey("Then status should not be nil.", func(c convey.C) {
				c.So(status, convey.ShouldEqual, -1)
			})
		})
	})
}

func TestV1GetBagList(t *testing.T) {
	convey.Convey("GetBagList", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(123214)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			bagList := s.GetBagList(ctx, uid)
			c.Convey("Then bagList should not be nil.", func(c convey.C) {
				c.So(bagList, convey.ShouldBeNil)
			})
		})
	})
}
func TestV1GetUserLevelInfo(t *testing.T) {
	convey.Convey("GetUserLevelInfo", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(88895029)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			levelInfo := s.GetUserLevelInfo(ctx, uid)
			c.Convey("Then levelInfo should not be nil.", func(c convey.C) {
				c.So(levelInfo, convey.ShouldNotBeNil)
			})
		})
	})
}
