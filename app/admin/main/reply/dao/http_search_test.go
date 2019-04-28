package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateReply(t *testing.T) {
	Convey("update search reply", t, WithDao(func(d *Dao) {
		err := d.UpSearchReply(context.Background(), map[int64]*model.Reply{111852176: &model.Reply{
			ID:    111852176,
			Oid:   10098544,
			Type:  1,
			CTime: xtime.Time(1534412702),
			MTime: xtime.Time(time.Now().Unix()),
			State: 0,
		}}, 3)
		So(err, ShouldBeNil)
		fmt.Println(err)
	}))
}

func TestSearchAdminLog(t *testing.T) {
	Convey("test search adminlog", t, WithDao(func(d *Dao) {
		res, err := d.SearchAdminLog(context.TODO(), []int64{111843721})
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		fmt.Printf("%+v", res[0])
	}))
}

func TestSearchMonitor(t *testing.T) {
	var (
		c  = context.Background()
		sp = &model.SearchMonitorParams{
			Mode: 0,
			Type: 1,
			Oid:  10099866,
			Sort: "unverify_num",
		}
		oid    int64 = 10099866
		typ    int32 = 1
		remark       = "remark"
	)
	Convey("test search monitor", t, WithDao(func(d *Dao) {
		res, err := d.SearchMonitor(c, sp, 1, 20)
		So(err, ShouldBeNil)
		So(len(res.Result), ShouldBeGreaterThan, 0)
		So(res.Result[0].Oid, ShouldEqual, sp.Oid)
	}))
	Convey("test add monitor", t, WithDao(func(d *Dao) {
		sub, _ := d.Subject(c, oid, typ)
		sub.AttrSet(model.AttrYes, model.SubAttrMonitor)
		err := d.UpSearchMonitor(c, sub, remark)
		So(err, ShouldBeNil)
	}))
	time.Sleep(5 * time.Second)
	Convey("test search monitor", t, WithDao(func(d *Dao) {
		res, err := d.SearchMonitor(c, sp, 1, 1)
		So(err, ShouldBeNil)
		So(len(res.Result), ShouldBeGreaterThan, 0)
		So(res.Result[0].Remark, ShouldEqual, remark)
	}))
}
