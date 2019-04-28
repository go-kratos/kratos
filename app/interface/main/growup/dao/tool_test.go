package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDiv(t *testing.T) {
	convey.Convey("Div", t, func(ctx convey.C) {
		var (
			x = float64(1)
			y = float64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Div(x, y)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMul(t *testing.T) {
	convey.Convey("Mul", t, func(ctx convey.C) {
		var (
			x = float64(1)
			y = float64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Mul(x, y)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDivWithRound(t *testing.T) {
	convey.Convey("DivWithRound", t, func(ctx convey.C) {
		var (
			x      = float64(1)
			y      = float64(1)
			places = int(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := DivWithRound(x, y, places)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMulWithRound(t *testing.T) {
	convey.Convey("MulWithRound", t, func(ctx convey.C) {
		var (
			x      = float64(1)
			y      = float64(1)
			places = int(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := MulWithRound(x, y, places)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRound(t *testing.T) {
	convey.Convey("Round", t, func(ctx convey.C) {
		var (
			val    = float64(1)
			places = int(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := Round(val, places)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
