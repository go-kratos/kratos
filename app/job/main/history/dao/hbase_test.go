package dao

import (
	"context"
	"go-common/app/interface/main/history/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohashRowKey(t *testing.T) {
	convey.Convey("hashRowKey", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := hashRowKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocolumn(t *testing.T) {
	convey.Convey("column", t, func(ctx convey.C) {
		var (
			aid = int64(14771787)
			typ = int8(3)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.column(aid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAdd(t *testing.T) {
	convey.Convey("Add", t, func(ctx convey.C) {
		var (
			h = &model.History{Mid: 14771787, Aid: 14771787}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.Add(context.Background(), h)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
