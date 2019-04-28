package dao

import (
	"context"
	"go-common/app/service/main/relation/model"
	xtime "go-common/library/time"
	"testing"
	gtime "time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaofollowingsKey(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("followingsKey", t, func(cv convey.C) {
		p1 := followingsKey(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaomonitorKey(t *testing.T) {
	convey.Convey("monitorKey", t, func(cv convey.C) {
		p1 := monitorKey()
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorecentFollower(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("recentFollower", t, func(cv convey.C) {
		p1 := recentFollower(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorecentFollowerNotify(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("recentFollowerNotify", t, func(cv convey.C) {
		p1 := recentFollowerNotify(mid)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaodailyNotifyCount(t *testing.T) {
	var (
		mid  = int64(1)
		date gtime.Time
	)
	convey.Convey("dailyNotifyCount", t, func(cv convey.C) {
		p1 := dailyNotifyCount(mid, date)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaopingRedis(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("pingRedis", t, func(cv convey.C) {
		err := d.pingRedis(c)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFollowingsCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("FollowingsCache", t, func(cv convey.C) {
		err := d.SetFollowingsCache(c, mid, []*model.Following{
			{Mid: 2},
		})
		cv.Convey("SetFollowingsCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})

		err = d.AddFollowingCache(c, mid, &model.Following{Mid: 2})
		cv.Convey("AddFollowingCache; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})

		followings, err := d.FollowingsCache(c, mid)
		cv.Convey("FollowingsCache; Then err should be nil.followings should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(followings, convey.ShouldNotBeNil)
		})

		err = d.DelFollowing(c, mid, &model.Following{Mid: 2})
		cv.Convey("DelFollowing; Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelFollowingsCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelFollowingsCache", t, func(cv convey.C) {
		err := d.DelFollowingsCache(c, mid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRelationsCache(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		fids = []int64{2}
	)
	convey.Convey("RelationsCache", t, func(cv convey.C) {
		resMap, err := d.RelationsCache(c, mid, fids)
		cv.Convey("Then err should be nil.resMap should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(resMap, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoencode(t *testing.T) {
	var (
		attribute = uint32(0)
		mtime     = xtime.Time(int64(0))
		tagids    = []int64{}
		special   = int32(0)
	)
	convey.Convey("encode", t, func(cv convey.C) {
		res, err := d.encode(attribute, mtime, tagids, special)
		cv.Convey("Then err should be nil.res should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaodecode(t *testing.T) {
	var (
		src = []byte("")
		v   = &model.FollowingTags{}
	)
	convey.Convey("decode", t, func(cv convey.C) {
		err := d.decode(src, v)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoMonitorCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("MonitorCache", t, func(cv convey.C) {
		exist, err := d.MonitorCache(c, mid)
		cv.Convey("Then err should be nil.exist should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(exist, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetMonitorCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("SetMonitorCache", t, func(cv convey.C) {
		err := d.SetMonitorCache(c, mid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelMonitorCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelMonitorCache", t, func(cv convey.C) {
		err := d.DelMonitorCache(c, mid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoLoadMonitorCache(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{}
	)
	convey.Convey("LoadMonitorCache", t, func(cv convey.C) {
		err := d.LoadMonitorCache(c, mids)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTodayNotifyCountCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("TodayNotifyCountCache", t, func(cv convey.C) {
		notifyCount, err := d.TodayNotifyCountCache(c, mid)
		cv.Convey("Then err should be nil.notifyCount should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(notifyCount, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoIncrTodayNotifyCountCache(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("IncrTodayNotifyCount", t, func(cv convey.C) {
		err := d.IncrTodayNotifyCount(c, mid)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}
