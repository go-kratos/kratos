package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoexpKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("expKey", t, func(ctx convey.C) {
		p1 := expKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaokeyInfo(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("keyInfo", t, func(ctx convey.C) {
		p1 := keyInfo(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaomcBaseKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("mcBaseKey", t, func(ctx convey.C) {
		key := d.mcBaseKey(mid)
		ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
			ctx.So(key, convey.ShouldNotBeNil)
		})
	})
}

func TestDaomoralKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("moralKey", t, func(ctx convey.C) {
		key := d.moralKey(mid)
		ctx.Convey("Then key should not be nil.", func(ctx convey.C) {
			ctx.So(key, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelMoralCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelMoralCache", t, func(ctx convey.C) {
		err := d.DelMoralCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelBaseInfoCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelBaseInfoCache", t, func(ctx convey.C) {
		err := d.DelBaseInfoCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelInfoCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelInfoCache", t, func(ctx convey.C) {
		err := d.DelInfoCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaorealnameApplyStatusKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("realnameApplyStatusKey", t, func(ctx convey.C) {
		p1 := realnameApplyStatusKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDeleteRealnameCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DeleteRealnameCache", t, func(ctx convey.C) {
		err := d.DeleteRealnameCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetExpCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
		exp = int64(0)
	)
	convey.Convey("SetExpCache", t, func(ctx convey.C) {
		err := d.SetExpCache(c, mid, exp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelExpCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelExpCache", t, func(ctx convey.C) {
		err := d.DelExpCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
