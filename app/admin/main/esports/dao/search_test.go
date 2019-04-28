package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/esports/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSearchArc(t *testing.T) {
	var (
		c = context.Background()
		p = &model.ArcListParam{Pn: 1, Ps: 1, Title: "dota"}
	)
	convey.Convey("SearchArc", t, func(ctx convey.C) {
		rs, total, err := d.SearchArc(c, p)
		ctx.Convey("Then err should be nil.rs,total should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(total, convey.ShouldNotBeNil)
			ctx.So(rs, convey.ShouldNotBeNil)
		})
	})
}
