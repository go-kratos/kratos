package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	mid     int64 = 1
	nowTs         = time.Now().Unix()
	lastID  int64
	lastID2 int64
)

func Test_AddNotice(t *testing.T) {
	c := context.Background()
	nt := model.Notice{
		Plat:       model.PlatAndroid,
		Condition:  model.ConditionGT,
		Build:      1113,
		Title:      "测试",
		Status:     model.StatusOffline,
		Content:    "测试内容",
		Link:       "http://www.bilibili.com",
		StartTime:  xtime.Time(nowTs),
		EndTime:    xtime.Time(nowTs + 10*3600),
		ClientType: "android",
	}
	nt2 := model.Notice{
		Plat:       model.PlatAndroid,
		Condition:  model.ConditionLT,
		Build:      2233,
		Title:      "测试2",
		Status:     model.StatusOffline,
		Content:    "测试内容2",
		Link:       "http://www.bilibili.com",
		StartTime:  xtime.Time(nowTs),
		EndTime:    xtime.Time(nowTs + 10*3600),
		ClientType: "",
	}
	Convey("add notice", t, WithDao(func(d *Dao) {
		var err error
		lastID, err = d.CreateNotice(c, &nt)
		So(err, ShouldBeNil)
		So(lastID, ShouldBeGreaterThan, 0)

		lastID2, err = d.CreateNotice(c, &nt2)
		So(err, ShouldBeNil)
		So(lastID2, ShouldBeGreaterThan, 0)

	}))
}

func Test_ListNotice(t *testing.T) {
	c := context.Background()
	Convey("list notice", t, WithDao(func(d *Dao) {

		nts, err := d.ListNotice(c, 1, 100)
		So(err, ShouldBeNil)
		So(len(nts), ShouldBeGreaterThan, 0)
		So(nts[0].StartTime.Time().Unix(), ShouldBeGreaterThanOrEqualTo, nowTs)

		count, err := d.CountNotice(c)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 1)

	}))
}

func Test_UpdateNotice(t *testing.T) {
	c := context.Background()
	Convey("update notice", t, WithDao(func(d *Dao) {
		data, err := d.Notice(c, uint32(lastID2))
		So(err, ShouldBeNil)
		So(data.Title, ShouldEqual, "测试2")

		data.ID = uint32(lastID2)
		data.Title = "测试3"
		rows, err := d.UpdateNotice(c, data)
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)
		nt2, err := d.Notice(c, uint32(lastID2))
		So(err, ShouldBeNil)
		So(nt2.Title, ShouldEqual, "测试3")

		rows, err = d.UpdateNoticeStatus(c, model.StatusOnline, uint32(lastID))
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)

		nts, err := d.RangeNotice(c, model.PlatAndroid, xtime.Time(nowTs)-3600, xtime.Time(nowTs+5*3600))
		So(err, ShouldBeNil)
		So(len(nts), ShouldBeGreaterThan, 0)
		var isFound bool
		for _, data = range nts {
			if data.ID == uint32(lastID) {
				isFound = true
			}
		}
		So(isFound, ShouldBeTrue)

	}))
}

func Test_DeleteNotice(t *testing.T) {
	c := context.Background()
	Convey("delete notice", t, WithDao(func(d *Dao) {
		rows, err := d.DeleteNotice(c, uint32(lastID))
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)
		rows, err = d.DeleteNotice(c, uint32(lastID2))
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)
	}))
}
