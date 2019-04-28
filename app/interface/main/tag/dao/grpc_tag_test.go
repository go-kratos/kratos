package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddReport(t *testing.T) {
	var (
		c      = context.TODO()
		oid    = int64(1234)
		tid    = int64(1)
		mid    = int64(123)
		typ    = int32(3)
		partID = int32(64)
		reason = int32(1)
		score  = int32(100)
	)
	convey.Convey("AddReport", t, func(ctx convey.C) {
		err := d.AddReport(c, oid, tid, mid, typ, partID, reason, score)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTagMap(t *testing.T) {
	var (
		c    = context.TODO()
		tids = []int64{1}
		mid  = int64(0)
	)
	convey.Convey("TagMap", t, func(ctx convey.C) {
		res, err := d.TagMap(c, tids, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoResTag(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(123)
		tp  = int32(3)
	)
	convey.Convey("ResTag", t, func(ctx convey.C) {
		res, err := d.ResTag(c, oid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoResTagMap(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(123)
		tp  = int32(3)
	)
	convey.Convey("ResTagMap", t, func(ctx convey.C) {
		res, err := d.ResTagMap(c, oid, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoResTags(t *testing.T) {
	var (
		c    = context.TODO()
		oids = []int64{1, 2, 3, 4}
		tp   = int32(3)
	)
	convey.Convey("ResTags", t, func(ctx convey.C) {
		res, err := d.ResTags(c, oids, tp)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
