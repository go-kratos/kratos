package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoReportDmGarbageLog(t *testing.T) {
	convey.Convey("ReportDmGarbageLog", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			dm = &model.DM{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.ReportDmGarbageLog(c, dm)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReportDmLog(t *testing.T) {
	convey.Convey("ReportDmLog", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			dm = &model.DM{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.ReportDmLog(c, dm)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoreportUserLog(t *testing.T) {
	convey.Convey("reportUserLog", t, func(ctx convey.C) {
		var (
			c             = context.Background()
			dm            = &model.DM{}
			userLogType   = int(0)
			userLogAction = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.reportUserLog(c, dm, userLogType, userLogAction)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
