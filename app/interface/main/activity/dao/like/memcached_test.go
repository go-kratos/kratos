package like

import (
	"context"
	"go-common/app/interface/main/activity/model/like"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikekeyInfo(t *testing.T) {
	convey.Convey("keyInfo", t, func(ctx convey.C) {
		var (
			sid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyInfo(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSetInfoCache(t *testing.T) {
	convey.Convey("SetInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			v   = &like.Subject{}
			sid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetInfoCache(c, v, sid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeInfoCache(t *testing.T) {
	convey.Convey("InfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			v, err := d.InfoCache(c, sid)
			ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", v)
			})
		})
	})
}
