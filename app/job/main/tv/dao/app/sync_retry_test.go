package app

import (
	"context"
	"go-common/app/job/main/tv/model/common"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppSetRetry(t *testing.T) {
	var (
		c     = context.Background()
		retry = &common.SyncRetry{}
	)
	convey.Convey("SetRetry", t, func(ctx convey.C) {
		err := d.SetRetry(c, retry)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppGetRetry(t *testing.T) {
	var (
		c   = context.Background()
		req = &common.SyncRetry{}
	)
	convey.Convey("GetRetry", t, func(ctx convey.C) {
		times, err := d.GetRetry(c, req)
		ctx.Convey("Then err should be nil.times should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(times, convey.ShouldNotBeNil)
		})
	})
}

func TestAppDelRetry(t *testing.T) {
	var (
		c   = context.Background()
		req = &common.SyncRetry{}
	)
	convey.Convey("DelRetry", t, func(ctx convey.C) {
		err := d.DelRetry(c, req)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
