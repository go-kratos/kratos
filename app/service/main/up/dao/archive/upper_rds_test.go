package archive

import (
	"context"
	"testing"

	"go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveupCntKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("upCntKey", t, func(ctx convey.C) {
		p1 := upCntKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveupPasKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("upPasKey", t, func(ctx convey.C) {
		p1 := upPasKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveAddUpperCountCache(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(1)
		count = int64(10)
	)
	convey.Convey("AddUpperCountCache", t, func(ctx convey.C) {
		err := d.AddUpperCountCache(c, mid, count)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveUpperCountCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("UpperCountCache", t, func(ctx convey.C) {
		count, err := d.UpperCountCache(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUppersCountCache(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{1, 2}
	)
	convey.Convey("UppersCountCache", t, func(ctx convey.C) {
		cached, missed, err := d.UppersCountCache(c, mids)
		ctx.Convey("Then err should be nil.cached,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missed, convey.ShouldNotBeNil)
			ctx.So(cached, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpperPassedCache(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(1)
		start = int(0)
		end   = int(10)
	)
	convey.Convey("UpperPassedCache", t, func(ctx convey.C) {
		_, err := d.UpperPassedCache(c, mid, start, end)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveUppersPassedCacheWithScore(t *testing.T) {
	var (
		c     = context.TODO()
		mids  = []int64{1, 2}
		start = int(0)
		end   = int(10)
	)
	convey.Convey("UppersPassedCacheWithScore", t, func(ctx convey.C) {
		aidm, err := d.UppersPassedCacheWithScore(c, mids, start, end)
		ctx.Convey("Then err should be nil.aidm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(aidm, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUppersPassedCache(t *testing.T) {
	var (
		c     = context.TODO()
		mids  = []int64{1, 2}
		start = int(0)
		end   = int(10)
	)
	convey.Convey("UppersPassedCache", t, func(ctx convey.C) {
		aidm, err := d.UppersPassedCache(c, mids, start, end)
		ctx.Convey("Then err should be nil.aidm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(aidm, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveExpireUpperPassedCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("ExpireUpperPassedCache", t, func(ctx convey.C) {
		ok, err := d.ExpireUpperPassedCache(c, mid)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveExpireUppersCountCache(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{1, 2}
	)
	convey.Convey("ExpireUppersCountCache", t, func(ctx convey.C) {
		cachedUp, missed, err := d.ExpireUppersCountCache(c, mids)
		ctx.Convey("Then err should be nil.cachedUp,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(missed, convey.ShouldNotBeNil)
			ctx.So(cachedUp, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveExpireUppersPassedCache(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{1, 2}
	)
	convey.Convey("ExpireUppersPassedCache", t, func(ctx convey.C) {
		res, err := d.ExpireUppersPassedCache(c, mids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveAddUpperPassedCache(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		aids       = []int64{1}
		ptimes     = []time.Time{1531553274}
		copyrights = []int8{1}
	)
	convey.Convey("AddUpperPassedCache", t, func(ctx convey.C) {
		err := d.AddUpperPassedCache(c, mid, aids, ptimes, copyrights)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveDelUpperPassedCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		aid = int64(1)
	)
	convey.Convey("DelUpperPassedCache", t, func(ctx convey.C) {
		err := d.DelUpperPassedCache(c, mid, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
