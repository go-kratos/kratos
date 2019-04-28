package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCallback(t *testing.T) {
	var (
		c          = context.Background()
		chall      = &model.Challenge{}
		businessID = int8(6)
	)
	convey.Convey("Callback", t, func(ctx convey.C) {
		err := d.Callback(c, chall, businessID)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
