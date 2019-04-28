package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoreverseID(t *testing.T) {
	var (
		id = "123"
		l  = int(3)
	)
	convey.Convey("reverseID", t, func(ctx convey.C) {
		p1 := reverseID(id, l)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaocheckIDLen(t *testing.T) {
	var (
		id = "123"
	)
	convey.Convey("checkIDLen", t, func(ctx convey.C) {
		p1 := checkIDLen(id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaodiffTs(t *testing.T) {
	var (
		ts = int64(0)
	)
	convey.Convey("diffTs", t, func(ctx convey.C) {
		p1 := diffTs(ts)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaodiffID(t *testing.T) {
	var (
		id = int64(0)
	)
	convey.Convey("diffID", t, func(ctx convey.C) {
		p1 := diffID(id)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
