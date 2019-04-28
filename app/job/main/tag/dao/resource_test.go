package dao

import (
	"context"
	"go-common/app/job/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateResTags(t *testing.T) {
	var (
		c  = context.Background()
		rt = &model.ResTag{
			Oid:  1,
			Tids: []int64{1, 2, 3},
		}
	)
	convey.Convey("UpdateResTags", t, func(ctx convey.C) {
		affect, err := d.UpdateResTags(c, rt)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertResTags(t *testing.T) {
	var (
		c  = context.Background()
		rt = &model.ResTag{
			Oid:  1,
			Tids: []int64{1, 2, 3},
		}
	)
	convey.Convey("InsertResTags", t, func(ctx convey.C) {
		affect, err := d.InsertResTags(c, rt)
		ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affect, convey.ShouldNotBeNil)
		})
	})
}
