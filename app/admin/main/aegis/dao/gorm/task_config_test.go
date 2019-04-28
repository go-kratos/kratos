package gorm

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/task"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoQueryConfigs(t *testing.T) {
	convey.Convey("QueryConfigs", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			params = &task.QueryParams{
				BusinessID: 1,
				FlowID:     1,
				Btime:      "0000-00-00 00:00:00",
				Etime:      "2018-12-12 12:12:12",
				ConfName:   "mid",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.QueryConfigs(c, params)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDeleteConfig(t *testing.T) {
	convey.Convey("DeleteConfig", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteConfig(c, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetStateConfig(t *testing.T) {
	convey.Convey("SetStateConfig", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetStateConfig(c, 0, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateConfig(t *testing.T) {
	convey.Convey("UpdateConfig", t, func(ctx convey.C) {
		var (
			config = &task.Config{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateConfig(cntx, config)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
