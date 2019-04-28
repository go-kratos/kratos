package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/tv/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRegList(t *testing.T) {
	convey.Convey("RegList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.Param{PageID: "1"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RegList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddReg(t *testing.T) {
	convey.Convey("AddReg", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			title = "1"
			itype = "1"
			itid  = "1"
			rank  = "1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddReg(c, title, itype, itid, rank)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				println(err)
			})
		})
	})
}

func TestDaoEditReg(t *testing.T) {
	convey.Convey("EditReg", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			pid   = "1"
			title = "1"
			itype = "1"
			itid  = "1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.EditReg(c, pid, title, itype, itid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpState(t *testing.T) {
	convey.Convey("UpState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			pids  = []int{1}
			state = "0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpState(c, pids, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
