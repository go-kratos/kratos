package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/ugcpay-rank/internal/model"
	"go-common/library/cache/memcache"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCASCacheElecPrepRank(t *testing.T) {
	convey.Convey("CASCacheElecPrepRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = &model.RankElecPrepUPProto{
				Count: 233,
				UPMID: 233,
				Size_: 10,
			}
			id      int64 = 233
			rawItem *memcache.Item
			err     error
			ver     = int64(0)
		)
		_, rawItem, err = d.CacheElecPrepUPRank(c, id, ver)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rawItem, convey.ShouldNotBeNil)

		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			var ok bool
			ok, err = d.CASCacheElecPrepRank(c, val, rawItem)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})

		data, _, err := d.CacheElecPrepUPRank(c, id, ver)
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldResemble, val)
	})
}
