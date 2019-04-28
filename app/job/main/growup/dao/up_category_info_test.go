package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoListUpNickname(t *testing.T) {
	convey.Convey("ListUpNickname", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{100}
			m    = make(map[int64]string)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ListUpNickname(c, mids, m)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetNickname(t *testing.T) {
	convey.Convey("GetNickname", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nickname, err := d.GetNickname(c, mid)
			ctx.Convey("Then err should be nil.nickname should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}
