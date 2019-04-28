package app

import (
	"context"
	commonMdl "go-common/app/job/main/tv/model/common"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppkeyZone(t *testing.T) {
	var (
		category = int(0)
	)
	convey.Convey("keyZone", t, func(ctx convey.C) {
		p1 := keyZone(category)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAppFlush(t *testing.T) {
	var (
		c        = context.Background()
		category = int(0)
		idxRanks = []*commonMdl.IdxRank{}
	)
	convey.Convey("Flush", t, func(ctx convey.C) {
		err := d.Flush(c, category, idxRanks)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppTimeTrans(t *testing.T) {
	var (
		stimeStr = "2018-09-16 08:00:01"
	)
	convey.Convey("TimeTrans", t, func(ctx convey.C) {
		stime, err := TimeTrans(stimeStr)
		ctx.Convey("Then err should be nil.stime should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(stime, convey.ShouldNotBeNil)
		})
	})
}

func TestAppZAddIdx(t *testing.T) {
	var (
		c        = context.Background()
		category = int(0)
		ctimeStr = "2018-09-16 08:00:01"
		id       = int64(0)
	)
	convey.Convey("ZAddIdx", t, func(ctx convey.C) {
		err := d.ZAddIdx(c, category, ctimeStr, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppZRemIdx(t *testing.T) {
	var (
		c        = context.Background()
		category = int(0)
		id       = int64(0)
	)
	convey.Convey("ZRemIdx", t, func(ctx convey.C) {
		err := d.ZRemIdx(c, category, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
