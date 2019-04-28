package archive

import (
	"context"
	"testing"

	arccli "go-common/app/service/main/archive/api"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivekeyArc(t *testing.T) {
	var (
		aid = int64(0)
	)
	convey.Convey("keyArc", t, func(ctx convey.C) {
		p1 := keyArc(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivekeyView(t *testing.T) {
	var (
		aid = int64(0)
	)
	convey.Convey("keyView", t, func(ctx convey.C) {
		p1 := keyView(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpArcCache(t *testing.T) {
	var (
		c = context.Background()
		a = &arccli.Arc{}
	)
	convey.Convey("UpArcCache", t, func(ctx convey.C) {
		err := d.UpArcCache(c, a)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveUpViewCache(t *testing.T) {
	var (
		c = context.Background()
		v = &arccli.ViewReply{
			Arc: &arccli.Arc{
				Aid: 123,
			},
		}
	)
	convey.Convey("UpViewCache", t, func(ctx convey.C) {
		err := d.UpViewCache(c, v)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
