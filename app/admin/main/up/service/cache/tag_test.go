package cache

import (
	"go-common/app/interface/main/creative/model/tag"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestCacheClearTagCache(t *testing.T) {
	convey.Convey("ClearTagCache", t, func(ctx convey.C) {
		var (
			tm = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ClearTagCache(tm)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestCacheAddTagCache(t *testing.T) {
	convey.Convey("AddTagCache", t, func(ctx convey.C) {
		var (
			meta = &tag.Meta{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			AddTagCache(meta)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestCacheGetTagCache(t *testing.T) {
	convey.Convey("GetTagCache", t, func(ctx convey.C) {
		var (
			ids = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, leftIDs := GetTagCache(ids)
			ctx.Convey("Then result,leftIDs should not be nil.", func(ctx convey.C) {
				ctx.So(leftIDs, convey.ShouldNotBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
