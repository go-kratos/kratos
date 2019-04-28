package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdminLog(t *testing.T) {
	var (
		oids    = []int64{1, 2, 3}
		rpIDs   = []int64{10, 20, 30}
		adminID = int64(100)
		typ     = int32(1)
		c       = context.Background()
		now     = time.Now()
	)
	Convey("add admin log", t, WithDao(func(d *Dao) {
		rows, err := d.AddAdminLog(c, oids, rpIDs, adminID, typ, model.AdminIsNew, model.AdminIsReport, model.AdminOperDelete, "result", "remark", now)
		So(err, ShouldBeNil)
		So(rows, ShouldNotEqual, 0)
		t.Log(rows)
		rows, err = d.UpAdminNotNew(c, rpIDs, now)
		So(err, ShouldBeNil)
		So(rows, ShouldNotEqual, 0)
		res, err := d.AdminLogsByRpID(c, rpIDs[0])
		So(err, ShouldBeNil)
		So(len(res), ShouldNotEqual, 0)
	}))
}
