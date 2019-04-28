package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSearchChall(t *testing.T) {
	convey.Convey("SearchChall", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			params = &model.SearchParams{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SearchChall(c, params)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSearchAppeal(t *testing.T) {
	convey.Convey("SearchAppeal", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			cond = model.AppealSearchCond{
				Fields:  []string{"id", "mid"},
				Bid:     []int{2, 28},
				TTimeTo: time.Now().AddDate(0, 0, -3).Format("2006-01-02 15:04:05"),
				PS:      1000,
				PN:      1,
				Order:   "id",
				Sort:    "desc",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SearchAppeal(c, cond)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
