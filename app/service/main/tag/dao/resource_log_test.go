package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoResourceLogByTid(t *testing.T) {
	convey.Convey("ResourceLogByTid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			tid = int64(1833)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceLogByTid(c, oid, tid, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResourceLogs(t *testing.T) {
	convey.Convey("ResourceLogs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
			ps  = int(20)
			pn  = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceLogs(c, oid, typ, ps, pn)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResourceLogsAdmin(t *testing.T) {
	convey.Convey("ResourceLogsAdmin", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int32(3)
			ps  = int(20)
			pn  = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ResourceLogsAdmin(c, oid, typ, ps, pn)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoResourceLog(t *testing.T) {
	convey.Convey("ResourceLog", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(28843596)
			logID = int64(2233)
			typ   = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.ResourceLog(c, oid, logID, typ)
		})
	})
}

func TestDaoAddResourceLog(t *testing.T) {
	convey.Convey("AddResourceLog", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tname = "搞笑"
			m     = &model.ResourceLog{
				Oid:  28843596,
				Type: 3,
				Tid:  1833,
				Mid:  35152246,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.AddResourceLog(c, tname, m)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
