package dao

import (
	"context"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTable(t *testing.T) {
	var (
		oid = int64(0)
	)
	convey.Convey("table", t, func(ctx convey.C) {
		p1 := table(oid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoShare(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(0)
		tp  = int(0)
	)
	convey.Convey("Share", t, func(ctx convey.C) {
		share, err := d.Share(c, oid, tp)
		ctx.Convey("Then err should be nil.share should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(share, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddShare(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(rand.Intn(10000000))
		tp  = int(3)
	)
	convey.Convey("AddShare", t, func(ctx convey.C) {
		err := d.AddShare(c, oid, tp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
