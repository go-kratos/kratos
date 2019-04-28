package academy

import (
	"context"
	"testing"

	"go-common/app/interface/main/creative/model/academy"

	"github.com/smartystreets/goconvey/convey"
)

func TestAcademyArchivesWithES(t *testing.T) {
	convey.Convey("ArchivesWithES", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aca = &academy.EsParam{
				OID:      int64(1),
				Tid:      []int64{1, 2},
				Business: 1,
				Pn:       10,
				Ps:       20,
				Keyword:  "string",
				Order:    "",
				IP:       "",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ArchivesWithES(c, aca)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAcademyKeywords(t *testing.T) {
	convey.Convey("Keywords", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Keywords(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
