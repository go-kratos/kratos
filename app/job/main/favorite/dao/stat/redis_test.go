package stat

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestStatipBanKey(t *testing.T) {
	convey.Convey("ipBanKey", t, func(convCtx convey.C) {
		var (
			id = int64(0)
			ip = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key := ipBanKey(id, ip)
			convCtx.Convey("Then key should not be nil.", func(convCtx convey.C) {
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatbuvidBanKey(t *testing.T) {
	convey.Convey("buvidBanKey", t, func(convCtx convey.C) {
		var (
			id    = int64(0)
			mid   = int64(0)
			ip    = ""
			buvid = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key := buvidBanKey(id, mid, ip, buvid)
			convCtx.Convey("Then key should not be nil.", func(convCtx convey.C) {
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatIPBan(t *testing.T) {
	convey.Convey("IPBan", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
			ip = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ban := d.IPBan(c, id, ip)
			convCtx.Convey("Then ban should not be nil.", func(convCtx convey.C) {
				convCtx.So(ban, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatBuvidBan(t *testing.T) {
	convey.Convey("BuvidBan", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			mid   = int64(0)
			buvid = ""
			ip    = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ban := d.BuvidBan(c, id, mid, ip, buvid)
			convCtx.Convey("Then ban should not be nil.", func(convCtx convey.C) {
				convCtx.So(ban, convey.ShouldNotBeNil)
			})
		})
	})
}
