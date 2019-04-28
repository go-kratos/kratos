package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPingMySQL(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("PingMySQL", t, func(ctx convey.C) {
		err := PingMySQL(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoClose(t *testing.T) {
	convey.Convey("Close", t, func(ctx convey.C) {
		Close()
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
