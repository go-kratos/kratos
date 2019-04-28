package dao

import (
	"context"
	relation "go-common/app/service/main/relation/model"
	xtime "go-common/library/time"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaostatKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("statKey", t, func(ctx convey.C) {
		p1 := statKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagsKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("tagsKey", t, func(ctx convey.C) {
		p1 := tagsKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaofollowingsKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("followingsKey", t, func(ctx convey.C) {
		p1 := followingsKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaofollowingKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("followingKey", t, func(ctx convey.C) {
		p1 := followingKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagCountKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("tagCountKey", t, func(ctx convey.C) {
		p1 := tagCountKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelStatCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelStatCache", t, func(ctx convey.C) {
		err := d.DelStatCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelFollowerCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelFollowerCache", t, func(ctx convey.C) {
		err := d.DelFollowerCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelFollowing(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(1)
		following = &relation.Following{}
	)
	convey.Convey("DelFollowing", t, func(ctx convey.C) {
		err := d.DelFollowing(c, mid, following)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelTagsCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelTagsCache", t, func(ctx convey.C) {
		err := d.DelTagsCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddFollowingCache(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(1)
		following = &relation.Following{}
	)
	convey.Convey("AddFollowingCache", t, func(ctx convey.C) {
		err := d.AddFollowingCache(c, mid, following)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoencode(t *testing.T) {
	var (
		attribute = uint32(1)
		mtime     xtime.Time
		tagids    = []int64{1}
		special   = int32(1)
	)
	convey.Convey("encode", t, func(ctx convey.C) {
		res, err := d.encode(attribute, mtime, tagids, special)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelFollowingCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelFollowingCache", t, func(ctx convey.C) {
		err := d.DelFollowingCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaodelFollowingCache(t *testing.T) {
	var (
		c   = context.Background()
		key = ""
	)
	convey.Convey("delFollowingCache", t, func(ctx convey.C) {
		err := d.delFollowingCache(c, key)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelTagCountCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelTagCountCache", t, func(ctx convey.C) {
		err := d.DelTagCountCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
