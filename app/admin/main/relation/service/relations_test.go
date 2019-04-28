package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/relation/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestFollowers(t *testing.T) {
	convey.Convey("Followers", t, func() {
		rl, err := s.Followers(context.TODO(), &model.FollowersParam{
			Fid: 1,
			Mid: 2,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(rl, convey.ShouldNotBeNil)
	})
}

func TestFollowings(t *testing.T) {
	convey.Convey("Followings", t, func() {
		rl, err := s.Followings(context.TODO(), &model.FollowingsParam{
			Fid: 1,
			Mid: 2,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(rl, convey.ShouldNotBeNil)
	})
}

func TestStat(t *testing.T) {
	convey.Convey("stat", t, func() {
		rl, err := s.Stat(context.TODO(), &model.ArgMid{
			Mid: 2,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(rl, convey.ShouldNotBeNil)
	})
}

func TestStats(t *testing.T) {
	convey.Convey("stats", t, func() {
		rl, err := s.Stats(context.TODO(), &model.ArgMids{
			Mids: []int64{2, 3},
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(rl, convey.ShouldNotBeNil)
	})
}
