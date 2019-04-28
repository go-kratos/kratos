package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohitUpFilter(t *testing.T) {
	convey.Convey("hitUpFilter", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.hitUpFilter(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitUserFilter(t *testing.T) {
	convey.Convey("hitUserFilter", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.hitUserFilter(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitUserFilterCnt(t *testing.T) {
	convey.Convey("hitUserFilterCnt", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.hitUserFilterCnt(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoaddSlashes(t *testing.T) {
	convey.Convey("addSlashes", t, func(ctx convey.C) {
		var (
			str = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := addSlashes(str)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddUserFilter(t *testing.T) {
	convey.Convey("AddUserFilter", t, func(ctx convey.C) {
		var (
			tx, _   = testDao.BeginBiliDMTrans(c)
			mid     = int64(0)
			fType   = int8(0)
			filter  = ""
			comment = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.AddUserFilter(tx, mid, fType, filter, comment)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUserFilter(t *testing.T) {
	convey.Convey("UserFilter", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			fType = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UserFilter(c, mid, fType)
		})
	})
}

func TestDaoUserFilters(t *testing.T) {
	convey.Convey("UserFilters", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := testDao.UserFilters(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUserFiltersByID(t *testing.T) {
	convey.Convey("UserFiltersByID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UserFiltersByID(c, mid, ids)
		})
	})
}

func TestDaoDelUserFilter(t *testing.T) {
	convey.Convey("DelUserFilter", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			ids   = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.DelUserFilter(tx, mid, ids)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoMultiAddUpFilter(t *testing.T) {
	convey.Convey("MultiAddUpFilter", t, func(ctx convey.C) {
		var (
			tx, _  = testDao.BeginBiliDMTrans(c)
			mid    = int64(0)
			fType  = int8(0)
			fltMap = map[string]string{
				"test": "test",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.MultiAddUpFilter(tx, mid, fType, fltMap)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUpFilter(t *testing.T) {
	convey.Convey("UpFilter", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			fType = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpFilter(c, mid, fType)
		})
	})
}

func TestDaoUpFilters(t *testing.T) {
	convey.Convey("UpFilters", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpFilters(c, mid)
		})
	})
}

func TestDaoUpdateUpFilter(t *testing.T) {
	convey.Convey("UpdateUpFilter", t, func(ctx convey.C) {
		var (
			tx, _   = testDao.BeginBiliDMTrans(context.TODO())
			mid     = int64(0)
			fType   = int8(0)
			active  = int8(0)
			filters = []string{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpdateUpFilter(tx, mid, fType, active, filters)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoAddGlobalFilter(t *testing.T) {
	convey.Convey("AddGlobalFilter", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			fType  = int8(0)
			filter = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			lastID, err := testDao.AddGlobalFilter(c, fType, filter)
			ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGlobalFilter(t *testing.T) {
	convey.Convey("GlobalFilter", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			fType  = int8(0)
			filter = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := testDao.GlobalFilter(c, fType, filter)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGlobalFilters(t *testing.T) {
	convey.Convey("GlobalFilters", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(0)
			limit = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.GlobalFilters(c, sid, limit)
		})
	})
}

func TestDaoDelGlobalFilters(t *testing.T) {
	convey.Convey("DelGlobalFilters", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.DelGlobalFilters(c, ids)
		})
	})
}

func TestDaoUserFilterCnt(t *testing.T) {
	convey.Convey("UserFilterCnt", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			tp    = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UserFilterCnt(c, tx, mid, tp)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoInsertUserFilterCnt(t *testing.T) {
	convey.Convey("InsertUserFilterCnt", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			tp    = int8(0)
			count = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.InsertUserFilterCnt(c, tx, mid, tp, count)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUpdateUserFilterCnt(t *testing.T) {
	convey.Convey("UpdateUserFilterCnt", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			tp    = int8(0)
			count = int64(0)
			limit = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpdateUserFilterCnt(c, tx, mid, tp, count, limit)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUpFilterCnt(t *testing.T) {
	convey.Convey("UpFilterCnt", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			tp    = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpFilterCnt(c, tx, mid, tp)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoInsertUpFilterCnt(t *testing.T) {
	convey.Convey("InsertUpFilterCnt", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			tp    = int8(0)
			count = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.InsertUpFilterCnt(c, tx, mid, tp, count)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}

func TestDaoUpdateUpFilterCnt(t *testing.T) {
	convey.Convey("UpdateUpFilterCnt", t, func(ctx convey.C) {
		var (
			tx, _ = testDao.BeginBiliDMTrans(c)
			mid   = int64(0)
			tp    = int8(0)
			count = int(0)
			limit = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpdateUpFilterCnt(c, tx, mid, tp, count, limit)
		})
		ctx.Reset(func() {
			tx.Commit()
		})
	})
}
