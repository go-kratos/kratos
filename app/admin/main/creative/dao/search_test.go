package dao

import (
	"context"
	"go-common/app/admin/main/creative/model/academy"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoArchivesWithES(t *testing.T) {
	convey.Convey("ArchivesWithES", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aca = &academy.EsParam{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ArchivesWithES(c, aca)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
