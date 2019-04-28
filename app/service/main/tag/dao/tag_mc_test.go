package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyTag(t *testing.T) {
	convey.Convey("keyTag", t, func(ctx convey.C) {
		var (
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyTag(tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyName(t *testing.T) {
	convey.Convey("keyName", t, func(ctx convey.C) {
		var (
			name = "搞笑"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyName(name)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyNames(t *testing.T) {
	convey.Convey("keyNames", t, func(ctx convey.C) {
		var (
			oid = int64(28843596)
			typ = int8(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyNames(oid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyCount(t *testing.T) {
	convey.Convey("keyCount", t, func(ctx convey.C) {
		var (
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyCount(tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagCache(t *testing.T) {
	convey.Convey("TagCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TagCache(c, tid)
		})
	})
}

func TestDaoTagsCaches(t *testing.T) {
	convey.Convey("TagsCaches", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1833}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TagsCaches(c, tids)
		})
	})
}

func TestDaoTagMapCaches(t *testing.T) {
	convey.Convey("TagMapCaches", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, missed, err := d.TagMapCaches(c, tids)
			ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagCacheByName(t *testing.T) {
	convey.Convey("TagCacheByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "搞笑"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TagCacheByName(c, name)
		})
	})
}

func TestDaoTagCachesByNames(t *testing.T) {
	convey.Convey("TagCachesByNames", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			names = []string{"搞笑", "数码"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TagCachesByNames(c, names)
		})
	})
}

func TestDaoAddTagCache(t *testing.T) {
	convey.Convey("AddTagCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tag = &model.Tag{
				Name: "unit 233",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddTagCache(c, tag)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddTagsCache(t *testing.T) {
	convey.Convey("AddTagsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tags = []*model.Tag{
				{
					Name: "unit 233",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddTagsCache(c, tags)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelTagCache(t *testing.T) {
	convey.Convey("DelTagCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelTagCache(c, tid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCountCache(t *testing.T) {
	convey.Convey("AddCountCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			count = &model.Count{
				Tid:  1844,
				Bind: 123,
				Sub:  456,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCountCache(c, count)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCountsCache(t *testing.T) {
	convey.Convey("AddCountsCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			counts = []*model.Count{
				{
					Tid:  1844,
					Bind: 123,
					Sub:  456,
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCountsCache(c, counts)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountCache(t *testing.T) {
	convey.Convey("CountCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.CountCache(c, tid)
		})
	})
}

func TestDaoCountMapCache(t *testing.T) {
	convey.Convey("CountMapCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, missed, err := d.CountMapCache(c, tids)
			ctx.Convey("Then err should be nil.res,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCountCache(t *testing.T) {
	convey.Convey("DelCountCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCountCache(c, tid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCountsCache(t *testing.T) {
	convey.Convey("DelCountsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1833, 12096}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCountsCache(c, tids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTagNamesCache(t *testing.T) {
	convey.Convey("TagNamesCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(28843596)
			typ = int8(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TagNamesCache(c, oid, typ)
		})
	})
}
