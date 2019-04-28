package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoRanking(t *testing.T) {
	convey.Convey("Ranking", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(0)
			rankType = int(0)
			day      = int(0)
			arcType  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.c.Host.Rank+_rankURL).Reply(200).JSON(`{"code":0,"num":8,"list":[{"aid":33986715,"score":704},{"aid":33913315,"score":571}]}`)
			res, err := d.Ranking(c, rid, rankType, day, arcType)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRankingIndex(t *testing.T) {
	convey.Convey("RankingIndex", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			day = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.c.Host.Rank+_rankURL).Reply(200).JSON(`{"code":0,"num":8,"list":[{"aid":33986715,"score":704},{"aid":33913315,"score":571}]}`)
			res, err := d.RankingIndex(c, day)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRankingRegion(t *testing.T) {
	convey.Convey("RankingRegion", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(0)
			day      = int(0)
			original = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.c.Host.Rank+_rankURL).Reply(200).JSON(`{"code":0,"num":8,"list":[{"aid":33986715,"score":704},{"aid":33913315,"score":571}]}`)
			res, err := d.RankingRegion(c, rid, day, original)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRankingRecommend(t *testing.T) {
	convey.Convey("RankingRecommend", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int16(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.c.Host.Rank+_rankURL).Reply(200).JSON(`{"code":0,"num":8,"list":[{"aid":33986715,"score":704},{"aid":33913315,"score":571}]}`)
			res, err := d.RankingRecommend(c, rid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRankingTag(t *testing.T) {
	convey.Convey("RankingTag", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int16(0)
			tagID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.c.Host.Rank).Reply(200).JSON(`{"code":0,"num":8,"list":[{"aid":33986715,"score":704},{"aid":33913315,"score":571}]}`)
			rs, err := d.RankingTag(c, rid, tagID)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRegionCustom(t *testing.T) {
	convey.Convey("RegionCustom", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.customURL).Reply(200).JSON(`{"code":0,"num":8,"list":[{"aid":33986715},{"aid":33913315}]}`)
			rs, err := d.RegionCustom(c)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaorankURI(t *testing.T) {
	convey.Convey("rankURI", t, func(ctx convey.C) {
		var (
			rid      = int16(0)
			rankType = ""
			day      = int(0)
			arcType  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := rankURI(rid, rankType, day, arcType)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
