package mcndao

import (
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaoadd(t *testing.T) {
	convey.Convey("add", t, func(ctx convey.C) {
		var (
			v = &mcnmodel.McnGetRecommendPoolInfo{}
			s = &RecommendPoolCache{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.add(v)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestMcndaoLen(t *testing.T) {
	var (
		s = &RecommendDataSorter{}
	)
	convey.Convey("Len", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := s.Len()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRecommendSortByFansDesc(t *testing.T) {
	convey.Convey("RecommendSortByFansDesc", t, func(ctx convey.C) {
		var (
			p1 = &mcnmodel.McnGetRecommendPoolInfo{}
			p2 = &mcnmodel.McnGetRecommendPoolInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := RecommendSortByFansDesc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRecommendSortByFansAsc(t *testing.T) {
	convey.Convey("RecommendSortByFansAsc", t, func(ctx convey.C) {
		var (
			p1 = &mcnmodel.McnGetRecommendPoolInfo{}
			p2 = &mcnmodel.McnGetRecommendPoolInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := RecommendSortByFansAsc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRecommendSortByMonthFansDesc(t *testing.T) {
	convey.Convey("RecommendSortByMonthFansDesc", t, func(ctx convey.C) {
		var (
			p1 = &mcnmodel.McnGetRecommendPoolInfo{}
			p2 = &mcnmodel.McnGetRecommendPoolInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := RecommendSortByMonthFansDesc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRecommendSortByArchiveCountDesc(t *testing.T) {
	convey.Convey("RecommendSortByArchiveCountDesc", t, func(ctx convey.C) {
		var (
			p1 = &mcnmodel.McnGetRecommendPoolInfo{}
			p2 = &mcnmodel.McnGetRecommendPoolInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := RecommendSortByArchiveCountDesc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaocacheKeyRecommend(t *testing.T) {
	convey.Convey("cacheKeyRecommend", t, func(ctx convey.C) {
		var (
			a1 = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := cacheKeyRecommend(a1)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaorawGetRecommendPool(t *testing.T) {
	convey.Convey("rawGetRecommendPool", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.rawGetRecommendPool()
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(res), convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestMcndaoloadRecommendPool(t *testing.T) {
	convey.Convey("loadRecommendPool", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.loadRecommendPool()
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaogetRecommendCache(t *testing.T) {
	convey.Convey("getRecommendCache", t, func(ctx convey.C) {
		var (
			keyCalc = cacheKeyRecommend
			load    = d.loadRecommendPool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			result, err := d.getRecommendCache(keyCalc, load)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoGetRecommendPool(t *testing.T) {
	convey.Convey("GetRecommendPool", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetRecommendPool()
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
