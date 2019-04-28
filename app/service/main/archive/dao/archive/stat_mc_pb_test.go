package archive

import (
	"context"
	"go-common/app/service/main/archive/api"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivestatPBKey(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("statPBKey", t, func(ctx convey.C) {
		p1 := statPBKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveclickPBKey(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("clickPBKey", t, func(ctx convey.C) {
		p1 := clickPBKey(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivestatCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("statCache3", t, func(ctx convey.C) {
		_, err := d.statCache3(c, aid)
		ctx.Convey("Then err should be nil.st should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivestatCaches3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{1, 2}
	)
	convey.Convey("statCaches3", t, func(ctx convey.C) {
		cached, missed, err := d.statCaches3(c, aids)
		ctx.Convey("Then err should be nil.cached,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missed, convey.ShouldNotBeNil)
			ctx.So(cached, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveaddStatCache3(t *testing.T) {
	var (
		c  = context.TODO()
		st = &api.Stat{}
	)
	convey.Convey("addStatCache3", t, func(ctx convey.C) {
		err := d.addStatCache3(c, st)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveclickCache3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("clickCache3", t, func(ctx convey.C) {
		_, err := d.clickCache3(c, aid)
		ctx.Convey("Then err should be nil.clk should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveaddClickCache3(t *testing.T) {
	var (
		c   = context.TODO()
		clk = &api.Click{}
	)
	convey.Convey("addClickCache3", t, func(ctx convey.C) {
		err := d.addClickCache3(c, clk)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
