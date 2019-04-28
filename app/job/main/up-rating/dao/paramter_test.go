package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetAllParamter(t *testing.T) {
	convey.Convey("GetAllParamter", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO rating_parameter(name,value) VALUES('test', 123)")
			paramters, err := d.GetAllParamter(c)
			ctx.Convey("Then err should be nil.paramters should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(paramters, convey.ShouldNotBeNil)
			})
		})
	})
}
