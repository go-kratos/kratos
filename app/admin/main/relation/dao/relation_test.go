package dao

import (
	"context"

	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestFollowers(t *testing.T) {
	convey.Convey("Followers", t, func() {
		rl, err := d.Followers(context.TODO(), 1, 2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rl, convey.ShouldNotBeNil)
	})
}

func TestFollowings(t *testing.T) {
	convey.Convey("Followings", t, func() {
		rl, err := d.Followings(context.TODO(), 1, 2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(rl, convey.ShouldNotBeNil)
	})
}

func TestStat(t *testing.T) {
	var (
		c         = context.Background()
		mid int64 = 2
	)
	convey.Convey("stat", t, func() {
		reply, err := d.Stat(c, mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(reply, convey.ShouldNotBeNil)
	})
}

func TestStats(t *testing.T) {
	var (
		c    = context.Background()
		mids = []int64{2, 3}
	)
	convey.Convey("stats", t, func() {
		replys, err := d.Stats(c, mids)
		convey.So(err, convey.ShouldBeNil)
		convey.So(replys, convey.ShouldNotBeNil)
	})
}
