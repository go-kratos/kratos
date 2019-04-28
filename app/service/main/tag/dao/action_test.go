package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoActions(t *testing.T) {
	convey.Convey("Actions", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
			typ = int32(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rs, err := d.Actions(c, mid, typ)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAction(t *testing.T) {
	convey.Convey("AddAction", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.ResourceAction{
				Oid:    28843596,
				Mid:    35152246,
				Type:   3,
				Tid:    16906,
				Action: 0,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.AddAction(c, a)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
