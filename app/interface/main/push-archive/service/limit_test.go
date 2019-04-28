package service

import (
	"testing"
	"time"

	"go-common/app/interface/main/push-archive/dao"

	"github.com/smartystreets/goconvey/convey"
)

func Test_limit(t *testing.T) {
	initd2()

	upper := int64(113)
	convey.Convey("upper主次数限制", t, func() {
		convey.So(s.limit(upper), convey.ShouldEqual, false)
		convey.So(s.limit(upper), convey.ShouldEqual, true)
		convey.So(s.limit(upper), convey.ShouldEqual, true)
		time.Sleep(time.Second * 2)
		convey.So(s.limit(upper), convey.ShouldEqual, false)
		convey.So(s.limit(upper), convey.ShouldEqual, true)
	})

	fan := int64(121)
	g := &dao.FanGroup{
		Limit:         2,
		PerUpperLimit: 0,
		LimitExpire:   2,
	}
	noLimitFans := &map[int64]int{}
	convey.Convey("粉丝推送次数限制", t, func() {
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, true)
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, true)
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, false)
		time.Sleep(time.Second * 2)
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, true)
	})

	g.PerUpperLimit = 1
	fan = int64(333)
	convey.Convey("粉丝推送upper主的次数限制", t, func() {
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, true)
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, false)
		time.Sleep(time.Second * 2)
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, true)
	})

	convey.Convey("粉丝不受pushlimit限制", t, func() {
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, false)
		(*noLimitFans)[fan] = 1
		convey.So(s.pushLimit(fan, upper, g, noLimitFans), convey.ShouldEqual, true)
	})
}
