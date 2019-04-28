package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaorowKeyContent(t *testing.T) {
	convey.Convey("rowKeyContent", t, func(ctx convey.C) {
		var (
			id  = int64(233)
			typ = "test_typ"
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			key := rowKeyContent(id, typ)
			ctx.Convey("Then key should equal id_typ.", func(ctx convey.C) {
				ctx.So(key, convey.ShouldEqual, "233_test_typ")
			})
		})
	})
}

func TestDaoContent(t *testing.T) {
	convey.Convey("Content", t, func(ctx convey.C) {
		var (
			c       = context.TODO()
			id      = int64(2333)
			typ     = "test_type"
			content = "test_content"
		)
		ctx.Convey("When SetContent", func(ctx convey.C) {
			err := d.SetContent(c, id, typ, content)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Convey("When get Content", func(ctx convey.C) {
					content2, err := d.Content(c, id, typ)
					ctx.Convey("Then err should be nil.content2 should equal content.", func(ctx convey.C) {
						ctx.So(err, convey.ShouldBeNil)
						ctx.So(content2, convey.ShouldEqual, content)
					})
				})
			})
		})
	})
}
