package stat

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestStatUpdateFav(t *testing.T) {
	convey.Convey("UpdateFav", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(111)
			count = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rows, err := d.UpdateFav(c, oid, count)
			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatUpdateShare(t *testing.T) {
	convey.Convey("UpdateShare", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(111)
			count = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rows, err := d.UpdateShare(c, oid, count)
			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatUpdatePlay(t *testing.T) {
	convey.Convey("UpdatePlay", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(111)
			count = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rows, err := d.UpdatePlay(c, oid, count)
			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatStat(t *testing.T) {
	convey.Convey("Stat", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			oid = int64(111)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			f, err := d.Stat(c, oid)
			convCtx.Convey("Then err should be nil.f should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(f, convey.ShouldNotBeNil)
			})
		})
	})
}
