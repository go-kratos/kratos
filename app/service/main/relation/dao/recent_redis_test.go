package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddRctFollower(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("AddRctFollower", t, func(cv convey.C) {
		err := d.AddRctFollower(c, mid, fid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelRctFollower(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		fid = int64(0)
	)
	convey.Convey("DelRctFollower", t, func(cv convey.C) {
		err := d.DelRctFollower(c, mid, fid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRctFollowerCount(t *testing.T) {
	var (
		ctx = context.Background()
		fid = int64(0)
	)
	convey.Convey("RctFollowerCount", t, func(cv convey.C) {
		p1, err := d.RctFollowerCount(ctx, fid)
		cv.Convey("Then err should be nil.p1 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoEmptyRctFollower(t *testing.T) {
	var (
		ctx = context.Background()
		fid = int64(0)
	)
	convey.Convey("EmptyRctFollower", t, func(cv convey.C) {
		err := d.EmptyRctFollower(ctx, fid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRctFollowerNotify(t *testing.T) {
	var (
		c   = context.Background()
		fid = int64(0)
	)
	convey.Convey("RctFollowerNotify", t, func(cv convey.C) {
		p1, err := d.RctFollowerNotify(c, fid)
		cv.Convey("Then err should be nil.p1 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetRctFollowerNotify(t *testing.T) {
	var (
		c    = context.Background()
		fid  = int64(0)
		flag bool
	)
	convey.Convey("SetRctFollowerNotify", t, func(cv convey.C) {
		err := d.SetRctFollowerNotify(c, fid, flag)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}
