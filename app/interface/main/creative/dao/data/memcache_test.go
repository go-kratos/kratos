package data

import (
	"context"
	"go-common/app/interface/main/creative/model/data"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatakeyStat(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("keyStat", t, func(ctx convey.C) {
		p1 := keyStat(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "s_1")
		})
	})
}

func TestDatakeyUpStat(t *testing.T) {
	var (
		mid  = int64(1)
		date = ""
	)
	convey.Convey("keyUpStat", t, func(ctx convey.C) {
		p1 := keyUpStat(mid, date)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "sup_1")
		})
	})
}

func TestDatastatCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("statCache", t, func(ctx convey.C) {
		_, err := d.statCache(c, mid)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataaddStatCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		st  = &data.Stat{}
	)
	convey.Convey("addStatCache", t, func(ctx convey.C) {
		err := d.addStatCache(c, mid, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataupBaseStatCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
	)
	convey.Convey("upBaseStatCache", t, func(ctx convey.C) {
		_, err := d.upBaseStatCache(c, mid, dt)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataaddUpBaseStatCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
		st  = &data.UpBaseStat{}
	)
	convey.Convey("addUpBaseStatCache", t, func(ctx convey.C) {
		err := d.addUpBaseStatCache(c, mid, dt, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDataDelUpBaseStatCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		dt  = ""
	)
	convey.Convey("DelUpBaseStatCache", t, func(ctx convey.C) {
		err := d.DelUpBaseStatCache(c, mid, dt)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
