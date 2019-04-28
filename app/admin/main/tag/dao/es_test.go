package dao

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoESearchTag(t *testing.T) {
	var (
		c     = context.TODO()
		esTag = &model.ESTag{}
	)
	convey.Convey("ESearchTag", t, func(ctx convey.C) {
		_, err := d.ESearchTag(c, esTag)
		ctx.Convey("Then err should be nil.tags should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoESearchArchives(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{10099867, 10099860}
	)
	convey.Convey("ESearchArchives", t, func(ctx convey.C) {
		arcs, err := d.ESearchArchives(c, aids)
		ctx.Convey("Then err should be nil. arcs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(arcs, convey.ShouldNotBeNil)
		})
	})
}
