package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubBigData(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		msg = interface{}("PubBigData")
	)
	convey.Convey("PubBigData", t, func(ctx convey.C) {
		err := d.PubBigData(c, aid, msg)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPubCoinJob(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		msg = interface{}("PubCoinJob")
	)
	convey.Convey("PubCoinJob", t, func(ctx convey.C) {
		err := d.PubCoinJob(c, aid, msg)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPubStat(t *testing.T) {
	var (
		c     = context.TODO()
		aid   = int64(1)
		tp    = int64(2)
		count = int64(10)
	)
	convey.Convey("PubStat", t, func(ctx convey.C) {
		err := d.PubStat(c, aid, tp, count)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
