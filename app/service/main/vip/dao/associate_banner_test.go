package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoEffectiveAssociateVips(t *testing.T) {
	convey.Convey("EffectiveAssociateVips", t, func(ctx convey.C) {
		res, err := d.EffectiveAssociateVips(context.Background())
		fmt.Println("res", res)
		ctx.Convey(" EffectiveAssociateVips ", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
