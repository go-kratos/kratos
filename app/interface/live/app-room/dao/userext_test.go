package dao

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/app-room/model"
	"math/rand"
	"testing"
	"time"
)

func getTestRandUid() int64 {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(10000000)
}

func TestDao_UserConf(t *testing.T) {
	Convey("normal", t, func() {
		Convey("conf", func() {
			So(7, ShouldEqual, testDao.c.Gift.RechargeTip.SilverTipDays[0])
			So(1, ShouldEqual, len(testDao.c.Gift.RechargeTip.SilverTipDays))
			t.Logf("%v", testDao.c.Gift.RechargeTip.SilverTipDays)
		})
		Convey("set and get", func() {
			uid := getTestRandUid()
			err := testDao.SetUserConf(context.Background(), uid, model.GoldTarget, 1)
			So(err, ShouldBeNil)
			m, err := testDao.GetUserConf(context.Background(), uid, model.GoldTarget, []int64{1})
			So(err, ShouldBeNil)
			v, ok := m[1]
			t.Logf("v: %v", v)
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, "1")
			So(m.IsSet(1), ShouldBeTrue)

			testDao.DelUserConf(context.Background(), uid, model.GoldTarget, 1)

			m, err = testDao.GetUserConf(context.Background(), uid, model.GoldTarget, []int64{1})
			So(err, ShouldBeNil)
			v, ok = m[1]
			t.Logf("v: %v", v)
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, "")
			So(m.IsSet(1), ShouldBeFalse)

		})
	})
}
