package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoregionArcKey(t *testing.T) {
	var (
		rid = int64(0)
		tid = int64(0)
	)
	convey.Convey("regionArcKey", t, func(ctx convey.C) {
		p1 := regionArcKey(rid, tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoregionOriArcKey(t *testing.T) {
	var (
		rid = int64(0)
		tid = int64(0)
	)
	convey.Convey("regionOriArcKey", t, func(ctx convey.C) {
		p1 := regionOriArcKey(rid, tid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddTagArcs(t *testing.T) {
	var (
		c      = context.TODO()
		tid    = int64(0)
		arcMap = map[int64]*model.SearchRes{
			0: {
				ID: 0,
			},
			1: {
				ID: 1,
			},
		}
	)
	convey.Convey("AddTagArcs", t, func(ctx convey.C) {
		err := d.AddTagArcs(c, tid, arcMap)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
