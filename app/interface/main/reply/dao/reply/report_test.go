package reply

import (
	"context"
	"go-common/app/interface/main/reply/conf"
	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewReportDao(t *testing.T) {
	convey.Convey("NewReportDao", t, func(ctx convey.C) {
		var (
			db = sql.NewMySQL(conf.Conf.MySQL.Reply)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewReportDao(db)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyReporthit(t *testing.T) {
	convey.Convey("hit", t, func(ctx convey.C) {
		var (
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.Report.hit(oid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyReportInsert(t *testing.T) {
	convey.Convey("Insert", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rpt = &model.Report{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.Report.Insert(c, rpt)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyReportGet(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpt, err := d.Report.Get(c, oid, rpID)
			ctx.Convey("Then err should be nil.rpt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyInsertUser(t *testing.T) {
	convey.Convey("InsertUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ru = &model.ReportUser{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.Report.InsertUser(c, ru)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
