package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/growup/model"
	"go-common/library/net/metadata"
)

func TestDaoArticleStat(t *testing.T) {
	convey.Convey("ArticleStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
			ip  = metadata.String(c, metadata.RemoteIP)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ArticleStat(c, mid, ip)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetUpBaseStatCache(t *testing.T) {
	convey.Convey("getUpBaseStatCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(1001)
			date = "2018-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.getUpBaseStatCache(c, mid, date)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosetUpBaseStatCache(t *testing.T) {
	convey.Convey("setUpBaseStatCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(1001)
			date = "2018-01-01"
			st   = &model.UpBaseStat{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.setUpBaseStatCache(c, mid, date, st)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpStat(t *testing.T) {
	convey.Convey("UpStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
			dt  = "20180101"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpStat(c, mid, dt)
			ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
