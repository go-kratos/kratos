package service

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/email"
)

var (
	aid    = a.ID
	action = email.RetryActionReply
	flags  = archive.ReplyOn
	flagA  = archive.ReplyOff
)

func TestServiceaddRetry(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("addRetry", t, func(ctx convey.C) {
		err := s.addRetry(c, aid, action, flags, flagA)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceremoveRetry(t *testing.T) {
	var (
		c = context.Background()
	)
	TestServiceaddRetry(t)
	convey.Convey("removeRetry", t, func(ctx convey.C) {
		err := s.removeRetry(c, aid, action)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceretryProc(t *testing.T) {
	convey.Convey("retryProc", t, func(ctx convey.C) {
		go func() {
			time.Sleep(1 * time.Second)
			s.closed = true
		}()
		s.retryProc()
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
