package like

import (
	"context"
	likemdl "go-common/app/interface/main/activity/model/like"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeCacheLike(t *testing.T) {
	convey.Convey("CacheLike", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheLike(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestLikeCacheLikes(t *testing.T) {
	convey.Convey("CacheLikes", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheLikes(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestLikeAddCacheLike(t *testing.T) {
	convey.Convey("AddCacheLike", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &likemdl.Item{ID: 0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheLike(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActSubject(t *testing.T) {
	convey.Convey("CacheActSubject", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActSubject(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestLikeAddCacheActSubject(t *testing.T) {
	convey.Convey("AddCacheActSubject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &likemdl.SubjectItem{ID: 0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActSubject(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActSubjectMaxID(t *testing.T) {
	convey.Convey("CacheActSubjectMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActSubjectMaxID(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheActSubjectMaxID(t *testing.T) {
	convey.Convey("AddCacheActSubjectMaxID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActSubjectMaxID(c, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheLikeMaxID(t *testing.T) {
	convey.Convey("CacheLikeMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheLikeMaxID(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheLikeMaxID(t *testing.T) {
	convey.Convey("AddCacheLikeMaxID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = int64(10586)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheLikeMaxID(c, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheLikeMissionBuff(t *testing.T) {
	convey.Convey("CacheLikeMissionBuff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheLikeMissionBuff(c, id, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheLikeMissionBuff(t *testing.T) {
	convey.Convey("AddCacheLikeMissionBuff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			val = int64(1)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheLikeMissionBuff(c, id, val, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheMissionGroupItems(t *testing.T) {
	convey.Convey("CacheMissionGroupItems", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMissionGroupItems(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheMissionGroupItems(t *testing.T) {
	convey.Convey("AddCacheMissionGroupItems", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = map[int64]*likemdl.MissionGroup{1: {ID: 1, Sid: 1}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMissionGroupItems(c, values)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActMission(t *testing.T) {
	convey.Convey("CacheActMission", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			lid = int64(77)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActMission(c, id, lid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheActMission(t *testing.T) {
	convey.Convey("AddCacheActMission", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			val = int64(1)
			lid = int64(77)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActMission(c, id, val, lid, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActLikeAchieves(t *testing.T) {
	convey.Convey("CacheActLikeAchieves", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActLikeAchieves(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeAddCacheActLikeAchieves(t *testing.T) {
	convey.Convey("AddCacheActLikeAchieves", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &likemdl.Achievements{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActLikeAchieves(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActMissionFriends(t *testing.T) {
	convey.Convey("CacheActMissionFriends", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			lid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActMissionFriends(c, id, lid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeDelCacheActMissionFriends(t *testing.T) {
	convey.Convey("DelCacheActMissionFriends", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			lid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheActMissionFriends(c, id, lid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeAddCacheActMissionFriends(t *testing.T) {
	convey.Convey("AddCacheActMissionFriends", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &likemdl.ActMissionGroups{}
			lid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActMissionFriends(c, id, val, lid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActUserAchieve(t *testing.T) {
	convey.Convey("CacheActUserAchieve", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActUserAchieve(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeAddCacheActUserAchieve(t *testing.T) {
	convey.Convey("AddCacheActUserAchieve", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &likemdl.ActLikeUserAchievement{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActUserAchieve(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheActUserAward(t *testing.T) {
	convey.Convey("CacheActUserAward", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActUserAward(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheActUserAward(t *testing.T) {
	convey.Convey("AddCacheActUserAward", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			val = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActUserAward(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheSubjectStat(t *testing.T) {
	convey.Convey("CacheSubjectStat", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheSubjectStat(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeAddCacheSubjectStat(t *testing.T) {
	convey.Convey("AddCacheSubjectStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &likemdl.SubjectStat{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheSubjectStat(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheViewRank(t *testing.T) {
	convey.Convey("CacheViewRank", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheViewRank(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheViewRank(t *testing.T) {
	convey.Convey("AddCacheViewRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheViewRank(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheLikeContent(t *testing.T) {
	convey.Convey("CacheLikeContent", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheLikeContent(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeAddCacheLikeContent(t *testing.T) {
	convey.Convey("AddCacheLikeContent", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values map[int64]*likemdl.LikeContent
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheLikeContent(c, values)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeAddCacheSourceItemData(t *testing.T) {
	convey.Convey("AddCacheSourceItemData", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			sid    = int64(213)
			values = []int64{10884, 10883, 10882, 10881, 10880}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheSourceItemData(c, sid, values)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCacheSourceItemData(t *testing.T) {
	convey.Convey("CacheSourceItemData", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(213)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheSourceItemData(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}
