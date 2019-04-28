package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoexpKey(t *testing.T) {
	convey.Convey("expKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := expKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaomcBaseKey(t *testing.T) {
	convey.Convey("mcBaseKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key := d.mcBaseKey(mid)
			convCtx.Convey("Then key should not be nil.", func(convCtx convey.C) {
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaomoralKey(t *testing.T) {
	convey.Convey("moralKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key := d.moralKey(mid)
			convCtx.Convey("Then key should not be nil.", func(convCtx convey.C) {
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetStartCache(t *testing.T) {
	convey.Convey("SetStartCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			key = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetStartCache(c, mid, key)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelMoralCache(t *testing.T) {
	convey.Convey("DelMoralCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelMoralCache(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelBaseInfoCache(t *testing.T) {
	convey.Convey("DelBaseInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelBaseInfoCache(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetExpCache(t *testing.T) {
	convey.Convey("SetExpCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			exp = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetExpCache(c, mid, exp)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaorealnameInfoKey(t *testing.T) {
	convey.Convey("realnameInfoKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := realnameInfoKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaorealnameApplyStatusKey(t *testing.T) {
	convey.Convey("realnameApplyStatusKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := realnameApplyStatusKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeleteRealnameCache(t *testing.T) {
	convey.Convey("DeleteRealnameCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DeleteRealnameCache(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
