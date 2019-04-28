package ugc

import (
	"context"
	"testing"

	"database/sql"
	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestUgcDeletedUp(t *testing.T) {
	var c = context.Background()
	convey.Convey("DeletedUp", t, func(ctx convey.C) {
		mid, err := d.DeletedUp(c)
		if err == sql.ErrNoRows {
			fmt.Println("No to delete data")
			return
		}
		ctx.Convey("Then err should be nil.mid should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(mid, convey.ShouldNotBeNil)
		})
	})
}

func TestUgcFinishDelUp(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("FinishDelUp", t, func(ctx convey.C) {
		err := d.FinishDelUp(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUgcPpDelUp(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("PpDelUp", t, func(ctx convey.C) {
		err := d.PpDelUp(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUgcCountUpArcs(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("CountUpArcs", t, func(ctx convey.C) {
		count, err := d.CountUpArcs(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestUgcUpArcs(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("UpArcs", t, func(ctx convey.C) {
		aids, err := d.UpArcs(c, mid)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(aids, convey.ShouldNotBeNil)
		})
	})
}
