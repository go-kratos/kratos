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

func deleteReply(d *Dao, oid int64, rpid int64) error {
	_delSQL := "Delete from reply_%d where id = ?"
	_, err := d.db.Exec(context.Background(), fmt.Sprintf(_delSQL, hit(oid)), rpid)
	if err != nil {
		return err
	}
	return nil
}

func TestDataReply(t *testing.T) {
	d := _d
	c := context.Background()
	now := time.Now()
	r := model.Reply{
		ID:    2,
		Oid:   86,
		Type:  1,
		Mid:   2233,
		Floor: 1,
		CTime: xtime.Time(now.Unix()),
		MTime: xtime.Time(now.Unix()),
	}
	d.InsertReply(c, &r)
	defer deleteReply(d, r.Oid, r.ID)
	r2 := model.Reply{
		ID:    1,
		Oid:   6,
		Type:  1,
		Mid:   2233,
		Floor: 1,
		CTime: xtime.Time(now.Unix()),
		MTime: xtime.Time(now.Unix()),
	}
	d.InsertReply(c, &r2)
	defer deleteReply(d, r2.Oid, r2.ID)

	Convey("reply test", t, WithDao(func(d *Dao) {
		Convey("get reply", WithDao(func(d *Dao) {
			rp, err := d.Reply(c, 86, r.ID)
			So(err, ShouldBeNil)
			So(rp, ShouldNotBeNil)
			So(rp.ID, ShouldEqual, r.ID)
			So(rp.Mid, ShouldEqual, 2233)
		}))
		Convey("increase reply RCount", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			now = time.Now()
			rows, err := d.TxIncrReplyRCount(t, 86, r.ID, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			rp, err := d.Reply(c, 86, r.ID)
			So(err, ShouldBeNil)
			So(rp, ShouldNotBeNil)
			So(rp.RCount, ShouldEqual, 1)
			So(rp.MTime, ShouldEqual, now.Unix())
			Convey("decrease reply RCount", WithDao(func(d *Dao) {
				t, err := d.BeginTran(c)
				So(err, ShouldBeNil)
				now = time.Now()
				rows, err := d.TxDecrReplyRCount(t, 86, r.ID, now)
				So(err, ShouldBeNil)
				So(rows, ShouldBeGreaterThan, 0)
				err = t.Commit()
				So(err, ShouldBeNil)
				rp, err := d.Reply(c, 86, r.ID)
				So(err, ShouldBeNil)
				So(rp, ShouldNotBeNil)
				So(rp.RCount, ShouldEqual, 0)
				So(rp.MTime, ShouldEqual, now.Unix())
			}))
		}))
		Convey("set reply attr", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			now = time.Now()
			rows, err := d.TxUpReplyAttr(t, 86, r.ID, 3, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			rp, err := d.Reply(c, 86, r.ID)
			So(err, ShouldBeNil)
			So(rp, ShouldNotBeNil)
			So(rp.Attr, ShouldEqual, 3)
			So(rp.MTime, ShouldEqual, now.Unix())
		}))
		Convey("set reply state", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			now = time.Now()
			rows, err := d.TxUpdateReplyState(t, 86, r.ID, 2, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			rp, err := d.Reply(c, 86, r.ID)
			So(err, ShouldBeNil)
			So(rp, ShouldNotBeNil)
			So(rp.State, ShouldEqual, 2)
			So(rp.MTime, ShouldEqual, now.Unix())
		}))
		Convey("get replies", WithDao(func(d *Dao) {
			rpMap, err := d.Replies(c, []int64{6, 86}, []int64{1, 2})
			So(err, ShouldBeNil)
			So(len(rpMap), ShouldEqual, 2)
			So(rpMap[1].ID, ShouldEqual, 1)
			So(rpMap[1].Oid, ShouldEqual, 6)
			So(rpMap[1].Mid, ShouldEqual, 2233)
			So(rpMap[2].ID, ShouldEqual, 2)
			So(rpMap[2].Oid, ShouldEqual, 86)
			So(rpMap[2].Mid, ShouldEqual, 2233)
		}))
		Convey("export replies test", WithDao(func(d *Dao) {
			var (
				ctime time.Time
				etime time.Time
				err   error
			)
			ctime, err = time.Parse("2006-01-02", "2008-03-03")
			So(err, ShouldBeNil)
			etime, err = time.Parse("2006-01-02", "2018-03-03")
			So(err, ShouldBeNil)
			_, err = d.ExportReplies(c, 3400, 0, 3, "", ctime, etime)
			So(err, ShouldBeNil)
		}))
	}))
}
