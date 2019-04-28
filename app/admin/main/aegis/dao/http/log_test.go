package http

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestHttpQueryLogSearch(t *testing.T) {
	convey.Convey("QueryLogSearch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			args = &model.ParamsQueryLog{
				Business:  231,
				Int0From:  "0",
				Int1:      []int64{0, 1},
				Int2:      []int64{0},
				Str0:      []string{"0"},
				CtimeFrom: "2019-01-01 00:00:00",
			}
			escm = model.EsCommon{Ps: 10, Pn: 1, Order: "ctime", Sort: "desc"}
		)

		ctx.Convey("success", func(ctx convey.C) {
			res, err := d.QueryLogSearch(c, args, escm)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
