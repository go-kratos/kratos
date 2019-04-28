package manager

import (
	"context"
	upgrpc "go-common/app/service/main/up/api/v1"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerAddCacheUpSpecial(t *testing.T) {
	convey.Convey("AddCacheUpSpecial", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &upgrpc.UpSpecial{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheUpSpecial(c, id, val)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestManagerCacheUpSpecial(t *testing.T) {
	convey.Convey("CacheUpSpecial", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CacheUpSpecial(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestManagerDelCacheUpSpecial(t *testing.T) {
	convey.Convey("DelCacheUpSpecial", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCacheUpSpecial(c, id)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestManagerAddCacheUpsSpecial(t *testing.T) {
	convey.Convey("AddCacheUpsSpecial", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			values map[int64]*upgrpc.UpSpecial
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheUpsSpecial(c, values)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestManagerCacheUpsSpecial(t *testing.T) {
	convey.Convey("CacheUpsSpecial", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{27515256}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CacheUpsSpecial(c, ids)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestManagerDelCacheUpsSpecial(t *testing.T) {
	convey.Convey("DelCacheUpsSpecial", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCacheUpsSpecial(c, ids)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
