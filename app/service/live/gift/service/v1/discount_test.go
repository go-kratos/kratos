package v1

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV1GetHighestGuardLevel(t *testing.T) {
	convey.Convey("GetHighestGuardLevel", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			lv, err := s.GetHighestGuardLevel(ctx, uid)
			c.Convey("Then err should be nil.lv should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(lv, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetDiscountList(t *testing.T) {
	convey.Convey("GetDiscountList", t, func(c convey.C) {
		var (
			ctx          = context.Background()
			platform     = ""
			userType     = int64(0)
			roomID       = int64(1)
			areaID       = int64(1)
			areaParentID = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			config, err := s.GetDiscountList(ctx, platform, userType, roomID, areaID, areaParentID)
			c.Convey("Then err should be nil.config should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(config, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetGiftPrice(t *testing.T) {
	convey.Convey("GetGiftPrice", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			giftID = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			price := s.GetGiftPrice(ctx, giftID)
			c.Convey("Then price should not be nil.", func(c convey.C) {
				c.So(price, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1LoadDiscountCache(t *testing.T) {
	convey.Convey("LoadDiscountCache", t, func(c convey.C) {
		var (
			ctx   = context.Background()
			force bool
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := s.LoadDiscountCache(ctx, force)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestV1SyncDiscountCache(t *testing.T) {
	convey.Convey("SyncDiscountCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := s.SyncDiscountCache(ctx)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestV1parsePlatform2Slice(t *testing.T) {
	convey.Convey("parsePlatform2Slice", t, func(c convey.C) {
		pf := s.parsePlatform2Slice(0)
		c.So(pf, convey.ShouldNotBeNil)
		c.So(pf, convey.ShouldContain, int64(0))

		pf = s.parsePlatform2Slice(6)
		c.So(pf, convey.ShouldNotBeNil)
		c.So(pf, convey.ShouldHaveLength, 2)
		c.So(pf, convey.ShouldContain, int64(2))
		c.So(pf, convey.ShouldContain, int64(4))

	})
}

func TestV1GetDiscountPlans(t *testing.T) {
	convey.Convey("GetDiscountPlans", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			dp, err := s.GetDiscountPlans(ctx)
			c.Convey("Then err should be nil.dp should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(dp, convey.ShouldNotBeNil)
			})
		})
	})
}
