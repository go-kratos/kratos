package archive

import (
	"context"
	"go-common/app/service/main/archive/api"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivergAllKey(t *testing.T) {
	var (
		rid = int16(97)
	)
	convey.Convey("rgAllKey", t, func(ctx convey.C) {
		p1 := rgAllKey(rid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivergOriginKey(t *testing.T) {
	var (
		rid = int16(97)
	)
	convey.Convey("rgOriginKey", t, func(ctx convey.C) {
		p1 := rgOriginKey(rid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivergTopKey(t *testing.T) {
	var (
		rid = int16(97)
	)
	convey.Convey("rgTopKey", t, func(ctx convey.C) {
		p1 := rgTopKey(rid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveAddRegionArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		rid  = int16(97)
		reid = int16(1)
		as   = &api.RegionArc{}
	)
	convey.Convey("AddRegionArcCache", t, func(ctx convey.C) {
		err := d.AddRegionArcCache(c, rid, reid, as)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveRegionTopArcsCache(t *testing.T) {
	var (
		c     = context.TODO()
		reid  = int16(1)
		start = int(0)
		end   = int(20)
	)
	convey.Convey("RegionTopArcsCache", t, func(ctx convey.C) {
		_, err := d.RegionTopArcsCache(c, reid, start, end)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveRegionArcsCache(t *testing.T) {
	var (
		c     = context.TODO()
		rid   = int16(97)
		start = int(0)
		end   = int(20)
	)
	convey.Convey("RegionArcsCache", t, func(ctx convey.C) {
		aids, err := d.RegionArcsCache(c, rid, start, end)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(aids, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveRegionOriginArcsCache(t *testing.T) {
	var (
		c     = context.TODO()
		rid   = int16(97)
		start = int(0)
		end   = int(20)
	)
	convey.Convey("RegionOriginArcsCache", t, func(ctx convey.C) {
		_, err := d.RegionOriginArcsCache(c, rid, start, end)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivezrange(t *testing.T) {
	var (
		c     = context.TODO()
		key   = "1_o"
		start = int(1)
		end   = int(10)
	)
	convey.Convey("zrange", t, func(ctx convey.C) {
		_, err := d.zrange(c, key, start, end)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveRegionTopCountCache(t *testing.T) {
	var (
		c     = context.TODO()
		reids = []int16{1}
		min   = int64(0)
		max   = int64(10)
	)
	convey.Convey("RegionTopCountCache", t, func(ctx convey.C) {
		recm, err := d.RegionTopCountCache(c, reids, min, max)
		ctx.Convey("Then err should be nil.recm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(recm, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveRegionAllCountCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("RegionAllCountCache", t, func(ctx convey.C) {
		count, err := d.RegionAllCountCache(c)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveRegionCountCache(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int16(97)
	)
	convey.Convey("RegionCountCache", t, func(ctx convey.C) {
		count, err := d.RegionCountCache(c, rid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveRegionOriginCountCache(t *testing.T) {
	var (
		c   = context.TODO()
		rid = int16(97)
	)
	convey.Convey("RegionOriginCountCache", t, func(ctx convey.C) {
		count, err := d.RegionOriginCountCache(c, rid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveDelRegionArcCache(t *testing.T) {
	var (
		c    = context.TODO()
		rid  = int16(97)
		reid = int16(1)
		aid  = int64(1)
	)
	convey.Convey("DelRegionArcCache", t, func(ctx convey.C) {
		err := d.DelRegionArcCache(c, rid, reid, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
