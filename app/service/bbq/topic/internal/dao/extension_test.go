package dao

import (
	"context"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/log"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawVideoExtension(t *testing.T) {
	convey.Convey("RawVideoExtension", t, func(convCtx convey.C) {
		var (
			ctx   = context.Background()
			svids = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.RawVideoExtension(ctx, svids)
			log.V(1).Infow(ctx, "res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCacheVideoExtension(t *testing.T) {
	convey.Convey("CacheVideoExtension", t, func(convCtx convey.C) {
		var (
			ctx   = context.Background()
			svids = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CacheVideoExtension(ctx, svids)
			log.V(1).Infow(ctx, "res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheVideoExtension(t *testing.T) {
	convey.Convey("AddCacheVideoExtension", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
		)
		extensions := make(map[int64]*api.VideoExtension)
		extensions[1] = &api.VideoExtension{Svid: 1, Extension: "{\"title_extra\":[{\"type\":1,\"name\":\"Test\",\"end\":4,\"schema\":\"qing://topic?topic_id=1\"}]}"}
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheVideoExtension(ctx, extensions)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheVideoExtension(t *testing.T) {
	convey.Convey("DelCacheVideoExtension", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			svid = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.DelCacheVideoExtension(ctx, svid)
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}

func TestDaoInsertExtension(t *testing.T) {
	convey.Convey("InsertExtension", t, func(convCtx convey.C) {
		var (
			ctx           = context.Background()
			svid          = int64(1)
			extensionType = int64(1)
			extension     = &api.Extension{TitleExtra: []*api.TitleExtraItem{{Name: "Test", Type: model.TitleExtraTypeTopic, Start: 0, End: 4, Scheme: "qing://topic?topic_id=1"}}}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rowsAffected, err := d.InsertExtension(ctx, svid, extensionType, extension)
			convCtx.Convey("Then err should be nil.rowsAffected should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rowsAffected, convey.ShouldNotBeNil)
			})
		})
	})
}
