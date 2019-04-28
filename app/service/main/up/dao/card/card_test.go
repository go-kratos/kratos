package card

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCountUpCard(t *testing.T) {
	var c = context.Background()
	convey.Convey("CountUpCard", t, func(ctx convey.C) {
		total, err := d.CountUpCard(c)
		ctx.Convey("Then err should be nil.total should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(total, convey.ShouldNotBeNil)
		})
	})
}
func TestPageListUpInfo(t *testing.T) {
	var (
		c      = context.Background()
		offset = uint(1)
		size   = uint(1)
	)
	convey.Convey("ListUpInfo", t, func(ctx convey.C) {
		infos, err := d.ListUpInfo(c, offset, size)
		ctx.Convey("Then err should be nil.infos should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(infos, convey.ShouldNotBeNil)
		})
	})
}
func TestListUpMID(t *testing.T) {
	var c = context.Background()
	convey.Convey("ListUpMID", t, func(ctx convey.C) {
		mids, err := d.ListUpMID(c)
		ctx.Convey("Then err should be nil.mids should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(mids, convey.ShouldNotBeNil)
		})
	})
}
func TestGetUpInfo(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1532165)
	)
	convey.Convey("GetUpInfo", t, func(ctx convey.C) {
		upInfo, err := d.GetUpInfo(c, mid)
		ctx.Convey("Then err should be nil.upInfo should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(upInfo, convey.ShouldNotBeNil)
		})
	})
}
func TestListUpAccount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1532165)
	)
	convey.Convey("ListUpAccount", t, func(ctx convey.C) {
		accounts, err := d.ListUpAccount(c, mid)
		ctx.Convey("Then err should be nil.accounts should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(accounts, convey.ShouldNotBeNil)
		})
	})
}
func TestListUpImage(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1532165)
	)
	convey.Convey("ListUpImage", t, func(ctx convey.C) {
		images, err := d.ListUpImage(c, mid)
		ctx.Convey("Then err should be nil.images should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(images, convey.ShouldNotBeNil)
		})
	})
}

func TestListAVID(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1532165)
	)
	convey.Convey("ListAVID", t, func(ctx convey.C) {
		avids, err := d.ListAVID(c, mid)
		ctx.Convey("Then err should be nil.avids should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(avids, convey.ShouldNotBeNil)
		})
	})
}

func TestMidUpInfoMap(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1532165}
	)
	convey.Convey("MidUpInfoMap", t, func(ctx convey.C) {
		res, err := d.MidUpInfoMap(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
func TestMidAccountsMap(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1532165}
	)
	convey.Convey("MidAccountsMap", t, func(ctx convey.C) {
		res, err := d.MidAccountsMap(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
func TestMidImagesMap(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1532165}
	)
	convey.Convey("MidImagesMap", t, func(ctx convey.C) {
		res, err := d.MidImagesMap(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
func TestMidAvidsMap(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{1532165}
	)
	convey.Convey("MidAvidsMap", t, func(ctx convey.C) {
		res, err := d.MidAvidsMap(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
