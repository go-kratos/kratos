package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveAddNetSafeMd5(t *testing.T) {
	var (
		c   = context.TODO()
		nid = int64(0)
		md5 = ""
	)
	convey.Convey("AddNetSafeMd5", t, func(ctx convey.C) {
		err := d.AddNetSafeMd5(c, nid, md5)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveNotifyNetSafe(t *testing.T) {
	var (
		c   = context.TODO()
		nid = int64(0)
	)
	convey.Convey("NotifyNetSafe", t, func(ctx convey.C) {
		err := d.NotifyNetSafe(c, nid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
