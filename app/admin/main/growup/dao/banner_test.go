package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTotalBannerCount(t *testing.T) {
	convey.Convey("TotalBannerCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.TotalBannerCount(c)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDupEditBanner(t *testing.T) {
	convey.Convey("DupEditBanner", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startAt = int64(0)
			endAt   = int64(0)
			now     = int64(0)
			id      = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dup, err := d.DupEditBanner(c, startAt, endAt, now, id)
			ctx.Convey("Then err should be nil.dup should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dup, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDupBanner(t *testing.T) {
	convey.Convey("DupBanner", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			startAt = int64(0)
			endAt   = int64(0)
			now     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dup, err := d.DupBanner(c, startAt, endAt, now)
			ctx.Convey("Then err should be nil.dup should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dup, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertBanner(t *testing.T) {
	convey.Convey("InsertBanner", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			image   = ""
			link    = ""
			startAt = int64(0)
			endAt   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertBanner(c, image, link, startAt, endAt)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBanners(t *testing.T) {
	convey.Convey("Banners", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			offset = int64(0)
			limit  = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Banners(c, offset, limit)
			ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateBanner(t *testing.T) {
	convey.Convey("UpdateBanner", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			image   = ""
			link    = ""
			startAt = int64(0)
			endAt   = int64(0)
			id      = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpdateBanner(c, image, link, startAt, endAt, id)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBannerEndAt(t *testing.T) {
	convey.Convey("UpdateBannerEndAt", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			endAt = int64(0)
			id    = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpdateBannerEndAt(c, endAt, id)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
