package dao

import (
	"context"
	"go-common/app/job/main/member/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddExpLog(t *testing.T) {
	convey.Convey("AddExpLog", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			ul  = &model.UserLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.AddExpLog(ctx, ul)
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}

func TestDaoAddMoralLog(t *testing.T) {
	convey.Convey("AddMoralLog", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			ul  = &model.UserLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.AddMoralLog(ctx, ul)
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}

func TestDaoaddLog(t *testing.T) {
	convey.Convey("addLog", t, func(convCtx convey.C) {
		var (
			ctx      = context.Background()
			business = int(0)
			action   = ""
			ul       = &model.UserLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.addLog(ctx, business, action, ul)
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}
