package app

import (
	"context"
	model "go-common/app/job/main/tv/model/pgc"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppEpCacheKey(t *testing.T) {
	var (
		epid = int(0)
	)
	convey.Convey("EpCacheKey", t, func(ctx convey.C) {
		p1 := EpCacheKey(epid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAppSeasonCacheKey(t *testing.T) {
	var (
		sid = int(0)
	)
	convey.Convey("SeasonCacheKey", t, func(ctx convey.C) {
		p1 := SeasonCacheKey(sid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAppSetEP(t *testing.T) {
	var (
		ctx = context.Background()
		res = &model.SimpleEP{}
	)
	convey.Convey("SetEP", t, func(cx convey.C) {
		err := d.SetEP(ctx, res)
		cx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppSetSeason(t *testing.T) {
	var (
		ctx = context.Background()
		res = &model.SimpleSeason{}
	)
	convey.Convey("SetSeason", t, func(cx convey.C) {
		err := d.SetSeason(ctx, res)
		cx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppCountEP(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("CountEP", t, func(cx convey.C) {
		count, err := d.CountEP(ctx)
		cx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAppCountSeason(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("CountSeason", t, func(cx convey.C) {
		count, err := d.CountSeason(ctx)
		cx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestAppRefreshEPMC(t *testing.T) {
	var (
		ctx    = context.Background()
		LastID = int(0)
		nbData = int(0)
	)
	convey.Convey("RefreshEPMC", t, func(cx convey.C) {
		myLast, err := d.RefreshEPMC(ctx, LastID, nbData)
		cx.Convey("Then err should be nil.myLast should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(myLast, convey.ShouldNotBeNil)
		})
	})
}

func TestAppRefreshSnMC(t *testing.T) {
	var (
		ctx    = context.Background()
		LastID = int(0)
		nbData = int(0)
	)
	convey.Convey("RefreshSnMC", t, func(cx convey.C) {
		myLast, err := d.RefreshSnMC(ctx, LastID, nbData)
		cx.Convey("Then err should be nil.myLast should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(myLast, convey.ShouldNotBeNil)
		})
	})
}

func TestAppPickSeason(t *testing.T) {
	var ctx = context.Background()
	convey.Convey("PickSeason", t, func(cx convey.C) {
		sid, errPick := pickEpSid()
		if errPick != nil {
			return
		}
		media, err := d.PickSeason(ctx, int(sid))
		cx.Convey("Then err should be nil.media should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(media, convey.ShouldNotBeNil)
		})
	})
}
