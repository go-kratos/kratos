package data

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	aid = int64(10098208)
	mid = int64(10920044)
)

func TestDataOpenReply(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("OpenReply", t, func(ctx convey.C) {
		err := d.OpenReply(c, aid, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataCloseReply(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("CloseReply", t, func(ctx convey.C) {
		err := d.CloseReply(c, aid, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataCheckReply(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("CheckReply", t, func(ctx convey.C) {
		replyState, err := d.CheckReply(c, aid)
		t.Logf("%d", replyState)
		ctx.Convey("Then err should be nil.replyState should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(replyState, convey.ShouldNotBeNil)
		})
	})
}
