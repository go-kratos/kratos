package cms

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsLoadSnsAuthMap(t *testing.T) {
	var (
		ctx  = context.Background()
		sids = []int64{}
	)
	convey.Convey("LoadSnsAuthMap", t, func(c convey.C) {
		resMetas, err := d.LoadSnsAuthMap(ctx, sids)
		c.Convey("Then err should be nil.resMetas should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resMetas, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadEpsAuthMap(t *testing.T) {
	var (
		ctx   = context.Background()
		epids = []int64{}
	)
	convey.Convey("LoadEpsAuthMap", t, func(c convey.C) {
		resMetas, err := d.LoadEpsAuthMap(ctx, epids)
		c.Convey("Then err should be nil.resMetas should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resMetas, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadSnsCMSMap(t *testing.T) {
	var (
		ctx  = context.Background()
		sids = []int64{}
	)
	convey.Convey("LoadSnsCMSMap", t, func(c convey.C) {
		resMetas, err := d.LoadSnsCMSMap(ctx, sids)
		c.Convey("Then err should be nil.resMetas should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resMetas, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadSnsCMS(t *testing.T) {
	var (
		ctx  = context.Background()
		sids = []int64{}
		err  error
	)
	convey.Convey("LoadSnsCMS", t, func(c convey.C) {
		c.Convey("Then err should be nil.seasons,newestEpids should not be nil.", func(cx convey.C) {
			if sids, err = pickIDs(d.db, _pickSids); err != nil || len(sids) == 0 {
				fmt.Println("Empty Sids ", err)
				return
			}
			seasons, newestEpids, err := d.LoadSnsCMS(ctx, sids)
			cx.So(err, convey.ShouldBeNil)
			cx.So(newestEpids, convey.ShouldNotBeNil)
			cx.So(seasons, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadSnCMS(t *testing.T) {
	var (
		ctx = context.Background()
		sid = int64(0)
	)
	convey.Convey("LoadSnCMS", t, func(c convey.C) {
		sn, err := d.LoadSnCMS(ctx, sid)
		c.Convey("Then err should be nil.sn should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(sn, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadEpCMS(t *testing.T) {
	var (
		ctx  = context.Background()
		epid = int64(0)
	)
	convey.Convey("LoadEpCMS", t, func(c convey.C) {
		ep, err := d.LoadEpCMS(ctx, epid)
		c.Convey("Then err should be nil.ep should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(ep, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsLoadEpsCMS(t *testing.T) {
	var (
		ctx   = context.Background()
		epids = []int64{}
	)
	convey.Convey("LoadEpsCMS", t, func(c convey.C) {
		resMetas, err := d.LoadEpsCMS(ctx, epids)
		c.Convey("Then err should be nil.resMetas should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(resMetas, convey.ShouldNotBeNil)
		})
	})
}
