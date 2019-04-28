package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBindInfoByMid(t *testing.T) {
	convey.Convey("BindInfoByMid", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			id    = int64(1)
			appID = int64(32)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.BindInfoByMid(c, id, appID)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOpenInfoByOpenID(t *testing.T) {
	convey.Convey("OpenInfoByOpenID", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			id    = "bdca8b71e7a6726885d40a395bf9ccd1"
			appID = int64(32)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.OpenInfoByOpenID(c, id, appID)
			fmt.Println("res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
