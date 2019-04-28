package cms

import (
	"context"
	"fmt"
	"go-common/app/interface/main/tv/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsSeaCacheKey(t *testing.T) {
	var (
		sid = int64(0)
	)
	convey.Convey("SeaCacheKey", t, func(c convey.C) {
		p1 := d.SeaCacheKey(sid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsEPCacheKey(t *testing.T) {
	var (
		epid = int64(0)
	)
	convey.Convey("EPCacheKey", t, func(c convey.C) {
		p1 := d.EPCacheKey(epid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsGetSeasonCache(t *testing.T) {
	var (
		ctx = context.Background()
		sid = int64(0)
	)
	convey.Convey("GetSeasonCache", t, func(c convey.C) {
		s, err := d.GetSeasonCache(ctx, sid)
		c.Convey("Then err should be nil.s should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(s, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsGetEPCache(t *testing.T) {
	var (
		ctx  = context.Background()
		epid = int64(0)
	)
	convey.Convey("GetEPCache", t, func(c convey.C) {
		ep, err := d.GetEPCache(ctx, epid)
		c.Convey("Then err should be nil.ep should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(ep, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsAddSnAuthCache(t *testing.T) {
	var (
		ctx = context.Background()
		s   = &model.SnAuth{}
	)
	convey.Convey("AddSnAuthCache", t, func(c convey.C) {
		err := d.AddSnAuthCache(ctx, s)
		c.Convey("Then err should be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCmsAddEpAuthCache(t *testing.T) {
	var (
		ctx = context.Background()
		ep  = &model.EpAuth{}
	)
	convey.Convey("AddEpAuthCache", t, func(c convey.C) {
		err := d.AddEpAuthCache(ctx, ep)
		c.Convey("Then err should be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCmssnAuthCache(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("snAuthCache", t, func(c convey.C) {
		sids, errPick := pickIDs(d.db, _pickSids)
		if errPick != nil || len(sids) == 0 {
			fmt.Println("Empty sids ", errPick)
			return
		}
		cached, missed, err := d.snAuthCache(ctx, sids)
		c.Convey("Then err should be nil.cached,missed should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(len(missed)+len(cached), convey.ShouldBeGreaterThan, 0)
		})
	})
}

func TestCmsSnAuth(t *testing.T) {
	var (
		ctx = context.Background()
		sid = int64(0)
	)
	convey.Convey("SnAuth", t, func(c convey.C) {
		sn, err := d.SnAuth(ctx, sid)
		c.Convey("Then err should be nil.sn should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(sn, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsEpAuth(t *testing.T) {
	var (
		ctx  = context.Background()
		epid = int64(0)
	)
	convey.Convey("EpAuth", t, func(c convey.C) {
		ep, err := d.EpAuth(ctx, epid)
		c.Convey("Then err should be nil.ep should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(ep, convey.ShouldNotBeNil)
		})
	})
}
