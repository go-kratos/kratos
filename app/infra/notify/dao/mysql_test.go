package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoLoadNotify(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("LoadNotify", t, func(ctx convey.C) {
		_, err := d.LoadNotify(c, "sh001")
		ctx.Convey("Then err should be nil.ns should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			//ctx.So(ns, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoLoadPub(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("LoadPub", t, func(ctx convey.C) {
		_, err := d.LoadPub(c)
		ctx.Convey("Then err should be nil.ps should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			//	ctx.So(ps, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoFilters(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(7)
	)
	convey.Convey("Filters", t, func(ctx convey.C) {
		fs, err := d.Filters(c, id)
		ctx.Convey("Then err should be nil.fs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fs, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddFailBk(t *testing.T) {
	var (
		c       = context.TODO()
		topic   = ""
		group   = ""
		cluster = ""
		msg     = ""
		index   = int64(0)
	)
	convey.Convey("AddFailBk", t, func(ctx convey.C) {
		id, err := d.AddFailBk(c, topic, group, cluster, msg, index)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelFailBk(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("DelFailBk", t, func(ctx convey.C) {
		affected, err := d.DelFailBk(c, id)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoLoadFailBk(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("LoadFailBk", t, func(ctx convey.C) {
		fbs, err := d.LoadFailBk(c)
		ctx.Convey("Then err should be nil.fbs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(fbs, convey.ShouldNotBeNil)
		})
	})
}
