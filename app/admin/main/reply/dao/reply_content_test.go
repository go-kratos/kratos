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

func deleteReplyContent(d *Dao, oid int64, rpid int64) error {
	_delSQL := "Delete from reply_content_%d where rpid = ?"
	_, err := d.db.Exec(context.Background(), fmt.Sprintf(_delSQL, hit(oid)), rpid)
	if err != nil {
		return err
	}
	return nil
}

func TestReplyContent(t *testing.T) {
	_inContSQL := "INSERT IGNORE INTO reply_content_%d (rpid,message,ats,ip,plat,device,version,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?)"
	d := _d
	c := context.Background()
	now := time.Now()
	rc := model.ReplyContent{
		ID:      1,
		Message: "test",
		Device:  "iphone",
		Version: "beta",
		CTime:   xtime.Time(now.Unix()),
		MTime:   xtime.Time(now.Unix()),
	}
	_, err := d.db.Exec(c, fmt.Sprintf(_inContSQL, hit(6)), rc.ID, rc.Message, rc.Ats, rc.IP, rc.Plat, rc.Device, rc.Version, rc.CTime, rc.MTime)
	if err != nil {
		t.Logf("insert reply content %v", err)
		t.FailNow()
	}
	defer deleteReplyContent(d, 6, 1)
	rc2 := model.ReplyContent{
		ID:      2,
		Message: "test2",
		Device:  "iphone2",
		Version: "beta2",
		CTime:   xtime.Time(now.Unix()),
		MTime:   xtime.Time(now.Unix()),
	}
	_, err = d.db.Exec(c, fmt.Sprintf(_inContSQL, hit(86)), rc2.ID, rc2.Message, rc2.Ats, rc2.IP, rc2.Plat, rc2.Device, rc2.Version, rc2.CTime, rc2.MTime)
	if err != nil {
		t.Logf("insert reply content %v", err)
		t.FailNow()
	}
	defer deleteReplyContent(d, 86, 2)
	Convey("reply_content test", t, WithDao(func(d *Dao) {
		Convey("get update content", WithDao(func(d *Dao) {
			now = time.Now()
			rows, err := d.UpReplyContent(c, 86, 2, "test3", now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
		}))
		Convey("get reply_content", WithDao(func(d *Dao) {
			content, err := d.ReplyContent(c, 86, 2)
			So(err, ShouldBeNil)
			So(content, ShouldNotBeNil)
			So(content.ID, ShouldEqual, 2)
			So(content.CTime, ShouldEqual, now.Unix())
			So(content.Message, ShouldEqual, "test3")
		}))
		Convey("find reply_contents", WithDao(func(d *Dao) {
			cMap, err := d.ReplyContents(c, []int64{6, 86}, []int64{1, 2})
			So(err, ShouldBeNil)
			So(len(cMap), ShouldEqual, 2)
			So(cMap[1].ID, ShouldEqual, 1)
			So(cMap[1].Message, ShouldEqual, "test")
			So(cMap[2].ID, ShouldEqual, 2)
			So(cMap[2].Message, ShouldEqual, "test3")
		}))
	}))
}
