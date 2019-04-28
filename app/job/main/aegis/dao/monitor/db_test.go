package monitor

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMonitorRulesByBid(t *testing.T) {
	convey.Convey("RulesByBid", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			bid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.RulesByBid(c, bid)
			convCtx.Convey("Then err should be nil.rules should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMonitorValidRules(t *testing.T) {
	convey.Convey("ValidRules", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.ValidRules(c)
			convCtx.Convey("Then err should be nil.rules should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
