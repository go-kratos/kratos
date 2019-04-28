package reply

import (
	"context"
	model "go-common/app/job/main/reply/model/reply"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyUpdate(t *testing.T) {
	convey.Convey("Update", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rpt = &model.Report{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Report.Update(c, rpt)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGet(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rpt, err := d.Report.Get(c, oid, rpID)
			ctx.Convey("Then err should be nil.rpt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetMapByOid(t *testing.T) {
	convey.Convey("GetMapByOid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			typ = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.Report.GetMapByOid(c, oid, typ)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyGetUsers(t *testing.T) {
	convey.Convey("GetUsers", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.Report.GetUsers(c, oid, tp, rpID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplySetUserReported(t *testing.T) {
	convey.Convey("SetUserReported", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			rpID = int64(0)
			now  = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Report.SetUserReported(c, oid, tp, rpID, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
