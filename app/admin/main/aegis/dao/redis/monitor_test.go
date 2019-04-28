package redis

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRedisMoniRuleStats(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("MoniRuleStats", t, func(ctx convey.C) {
		_, err := d.MoniRuleStats(c, 1, 0, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
func TestRedisMoniRuleOids(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("MoniRuleOids", t, func(ctx convey.C) {
		_, err := d.MoniRuleOids(c, 1, 0, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
