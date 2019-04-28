package up

import (
	"context"
	"go-common/app/service/main/up/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpSetSwitch(t *testing.T) {
	var (
		c = context.Background()
		u = &model.UpSwitch{
			MID:       int64(1),
			Attribute: 0,
		}
	)
	convey.Convey("SetSwitch", t, func(ctx convey.C) {
		id, err := d.SetSwitch(c, u)
		println(111, id)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestUpRawUpSwitch(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("RawUpSwitch", t, func(ctx convey.C) {
		_, err := d.RawUpSwitch(c, mid)
		ctx.Convey("Then err should be nil.u should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			//ctx.So(u, convey.ShouldNotBeNil)
		})
	})
}
