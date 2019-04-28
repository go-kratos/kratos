package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInTaskQA(t *testing.T) {
	var (
		tx, _    = d.BeginTran(context.TODO())
		uid      = int64(421)
		detailID = int64(1)
		taskType = int8(1)
	)
	convey.Convey("InTaskQA", t, func(ctx convey.C) {
		id, err := d.InTaskQA(tx, uid, detailID, taskType)
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpTask(t *testing.T) {
	var (
		c     = context.TODO()
		id    = int64(41)
		state = int16(2)
		ftime = time.Now()
	)
	convey.Convey("UpTask", t, func(ctx convey.C) {
		_, err := d.UpTask(c, id, state, ftime)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
