package dao

import (
	"context"
	"go-common/app/admin/main/tv/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoArcES(t *testing.T) {
	var (
		c   = context.Background()
		req = &model.ReqArcES{
			AID:          "10110475",
			Valid:        "1",
			Result:       "1",
			Mids:         []int64{477132},
			Typeids:      []int32{24},
			MtimeOrder:   "1",
			PubtimeOrder: "1",
		}
	)
	convey.Convey("ArcES", t, func(ctx convey.C) {
		data, err := d.ArcES(c, req)
		ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}
