package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/card/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawEquip(t *testing.T) {
	convey.Convey("RawEquip", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			r, err := d.RawEquip(c, mid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawEquips(t *testing.T) {
	convey.Convey("RawEquips", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2, 3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawEquips(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoEffectiveCards(t *testing.T) {
	convey.Convey("EffectiveCards", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.EffectiveCards(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoEffectiveGroups(t *testing.T) {
	convey.Convey("EffectiveGroups", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.EffectiveGroups(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCardEquip(t *testing.T) {
	convey.Convey("CardEquip", t, func(ctx convey.C) {
		var (
			c = context.Background()
			e = &model.UserEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.CardEquip(c, e)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDeleteEquip(t *testing.T) {
	convey.Convey("DeleteEquip", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteEquip(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
