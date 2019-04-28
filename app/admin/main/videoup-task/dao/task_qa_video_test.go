package dao

import (
	"context"
	"go-common/app/admin/main/videoup-task/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertQAVideo(t *testing.T) {
	var (
		tx, _ = d.BeginTran(context.TODO())
		dt    = &model.VideoDetail{
			UPGroups: []int64{0},
		}
	)
	convey.Convey("InsertQAVideo", t, func(ctx convey.C) {
		id, err := d.InsertQAVideo(tx, dt)
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

func TestDaoQAVideoDetail(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{437}
	)
	convey.Convey("QAVideoDetail", t, func(ctx convey.C) {
		list, arr, err := d.QAVideoDetail(c, ids)
		ctx.Convey("Then err should be nil.list,arr should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(arr, convey.ShouldNotBeNil)
			ctx.So(list, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetQAVideoID(t *testing.T) {
	var (
		c      = context.TODO()
		aid    = int64(10110610)
		cid    = int64(10134188)
		taskID = int64(8725)
	)
	convey.Convey("GetQAVideoID", t, func(ctx convey.C) {
		id, err := d.GetQAVideoID(c, aid, cid, taskID)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateQAVideoUTime(t *testing.T) {
	var (
		c      = context.TODO()
		aid    = int64(10110610)
		cid    = int64(10134188)
		taskID = int64(8725)
		utime  = int64(10)
	)
	convey.Convey("UpdateQAVideoUTime", t, func(ctx convey.C) {
		err := d.UpdateQAVideoUTime(c, aid, cid, taskID, utime)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelQAVideo(t *testing.T) {
	var (
		c     = context.TODO()
		mtime = time.Now().AddDate(-1, -1, 0)
		limit = int(1)
	)
	convey.Convey("DelQAVideo", t, func(ctx convey.C) {
		_, err := d.DelQAVideo(c, mtime, limit)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelQATask(t *testing.T) {
	var (
		c     = context.TODO()
		mtime = time.Now().AddDate(-1, -1, 0)
		limit = int(1)
	)
	convey.Convey("DelQATask", t, func(ctx convey.C) {
		_, err := d.DelQATask(c, mtime, limit)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoQATaskVideoByID(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("QATaskVideoByID", t, func(ctx convey.C) {
		_, err := d.QATaskVideoByID(c, 0)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoQATaskVideoSimpleByID(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("QATaskVideoSimpleByID", t, func(ctx convey.C) {
		_, err := d.QATaskVideoSimpleByID(c, 0)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
