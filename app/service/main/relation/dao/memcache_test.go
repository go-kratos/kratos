package dao

import (
	"context"
	"go-common/app/service/main/relation/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaostatKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("statKey", t, func(cv convey.C) {
		p1 := statKey(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagsKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("tagsKey", t, func(cv convey.C) {
		p1 := tagsKey(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaofollowingKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("followingKey", t, func(cv convey.C) {
		p1 := followingKey(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaofollowerKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("followerKey", t, func(cv convey.C) {
		p1 := followerKey(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaotagCountKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("tagCountKey", t, func(cv convey.C) {
		p1 := tagCountKey(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoglobalHotKey(t *testing.T) {
	convey.Convey("globalHotKey", t, func(cv convey.C) {
		p1 := globalHotKey()
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaofollowerNotifySetting(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("followerNotifySetting", t, func(cv convey.C) {
		p1 := followerNotifySetting(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaopingMC(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("pingMC", t, func(cv convey.C) {
		err := d.pingMC(c)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetFollowingCache(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(1)
		followings = []*model.Following{}
	)
	convey.Convey("SetFollowingCache", t, func(cv convey.C) {
		err := d.SetFollowingCache(c, mid, followings)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFollowingCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("FollowingCache", t, func(cv convey.C) {
		followings, err := d.FollowingCache(c, mid)
		cv.Convey("Then err should be nil.followings should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(followings, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelFollowingCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelFollowingCache", t, func(cv convey.C) {
		err := d.DelFollowingCache(c, mid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetFollowerCache(t *testing.T) {
	var (
		c         = context.Background()
		mid       = int64(1)
		followers = []*model.Following{}
	)
	convey.Convey("SetFollowerCache", t, func(cv convey.C) {
		err := d.SetFollowerCache(c, mid, followers)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFollowerCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("FollowerCache", t, func(cv convey.C) {
		followers, err := d.FollowerCache(c, mid)
		cv.Convey("Then err should be nil.followers should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(followers, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelFollowerCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("DelFollowerCache", t, func(cv convey.C) {
		err := d.DelFollowerCache(c, mid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaosetFollowingCache(t *testing.T) {
	var (
		c          = context.Background()
		key        = followingKey(1)
		followings = []*model.Following{}
	)
	convey.Convey("setFollowingCache", t, func(cv convey.C) {
		err := d.setFollowingCache(c, key, followings)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaofollowingCache(t *testing.T) {
	var (
		c   = context.Background()
		key = followingKey(1)
	)
	convey.Convey("followingCache", t, func(cv convey.C) {
		followings, err := d.followingCache(c, key)
		cv.Convey("Then err should be nil.followings should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(followings, convey.ShouldNotBeNil)
		})
	})
}

func TestDaodelFollowingCache(t *testing.T) {
	var (
		c   = context.Background()
		key = followingKey(1)
	)
	convey.Convey("delFollowingCache", t, func(cv convey.C) {
		err := d.delFollowingCache(c, key)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTagCountCache(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		tagCount = []*model.TagCount{
			{
				Tagid: 1,
				Name:  "test",
				Count: 1,
			},
		}
	)
	convey.Convey("TagCountCache", t, func(cv convey.C) {
		err := d.SetTagCountCache(c, mid, tagCount)
		cv.So(err, convey.ShouldBeNil)

		tagCount, err := d.TagCountCache(c, mid)
		cv.So(err, convey.ShouldBeNil)
		cv.So(tagCount, convey.ShouldNotBeNil)

		err = d.DelTagCountCache(c, mid)
		cv.So(err, convey.ShouldBeNil)
	})
}

func TestDaoTagsCache(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		tags = &model.Tags{
			Tags: map[int64]*model.Tag{
				1: {
					Id:     1,
					Name:   "1",
					Status: 1,
				},
			},
		}
	)
	convey.Convey("TagsCache", t, func(cv convey.C) {
		err := d.SetTagsCache(c, mid, tags)
		cv.Convey("SetTagsCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})

		tags, err := d.TagsCache(c, mid)
		cv.Convey("TagsCache; Then err should be nil.tags should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(tags, convey.ShouldNotBeNil)
		})

		err = d.DelTagsCache(c, mid)
		cv.Convey("DelTagsCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetGlobalHotRecCache(t *testing.T) {
	var (
		c    = context.Background()
		fids = []int64{1}
	)
	convey.Convey("SetGlobalHotRecCache", t, func(cv convey.C) {
		err := d.SetGlobalHotRecCache(c, fids)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGlobalHotRecCache(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("GlobalHotRecCache", t, func(cv convey.C) {
		fs, err := d.GlobalHotRecCache(c)
		cv.Convey("Then err should be nil.fs should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(fs, convey.ShouldNotBeNil)
		})
	})
}

func TestDaostatCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		st  = &model.Stat{
			Mid:       1,
			Follower:  1,
			Following: 1,
			Black:     1,
			Whisper:   1,
		}
	)
	convey.Convey("statCache", t, func(cv convey.C) {
		err := d.SetStatCache(c, mid, st)
		cv.Convey("SetStatCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})

		s1, err := d.statCache(c, mid)
		cv.Convey("statCache; Then err should be nil.p1 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(s1, convey.ShouldNotBeNil)
		})

		p1, p2, err := d.statsCache(c, []int64{1})
		cv.Convey("statsCache; Then err should be nil.p1,p2 should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(p2, convey.ShouldNotBeNil)
			cv.So(p1, convey.ShouldNotBeNil)
		})

		err = d.DelStatCache(c, mid)
		cv.Convey("DelStatCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetFollowerNotifyCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		val = &model.FollowerNotifySetting{
			Mid:     1,
			Enabled: true,
		}
	)
	convey.Convey("GetFollowerNotifyCache", t, func(cv convey.C) {
		err := d.SetFollowerNotifyCache(c, mid, val)
		cv.Convey("SetFollowerNotifyCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})

		res, err := d.GetFollowerNotifyCache(c, mid)
		cv.Convey("GetFollowerNotifyCache; Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})

		err = d.DelFollowerNotifyCache(c, mid)
		cv.Convey("DelFollowerNotifyCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}
