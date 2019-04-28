package service

import (
	bm "go-common/library/net/http/blademaster"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceBmHTTPErrorWithMsg(t *testing.T) {
	convey.Convey("BmHTTPErrorWithMsg", t, func(ctx convey.C) {
		var (
			c   = &bm.Context{}
			err error
			msg = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			BmHTTPErrorWithMsg(c, err, msg)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceBmGetStringOrDefault(t *testing.T) {
	convey.Convey("BmGetStringOrDefault", t, func(ctx convey.C) {
		var (
			c      = &bm.Context{}
			key    = ""
			defaul = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			value, exist := BmGetStringOrDefault(c, key, defaul)
			ctx.Convey("Then value,exist should not be nil.", func(ctx convey.C) {
				ctx.So(exist, convey.ShouldNotBeNil)
				ctx.So(value, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceBmGetInt64OrDefault(t *testing.T) {
	convey.Convey("BmGetInt64OrDefault", t, func(ctx convey.C) {
		var (
			c      = &bm.Context{}
			key    = ""
			defaul = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			value, exist := BmGetInt64OrDefault(c, key, defaul)
			ctx.Convey("Then value,exist should not be nil.", func(ctx convey.C) {
				ctx.So(exist, convey.ShouldNotBeNil)
				ctx.So(value, convey.ShouldNotBeNil)
			})
		})
	})
}
