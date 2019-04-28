package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/tag/model"
	"go-common/app/service/main/archive/api"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoarcKey(t *testing.T) {
	convey.Convey("arcKey", t, func(ctx convey.C) {
		var (
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := arcKey(tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoregionArcKey(t *testing.T) {
	convey.Convey("regionArcKey", t, func(ctx convey.C) {
		var (
			rid = int16(17)
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := regionArcKey(rid, tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoregionOriArcKey(t *testing.T) {
	convey.Convey("regionOriArcKey", t, func(ctx convey.C) {
		var (
			rid = int16(17)
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := regionOriArcKey(rid, tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddTagNewArcCache(t *testing.T) {
	convey.Convey("AddTagNewArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			arc  = &model.Archive{Aid: 29661790, PubTime: "2018-10-23 15:59:45"}
			tids = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddTagNewArcCache(c, arc, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemTidArcCache(t *testing.T) {
	convey.Convey("RemTidArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aid  = int64(29661790)
			tids = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemTidArcCache(c, aid, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoexpireTagNewArc(t *testing.T) {
	convey.Convey("expireTagNewArc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.expireTagNewArc(c, tid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddRegionTagNewArcCache(t *testing.T) {
	convey.Convey("AddRegionTagNewArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			arc  = &model.Archive{Aid: 29661790, PubTime: "2018-10-23 15:59:45"}
			tids = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddRegionTagNewArcCache(c, arc, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemoveRegionNewArcCache(t *testing.T) {
	convey.Convey("RemoveRegionNewArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			arc  = &model.Archive{Aid: 29661790, PubTime: "2018-10-23 15:59:45"}
			tids = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemoveRegionNewArcCache(c, arc, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoexpireRegionNewArc(t *testing.T) {
	convey.Convey("expireRegionNewArc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int16(17)
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.expireRegionNewArc(c, rid, tid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoexpireRegionOriArc(t *testing.T) {
	convey.Convey("expireRegionOriArc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int16(17)
			tid = int64(9222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.expireRegionOriArc(c, rid, tid)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetPubTime(t *testing.T) {
	convey.Convey("getPubTime", t, func(ctx convey.C) {
		var (
			aid     = int64(29661790)
			pubDate = "2018-10-23 15:59:45"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			timeStamp, err := d.getPubTime(aid, pubDate)
			ctx.Convey("Then err should be nil.timeStamp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(timeStamp, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateTagNewArcCache(t *testing.T) {
	convey.Convey("UpdateTagNewArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tid  = int64(9222)
			arcs = make(map[int64]*api.Arc)
		)
		arcs[29661790] = &api.Arc{Aid: 29661790}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateTagNewArcCache(c, tid, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
