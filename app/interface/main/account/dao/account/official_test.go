package account

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestKeyMonthlyOfficialSubmittedTimes(t *testing.T) {
	convey.Convey("keyMonthlyOfficialSubmittedTimes", t, func(ctx convey.C) {
		now := time.Now()
		key := keyMonthlyOfficialSubmittedTimes(now, 1)
		ctx.So(key, convey.ShouldEqual, fmt.Sprintf("ot_%d_%d", now.Month(), 1))
	})
}

func TestIncreaseMonthlyOfficialSubmittedTimes(t *testing.T) {
	convey.Convey("IncreaseMonthlyOfficialSubmittedTimes", t, func(ctx convey.C) {
		p1, err := d.IncreaseMonthlyOfficialSubmittedTimes(context.Background(), 1)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestGetMonthlyOfficialSubmittedTimes(t *testing.T) {
	convey.Convey("GetMonthlyOfficialSubmittedTimes", t, func(ctx convey.C) {
		p1, err := d.GetMonthlyOfficialSubmittedTimes(context.Background(), 1)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
