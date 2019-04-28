package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivestatTbl(t *testing.T) {
	var (
		aid = int64(1)
	)
	convey.Convey("statTbl", t, func(ctx convey.C) {
		p1 := statTbl(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivestat3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
	)
	convey.Convey("stat3", t, func(ctx convey.C) {
		_, err := d.stat3(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivestats3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{10097272}
	)
	convey.Convey("stats3", t, func(ctx convey.C) {
		_, err := d.stats3(c, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
