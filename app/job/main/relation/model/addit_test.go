package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_AllAchieveFromFollower(t *testing.T) {
	Convey("AllAchieveFromFollower", t, func() {
		flags := AllAchieveFromFollower(500)
		So(flags, ShouldBeEmpty)

		flags = AllAchieveFromFollower(1000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k})

		flags = AllAchieveFromFollower(2000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k})

		flags = AllAchieveFromFollower(5000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k})

		flags = AllAchieveFromFollower(5001)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k})

		flags = AllAchieveFromFollower(10000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k, FollowerAchieve10k})

		flags = AllAchieveFromFollower(100000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k, FollowerAchieve10k, FollowerAchieve10k << 1})

		flags = AllAchieveFromFollower(200000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k, FollowerAchieve10k, FollowerAchieve10k << 1, FollowerAchieve10k << 2})

		flags = AllAchieveFromFollower(300000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k, FollowerAchieve10k, FollowerAchieve10k << 1, FollowerAchieve10k << 2, FollowerAchieve10k << 3})

		flags = AllAchieveFromFollower(305000)
		So(flags, ShouldResemble, []AchieveFlag{FollowerAchieve1k, FollowerAchieve5k, FollowerAchieve10k, FollowerAchieve10k << 1, FollowerAchieve10k << 2, FollowerAchieve10k << 3})
	})
}

func TestDao_AchieveFromFollower(t *testing.T) {
	Convey("AchieveFromFollower", t, func() {
		flag := AchieveFromFollower(500)
		So(flag, ShouldBeZeroValue)

		flag = AchieveFromFollower(1000)
		So(flag, ShouldEqual, FollowerAchieve1k)

		flag = AchieveFromFollower(2000)
		So(flag, ShouldEqual, FollowerAchieve1k)

		flag = AchieveFromFollower(5000)
		So(flag, ShouldEqual, FollowerAchieve5k)

		flag = AchieveFromFollower(10000)
		So(flag, ShouldEqual, FollowerAchieve10k)

		flag = AchieveFromFollower(100000)
		So(flag, ShouldEqual, FollowerAchieve10k<<1)

		flag = AchieveFromFollower(200000)
		So(flag, ShouldEqual, FollowerAchieve10k<<2)

		flag = AchieveFromFollower(305000)
		So(flag, ShouldEqual, FollowerAchieve10k<<3)
	})
}
