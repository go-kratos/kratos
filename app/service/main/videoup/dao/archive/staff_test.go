package archive

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/smartystreets/goconvey/convey"
)

func TestStaffMid(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("MidCount", t, func(ctx convey.C) {
		count, err := d.MidCount(c, 2880441)
		spew.Dump(count)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestApply(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Apply", t, func(ctx convey.C) {
		data, err := d.Apply(c, 1)
		if err == nil {
			spew.Dump(data)
		}
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestApplys(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Applys", t, func(ctx convey.C) {
		data, err := d.Applys(c, []int64{1, 11})
		if err == nil {
			spew.Dump(data)
		}
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFilterApplys(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("FilterApplys", t, func(ctx convey.C) {
		data, err := d.FilterApplys(c, []int64{23213, 4052032}, 4052032)
		if err == nil {
			spew.Dump(data)
		}
		if err != nil {
			spew.Dump(err)
		}
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
