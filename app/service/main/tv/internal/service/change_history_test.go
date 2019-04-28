package service

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceChangeHistory(t *testing.T) {
	convey.Convey("ChangeHistory", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			hid = int32(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ch, err := s.ChangeHistory(c, hid)
			ctx.Convey("Then err should be nil.ch should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ch, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceChangeHistorys(t *testing.T) {
	convey.Convey("ChangeHistorys", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(27515308)
			from = int32(0)
			to   = int32(time.Now().Unix())
			pn   = int32(1)
			ps   = int32(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			chs, total, err := s.ChangeHistorys(c, mid, from, to, pn, ps)
			ctx.Convey("Then err should be nil.chs,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(chs, convey.ShouldNotBeNil)
			})
		})
	})
}
