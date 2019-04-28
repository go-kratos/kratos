package recommend

import (
	"context"
	"testing"
	"time"

	cdm "go-common/app/interface/main/app-card/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestRecommend(t *testing.T) {
	var (
		c          = context.TODO()
		plat       = int8(1)
		buvid      = ""
		mid        = int64(1)
		build      = int(1111)
		loginEvent = int(0)
		parentMode = int(0)
		recsysMode = int(1)
		zoneID     = int64(0)
		group      = int(0)
		interest   = ""
		network    = ""
		style      = int(0)
		column     cdm.ColumnStatus
		flush      = int(0)
		autoPlay   = ""
		now        time.Time
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, _, _, _, err := d.Recommend(c, plat, buvid, mid, build, loginEvent, parentMode, recsysMode, zoneID, group, interest, network, style, column, flush, autoPlay, now)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestHots(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.Hots(c)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestTagTop(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		tid = int64(1)
		rn  = int(0)
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.TagTop(c, mid, tid, rn)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestGroup(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		d.Group(c)
	})
}
