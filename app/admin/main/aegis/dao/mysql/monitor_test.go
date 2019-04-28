package mysql

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMysqlMoniBizRules(t *testing.T) {
	convey.Convey("MoniBizRules", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			bid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.MoniBizRules(c, bid)
			convCtx.Convey("Then err should be nil.rules should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMysqlMoniRule(t *testing.T) {
	convey.Convey("MoniRule", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.MoniRule(c, id)
			convCtx.Convey("Then err should be nil.rules should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
