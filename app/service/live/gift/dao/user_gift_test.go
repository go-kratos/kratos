package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetBag(t *testing.T) {
	convey.Convey("GetBag", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(0)
			giftID   = int64(0)
			expireAt = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetBag(ctx, uid, giftID, expireAt)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBagNum(t *testing.T) {
	convey.Convey("UpdateBagNum", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
			id  = int64(0)
			num = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			affected, err := d.UpdateBagNum(ctx, uid, id, num)
			c.Convey("Then err should be nil.affected should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBag(t *testing.T) {
	convey.Convey("AddBag", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(0)
			giftID   = int64(0)
			giftNum  = int64(0)
			expireAt = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			affected, err := d.AddBag(ctx, uid, giftID, giftNum, expireAt)
			c.Convey("Then err should be nil.affected should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetBagByID(t *testing.T) {
	convey.Convey("GetBagByID", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
			id  = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetBagByID(ctx, uid, id)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetPostFix(t *testing.T) {
	convey.Convey("getPostFix", t, func(c convey.C) {
		var (
			uid = int64(1)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := getPostFix(uid)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
				c.So(p1, convey.ShouldEqual, "c")
			})
		})
	})
}

func TestDaoGetBagList(t *testing.T) {
	convey.Convey("GetBagList", t, func(c convey.C) {
		var (
			ctx = context.Background()
			uid = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			list, err := d.GetBagList(ctx, uid)
			c.Convey("Then err should be nil.list should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(list, convey.ShouldBeNil)
			})
		})
	})
}
