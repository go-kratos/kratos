package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
)

var a = &archive.Archive{
	ID:        10098208,
	Mid:       10920044,
	Attribute: 1097728,
	State:     0,
	Round:     99,
	TypeID:    21,
}

func TestServicearcReply(t *testing.T) {
	var (
		c           = context.Background()
		replySwitch = int64(1)
	)
	convey.Convey("arcReply", t, func(ctx convey.C) {
		err := s.arcReply(c, a, replySwitch)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceopenReply(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("openReply", t, func(ctx convey.C) {
		err := s.openReply(c, a, archive.ReplyOff)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicecloseReply(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("closeReply", t, func(ctx convey.C) {
		err := s.closeReply(c, a, archive.ReplyOn)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
