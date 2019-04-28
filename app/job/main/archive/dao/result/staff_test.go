package result

import (
	"context"
	"testing"

	"go-common/app/job/main/archive/model/archive"

	"github.com/smartystreets/goconvey/convey"
)

func TestTxDelStaff(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(4052032)
	)
	convey.Convey("TxDelStaff", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxDelStaff(c, tx, aid)
			ctx.So(err, convey.ShouldBeNil)
			err = tx.Commit()
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestTxAddStaff(t *testing.T) {
	var (
		c     = context.TODO()
		aid   = int64(4052032)
		staff []*archive.Staff
	)
	convey.Convey("TxAddStaff", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			staff = append(staff, &archive.Staff{Aid: aid, Mid: 3333, Title: "哈哈", Ctime: "2018-11-28T16:50:14+08:00", Mtime: "2018-12-21T11:41:37+08:00"})
			staff = append(staff, &archive.Staff{Aid: aid, Mid: 4444, Title: "2223", Ctime: "2018-11-28T16:50:14+08:00", Mtime: "2018-12-21T11:41:38+08:00"})
			tx, err := d.BeginTran(c)
			ctx.So(err, convey.ShouldBeNil)
			err = d.TxAddStaff(c, tx, aid, staff)
			ctx.So(err, convey.ShouldBeNil)
			err = tx.Commit()
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
