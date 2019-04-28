package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInfos(t *testing.T) {
	convey.Convey("Infos", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = _testMids
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			for _, v := range keys {
				err := d.DelInfoCache(c, v)
				ctx.So(err, convey.ShouldBeNil)
			}
			res, err := d.Infos(c, keys)
			fmt.Println("res", res)
			time.Sleep(100 * time.Millisecond)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			_, err = d.Infos(c, keys)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoInfo(t *testing.T) {
	convey.Convey("Info", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = _testMid
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelInfoCache(c, _testMid)
			ctx.So(err, convey.ShouldBeNil)
			res, err := d.Info(c, id)
			time.Sleep(100 * time.Millisecond)
			fmt.Println("res", res.VipType, res.VipStatus)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			res, err = d.Info(c, id)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
