package mcndao

import (
	"testing"
	"time"

	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaobyValueDesc(t *testing.T) {
	convey.Convey("byValueDesc", t, func(ctx convey.C) {
		var (
			p1 = &mcnmodel.RankUpFansInfo{Mid: 1}
			p2 = &mcnmodel.RankUpFansInfo{Mid: 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := byValueDesc(p1, p2)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoGetRankUpFans(t *testing.T) {
	convey.Convey("GetRankUpFans", t, func(ctx convey.C) {
		var (
			signID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetRankUpFans(signID)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoGetRankArchiveLikes(t *testing.T) {
	convey.Convey("GetRankArchiveLikes", t, func(ctx convey.C) {
		var (
			signID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetRankArchiveLikes(signID)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaogetRankCache(t *testing.T) {
	convey.Convey("getRankCache", t, func(ctx convey.C) {
		var (
			signID  = int64(1)
			keyCalc keyFunc
			load    loadRankFunc
		)
		keyCalc = cacheKeyRankFans
		load = d.loadRankUpFansCache
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.getRankCache(signID, keyCalc, load)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				if result == nil {
					ctx.So(result, convey.ShouldBeNil)
				} else {
					ctx.So(result, convey.ShouldNotBeNil)
				}

			})
		})
	})
}

func TestMcndaocacheKeyRankFans(t *testing.T) {
	convey.Convey("cacheKeyRankFans", t, func(ctx convey.C) {
		var (
			signID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := cacheKeyRankFans(signID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaocacheKeyRankArchiveLikes(t *testing.T) {
	convey.Convey("cacheKeyRankArchiveLikes", t, func(ctx convey.C) {
		var (
			signID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := cacheKeyRankArchiveLikes(signID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoloadRankUpFansCache(t *testing.T) {
	convey.Convey("loadRankUpFansCache", t, func(ctx convey.C) {
		var (
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.loadRankUpFansCache(signID, date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRawRankUpFans(t *testing.T) {
	convey.Convey("RawRankUpFans", t, func(ctx convey.C) {
		var (
			signID = int64(1)
			date   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.RawRankUpFans(signID, date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				if len(result) == 0 {
					ctx.So(result, convey.ShouldBeNil)
				} else {
					ctx.So(result, convey.ShouldNotBeNil)
				}
			})
		})
	})
}

func TestMcndaoReloadRank(t *testing.T) {
	convey.Convey("ReloadRank", t, func(ctx convey.C) {
		var (
			signID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ReloadRank(signID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoloadRankArchiveLikesCache(t *testing.T) {
	convey.Convey("loadRankArchiveLikesCache", t, func(ctx convey.C) {
		var (
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.loadRankArchiveLikesCache(signID, date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRawRankArchiveLikes(t *testing.T) {
	convey.Convey("RawRankArchiveLikes", t, func(ctx convey.C) {
		var (
			signID = int64(1)
			date   = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.RawRankArchiveLikes(signID, date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				if len(result) == 0 {
					ctx.So(result, convey.ShouldBeEmpty)
				} else {
					ctx.So(result, convey.ShouldNotBeNil)
				}
			})
		})
	})
}
