package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_QQUnionID(t *testing.T) {
	var (
		c      = context.Background()
		openID = "2A9FE674CE0810761DC3F420239A8CD7"
	)
	convey.Convey("QQUnionID", t, func(ctx convey.C) {
		res, err := d.QQUnionID(c, openID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		fmt.Printf("(%+v) error(%+v)", res, err)
	})
}
