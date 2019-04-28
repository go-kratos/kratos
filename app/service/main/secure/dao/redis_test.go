package dao

import (
	"context"
	"testing"
	"time"

	model "go-common/app/service/main/secure/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestRedis(t *testing.T) {
	convey.Convey("TestRedis", t, func() {
		lput := &model.Log{
			Mid:        2,
			Location:   "上海",
			LocationID: 11,
			Time:       xtime.Time(time.Now().Unix()),
		}
		d.AddExpectionMsg(context.TODO(), lput)
		lget, _ := d.ExpectionMsg(context.TODO(), 2)

		convey.So(lget.LocationID, convey.ShouldEqual, lput.LocationID)

		d.AddUnNotify(context.TODO(), 2, "1234")
		b, _ := d.UnNotify(context.TODO(), 2, "1234")
		convey.So(b, convey.ShouldBeTrue)

		d.DelUnNotify(context.TODO(), 2)
		b, _ = d.UnNotify(context.TODO(), 2, "1234")
		convey.So(b, convey.ShouldBeFalse)

		d.AddCount(context.TODO(), 2, "1234")
		d.AddCount(context.TODO(), 2, "1234")
		count, _ := d.Count(context.TODO(), 2, "1234")
		convey.So(count, convey.ShouldEqual, 2)

		d.DelCount(context.TODO(), 2)
		count, _ = d.Count(context.TODO(), 2, "1234")
		convey.So(count, convey.ShouldEqual, 0)

	})
}

func TestDoubleCheckKey(t *testing.T) {
	convey.Convey("doubleCheckKey", t, func() {
		p1 := doubleCheckKey(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestChangePWDKey(t *testing.T) {
	convey.Convey("changePWDKey", t, func() {
		p1 := changePWDKey(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestMsgKey(t *testing.T) {
	convey.Convey("msgKey", t, func() {
		p1 := msgKey(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestUnnotifyKey(t *testing.T) {
	convey.Convey("unnotifyKey", t, func() {
		p1 := unnotifyKey(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestCountKey(t *testing.T) {
	convey.Convey("countKey", t, func() {
		p1 := countKey(0)
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaoAddExpectionMsg(t *testing.T) {
	var mid int64 = 7593623
	convey.Convey("AddExpectionMsg", t, func() {
		l := &model.Log{Mid: mid}
		err := d.AddExpectionMsg(context.TODO(), l)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("ExpectionMsg", t, func() {
		msg, err := d.ExpectionMsg(context.TODO(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(msg, convey.ShouldNotBeNil)
	})
}

func TestDaoAddUnNotify(t *testing.T) {
	convey.Convey("AddUnNotify", t, func() {
		err := d.AddUnNotify(context.TODO(), 0, "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoDelUnNotify(t *testing.T) {
	convey.Convey("DelUnNotify", t, func() {
		err := d.DelUnNotify(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoUnNotify(t *testing.T) {
	convey.Convey("UnNotify", t, func() {
		b, err := d.UnNotify(context.TODO(), 0, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(b, convey.ShouldNotBeNil)
	})
}

func TestDaoCount(t *testing.T) {
	convey.Convey("Count", t, func() {
		count, err := d.Count(context.TODO(), 0, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(count, convey.ShouldNotBeNil)
	})
}

func TestDaoAddCount(t *testing.T) {
	convey.Convey("AddCount", t, func() {
		err := d.AddCount(context.TODO(), 0, "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoAddChangePWDRecord(t *testing.T) {
	convey.Convey("AddChangePWDRecord", t, func() {
		err := d.AddChangePWDRecord(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoChangePWDRecord(t *testing.T) {
	convey.Convey("ChangePWDRecord", t, func() {
		b, err := d.ChangePWDRecord(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(b, convey.ShouldNotBeNil)
	})
}

func TestDaoDelCount(t *testing.T) {
	convey.Convey("DelCount", t, func() {
		err := d.DelCount(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoAddDCheckCache(t *testing.T) {
	convey.Convey("AddDCheckCache", t, func() {
		err := d.AddDCheckCache(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoDCheckCache(t *testing.T) {
	convey.Convey("DCheckCache", t, func() {
		b, err := d.DCheckCache(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(b, convey.ShouldNotBeNil)
	})
}
