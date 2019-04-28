package cms

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsArcCMSCacheKey(t *testing.T) {
	var (
		aid = int64(0)
	)
	convey.Convey("ArcCMSCacheKey", t, func(c convey.C) {
		p1 := d.ArcCMSCacheKey(aid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsVideoCMSCacheKey(t *testing.T) {
	var (
		cid = int64(0)
	)
	convey.Convey("VideoCMSCacheKey", t, func(c convey.C) {
		p1 := d.VideoCMSCacheKey(cid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsArcsMetaCache(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("ArcsMetaCache", t, func(c convey.C) {
		sids, errPick := pickIDs(d.db, _pickAids)
		if errPick != nil || len(sids) == 0 {
			fmt.Println("Empty sids ", errPick)
			return
		}
		cached, missed, err := d.ArcsMetaCache(ctx, sids)
		c.Convey("Then err should be nil.cached,missed should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(len(missed)+len(cached), convey.ShouldEqual, len(sids))
		})
	})
}

func TestCmsVideosMetaCache(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("VideosMetaCache", t, func(c convey.C) {
		sids, errPick := pickIDs(d.db, _pickCids)
		if errPick != nil || len(sids) == 0 {
			fmt.Println("Empty sids ", errPick)
			return
		}
		cached, missed, err := d.VideosMetaCache(ctx, sids)
		c.Convey("Then err should be nil.cached,missed should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(len(missed)+len(cached), convey.ShouldEqual, len(sids))
		})
	})
}
