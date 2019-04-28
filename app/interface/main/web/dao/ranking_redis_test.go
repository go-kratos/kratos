package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyRkList(t *testing.T) {
	convey.Convey("keyRkList", t, func(ctx convey.C) {
		var (
			rid      = int16(1)
			rankType = int(0)
			day      = int(1)
			arcType  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkList(rid, rankType, day, arcType)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkListBak(t *testing.T) {
	convey.Convey("keyRkListBak", t, func(ctx convey.C) {
		var (
			rid      = int16(0)
			rankType = int(0)
			day      = int(0)
			arcType  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkListBak(rid, rankType, day, arcType)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkIndex(t *testing.T) {
	convey.Convey("keyRkIndex", t, func(ctx convey.C) {
		var (
			day = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkIndex(day)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkIndexBak(t *testing.T) {
	convey.Convey("keyRkIndexBak", t, func(ctx convey.C) {
		var (
			day = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkIndexBak(day)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkRegionList(t *testing.T) {
	convey.Convey("keyRkRegionList", t, func(ctx convey.C) {
		var (
			rid      = int16(1)
			day      = int(3)
			original = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkRegionList(rid, day, original)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkRegionListBak(t *testing.T) {
	convey.Convey("keyRkRegionListBak", t, func(ctx convey.C) {
		var (
			rid      = int16(1)
			day      = int(3)
			original = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkRegionListBak(rid, day, original)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkRecommendList(t *testing.T) {
	convey.Convey("keyRkRecommendList", t, func(ctx convey.C) {
		var (
			rid = int16(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkRecommendList(rid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkRecommendListBak(t *testing.T) {
	convey.Convey("keyRkRecommendListBak", t, func(ctx convey.C) {
		var (
			rid = int16(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkRecommendListBak(rid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkTagList(t *testing.T) {
	convey.Convey("keyRkTagList", t, func(ctx convey.C) {
		var (
			rid   = int16(1)
			tagID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkTagList(rid, tagID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyRkTagListBak(t *testing.T) {
	convey.Convey("keyRkTagListBak", t, func(ctx convey.C) {
		var (
			rid   = int16(1)
			tagID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRkTagListBak(rid, tagID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRankingCache(t *testing.T) {
	convey.Convey("RankingCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(1)
			rankType = int(1)
			day      = int(3)
			arcType  = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RankingCache(c, rid, rankType, day, arcType)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestDaoRankingBakCache(t *testing.T) {
	convey.Convey("RankingBakCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(1)
			rankType = int(1)
			day      = int(3)
			arcType  = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.RankingBakCache(c, rid, rankType, day, arcType)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestDaoRankingIndexCache(t *testing.T) {
	convey.Convey("RankingIndexCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			day = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingIndexCache(c, day)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingIndexBakCache(t *testing.T) {
	convey.Convey("RankingIndexBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			day = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingIndexBakCache(c, day)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingRegionCache(t *testing.T) {
	convey.Convey("RankingRegionCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(1)
			day      = int(3)
			original = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingRegionCache(c, rid, day, original)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingRegionBakCache(t *testing.T) {
	convey.Convey("RankingRegionBakCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(1)
			day      = int(3)
			original = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingRegionBakCache(c, rid, day, original)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingRecommendCache(t *testing.T) {
	convey.Convey("RankingRecommendCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int16(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingRecommendCache(c, rid)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingRecommendBakCache(t *testing.T) {
	convey.Convey("RankingRecommendBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int16(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingRecommendBakCache(c, rid)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingTagCache(t *testing.T) {
	convey.Convey("RankingTagCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int16(0)
			tagID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingTagCache(c, rid, tagID)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoRankingTagBakCache(t *testing.T) {
	convey.Convey("RankingTagBakCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int16(0)
			tagID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.RankingTagBakCache(c, rid, tagID)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoSetRegionCustomCache(t *testing.T) {
	convey.Convey("SetRegionCustomCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = []*model.Custom{{Aid: 1111}, {Aid: 2222}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRegionCustomCache(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRegionCustomCache(t *testing.T) {
	convey.Convey("RegionCustomCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RegionCustomCache(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRegionCustomBakCache(t *testing.T) {
	convey.Convey("RegionCustomBakCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RegionCustomBakCache(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetRankingCache(t *testing.T) {
	convey.Convey("SetRankingCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(0)
			rankType = int(0)
			day      = int(0)
			arcType  = int(0)
			data     = &model.RankData{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRankingCache(c, rid, rankType, day, arcType, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetRankingIndexCache(t *testing.T) {
	convey.Convey("SetRankingIndexCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			day  = int(0)
			arcs = []*model.IndexArchive{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRankingIndexCache(c, day, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetRankingRegionCache(t *testing.T) {
	convey.Convey("SetRankingRegionCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			rid      = int16(0)
			day      = int(0)
			original = int(0)
			arcs     = []*model.RegionArchive{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRankingRegionCache(c, rid, day, original, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetRankingRecommendCache(t *testing.T) {
	convey.Convey("SetRankingRecommendCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rid  = int16(0)
			arcs = []*model.IndexArchive{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRankingRecommendCache(c, rid, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetRankingTagCache(t *testing.T) {
	convey.Convey("SetRankingTagCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int16(0)
			tagID = int64(0)
			arcs  = []*model.TagArchive{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRankingTagCache(c, rid, tagID, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
