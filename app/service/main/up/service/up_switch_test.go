package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceSetSwitch(t *testing.T) {
	convey.Convey("SetSwitch", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			state = int(0)
			from  = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			row, err := s.SetSwitch(c, mid, state, from)
			ctx.Convey("Then err should be nil.row should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(row, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpSwitch(t *testing.T) {
	convey.Convey("UpSwitch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			from = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			state, err := s.UpSwitchs(c, mid, from)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceRawUpSwitch(t *testing.T) {
	convey.Convey("RawUpSwitch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			from = uint8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			state, err := s.RawUpSwitch(c, mid, from)
			ctx.Convey("Then err should be nil.state should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}
