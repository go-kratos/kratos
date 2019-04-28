package dao

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
)

func TestDaoMainVip(t *testing.T) {
	convey.Convey("MainVip", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515308)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mv, err := d.MainVip(c, mid)
			ctx.Convey("Then err should be nil.mv should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mv, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGiveMVipGift(t *testing.T) {
	convey.Convey("GiveMVipGift", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(27515308)
			batchId = int(21)
			orderNo = "1" + strconv.Itoa(rand.Int()/100000)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.GiveMVipGift(c, mid, batchId, orderNo)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
