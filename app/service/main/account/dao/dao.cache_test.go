package dao

import (
	"context"
	"testing"

	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInfo(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1405)
	)
	convey.Convey("Get base-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Info(c, id)
		ctx.Convey("Then err should be nil and res should be not nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInfos(t *testing.T) {
	var (
		c    = context.TODO()
		keys = []int64{110020171, 110019841}
	)
	convey.Convey("Batch get base-infos from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Infos(c, keys)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCard(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(110020171)
	)
	convey.Convey("Get card-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Card(c, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCards(t *testing.T) {
	var (
		c    = context.TODO()
		keys = []int64{110020171, 110019841}
	)
	convey.Convey("Batch get card-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Cards(c, keys)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVip(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(110018881)
	)
	convey.Convey("Get vip-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Vip(c, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAutoRenewVip(t *testing.T) {
	var (
		c            = context.TODO()
		autoRenewMid = int64(27515232)
	)
	convey.Convey("Get vip-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Vip(c, autoRenewMid)
		t.Logf("data(%+v)", res)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVips(t *testing.T) {
	var (
		c    = context.TODO()
		keys = []int64{1405, 2205}
	)
	convey.Convey("Batch get vip-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Vips(c, keys)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoProfile(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(111003471)
	)
	convey.Convey("Get profile-info from cache if miss will call source method", t, func(ctx convey.C) {
		res, err := d.Profile(c, id)
		if err == ecode.Degrade {
			err = nil
		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
