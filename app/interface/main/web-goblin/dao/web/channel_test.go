package web

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestWebChCard(t *testing.T) {
	convey.Convey("ChCard", t, func(c convey.C) {
		var (
			ctx = context.Background()
			now = time.Now()
		)
		res, err := d.ChCard(ctx, now)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
		convey.Printf("%+v", res)
	})
}
