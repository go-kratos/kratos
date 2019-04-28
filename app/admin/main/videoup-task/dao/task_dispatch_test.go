package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpGtimeByID(t *testing.T) {
	convey.Convey("UpGtimeByID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			gtime = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpGtimeByID(c, id, gtime)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserUndoneSpecTask(t *testing.T) {
	convey.Convey("UserUndoneSpecTask", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			uid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UserUndoneSpecTask(c, uid)
			ctx.Convey("Then err should be nil.tasks should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoListByCondition(t *testing.T) {
	convey.Convey("ListByCondition", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			uid    = int64(0)
			pn     = int(0)
			ps     = int(0)
			ltype  = int8(0)
			leader = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tasks, err := d.ListByCondition(c, uid, pn, ps, ltype, leader)
			ctx.Convey("Then err should be nil.tasks should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tasks, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaosqlHelper(t *testing.T) {
	convey.Convey("sqlHelper", t, func(ctx convey.C) {
		var (
			uid    = int64(0)
			pn     = int(0)
			ps     = int(0)
			ltype  = int8(0)
			leader = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.sqlHelper(uid, pn, ps, ltype, leader)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetWeightDB(t *testing.T) {
	convey.Convey("GetWeightDB", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetWeightDB(c, []int64{0, 1, 2})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTaskByID(t *testing.T) {
	convey.Convey("TaskByID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TaskByID(c, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpTaskByID(t *testing.T) {
	convey.Convey("TxUpTaskByID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		tx, _ := d.BeginTran(c)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TxUpTaskByID(tx, 0, map[string]interface{}{"uid": 0})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxReleaseByID(t *testing.T) {
	convey.Convey("TxReleaseByID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		tx, _ := d.BeginTran(c)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TxReleaseByID(tx, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoMulReleaseMtime(t *testing.T) {
	convey.Convey("MulReleaseMtime", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.MulReleaseMtime(c, []int64{1}, time.Now())
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetTimeOutTask(t *testing.T) {
	convey.Convey("GetTimeOutTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetTimeOutTask(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetRelTask(t *testing.T) {
	convey.Convey("GetRelTask", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.GetRelTask(c, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
