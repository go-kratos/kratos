package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAvgUtimeByUID(t *testing.T) {
	var (
		c     = context.TODO()
		uid   = int64(421)
		stime = time.Now()
		etime = time.Now()
	)
	convey.Convey("AvgUtimeByUID", t, func(ctx convey.C) {
		utime, err := d.AvgUtimeByUID(c, uid, stime, etime)
		ctx.Convey("Then err should be nil.utime should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(utime, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSumDurationByUID(t *testing.T) {
	var (
		c     = context.TODO()
		uid   = int64(421)
		stime = time.Now()
		etime = time.Now()
	)
	convey.Convey("SumDurationByUID", t, func(ctx convey.C) {
		duration, err := d.SumDurationByUID(c, uid, stime, etime)
		ctx.Convey("Then err should be nil.duration should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(duration, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoActionCountByUID(t *testing.T) {
	var (
		c     = context.TODO()
		uid   = int64(421)
		stime = time.Now()
		etime = time.Now()
	)
	convey.Convey("ActionCountByUID", t, func(ctx convey.C) {
		mapAction, err := d.ActionCountByUID(c, uid, stime, etime)
		ctx.Convey("Then err should be nil.mapAction should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(mapAction, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPassCountByUID(t *testing.T) {
	var (
		c     = context.TODO()
		uid   = int64(421)
		stime = time.Now()
		etime = time.Now()
	)
	convey.Convey("PassCountByUID", t, func(ctx convey.C) {
		count, err := d.PassCountByUID(c, uid, stime, etime)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSubjectCountByUID(t *testing.T) {
	var (
		c     = context.TODO()
		uid   = int64(421)
		stime = time.Now()
		etime = time.Now()
	)
	convey.Convey("SubjectCountByUID", t, func(ctx convey.C) {
		count, err := d.SubjectCountByUID(c, uid, stime, etime)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoActiveUids(t *testing.T) {
	var (
		c     = context.TODO()
		stime = time.Now().AddDate(-1, 0, 0)
		etime = time.Now()
	)
	convey.Convey("ActiveUids", t, func(ctx convey.C) {
		_, err := d.ActiveUids(c, stime, etime)
		ctx.Convey("Then err should be nil.uids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
