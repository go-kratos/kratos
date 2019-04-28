package dao

import (
	"context"
	xsql "database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"go-common/app/interface/main/feedback/model"
	"go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoJudgeSsnRecord(t *testing.T) {
	convey.Convey("JudgeSsnRecord", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(1)
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			cnt, err := d.JudgeSsnRecord(c, sid)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.selSSnID.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.selSSnID), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.selSSnID.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.JudgeSsnRecord(c, sid)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSession(t *testing.T) {
	convey.Convey("Session", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			buvid   = ""
			system  = ""
			version = ""
			mid     = int64(1)
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			ssn, err := d.Session(c, buvid, system, version, mid)
			ctx.Convey("Then err should be nil.ssn should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(ssn, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSessionCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("SessionCount", t, func(ctx convey.C) {
		cnt, err := d.SessionCount(c, mid)
		ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cnt, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpSsnMtime(t *testing.T) {
	convey.Convey("UpSsnMtime", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			now = time.Now()
			id  = int64(1)
		)
		convey.Convey("When everything is correct", func(ctx convey.C) {
			err := d.UpSsnMtime(c, now, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		convey.Convey("When d.upSsnMtime.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.upSsnMtime), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.upSsnMtime.Exec Error")
				})
			defer guard.Unpatch()
			err := d.UpSsnMtime(c, now, id)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpSsnMtime(t *testing.T) {
	convey.Convey("TxUpSsnMtime", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			now   = time.Now()
			id    = int64(0)
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			err := d.TxUpSsnMtime(tx, now, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec",
				func(_ *sql.Tx, _ string, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("tx.Exec Error")
				})
			defer guard.Unpatch()
			err := d.TxUpSsnMtime(tx, now, id)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestDaoSessionIDByTagID(t *testing.T) {
	convey.Convey("SessionIDByTagID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tagID = []int64{1, 2}
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			sid, err := d.SessionIDByTagID(c, tagID)
			ctx.Convey("Then err should be nil.sid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(sid, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.dbMs.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.SessionIDByTagID(c, tagID)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSessionBySsnID(t *testing.T) {
	convey.Convey("SessionBySsnID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = []int64{1, 2}
			state = "3"
			start = time.Now()
			end   = time.Now()
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			ssns, err := d.SessionBySsnID(c, sid, state, start, end)
			ctx.Convey("Then err should be nil.ssns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(ssns, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.dbMs.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.SessionBySsnID(c, sid, state, start, end)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSSnBySsnIDAllSate(t *testing.T) {
	convey.Convey("SSnBySsnIDAllSate", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = []int64{1, 2}
			start = time.Now()
			end   = time.Now()
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			ssns, err := d.SSnBySsnIDAllSate(c, sid, start, end)
			ctx.Convey("Then err should be nil.ssns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(ssns, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.dbMs.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.SSnBySsnIDAllSate(c, sid, start, end)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSessionByMid(t *testing.T) {
	convey.Convey("SessionByMid", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(1)
			platform = "ios"
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			ssns, err := d.SessionByMid(c, mid, platform)
			ctx.Convey("Then err should be nil.ssns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ssns, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.dbMs.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.SessionByMid(c, mid, platform)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateSessionState(t *testing.T) {
	convey.Convey("TxUpdateSessionState", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			state = int(0)
			sid   = int64(0)
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			err := d.TxUpdateSessionState(tx, state, sid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec",
				func(_ *sql.Tx, _ string, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("tx.Exec Error")
				})
			defer guard.Unpatch()
			err := d.TxUpdateSessionState(tx, state, sid)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestDaoUpdateSessionState(t *testing.T) {
	convey.Convey("UpdateSessionState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			state = int(0)
			sid   = int64(0)
		)
		ctx.Convey("UpdateSessionState", func(ctx convey.C) {
			err := d.UpdateSessionState(c, state, sid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When d.upSsnSta.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.upSsnSta), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.upSsnSta.Exec Error")
				})
			defer guard.Unpatch()
			err := d.UpdateSessionState(c, state, sid)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTags(t *testing.T) {
	convey.Convey("Tags", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mold     = int(1)
			platform = "ios"
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			tMap, err := d.Tags(c, mold, platform)
			ctx.Convey("Then err should be nil.tMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tMap, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.dbMs.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.Tags(c, mold, platform)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddSession(t *testing.T) {
	convey.Convey("AddSession", t, func(ctx convey.C) {
		var (
			c = context.Background()
			s = &model.Session{}
		)
		ctx.Convey("When everthing is correct", func(ctx convey.C) {
			id, err := d.AddSession(c, s)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.inSsn.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.inSsn), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.inSsn.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.AddSession(c, s)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddSession(t *testing.T) {
	convey.Convey("TxAddSession", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			s     = &model.Session{}
		)
		ctx.Convey("When everthing is correct", func(ctx convey.C) {
			id, err := d.TxAddSession(tx, s)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec",
				func(_ *sql.Tx, _ string, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("tx.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.TxAddSession(tx, s)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestDaoAddSessionTag(t *testing.T) {
	convey.Convey("AddSessionTag", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			sessionID = int64(0)
			tagID     = int64(0)
			now       = time.Now()
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			id, err := d.AddSessionTag(c, sessionID, tagID, now)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.inSsnTag.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.inSsnTag), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.inSsnTag.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.AddSessionTag(c, sessionID, tagID, now)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddSessionTag(t *testing.T) {
	convey.Convey("TxAddSessionTag", t, func(ctx convey.C) {
		var (
			tx, _     = d.BeginTran(context.Background())
			sessionID = int64(0)
			tagID     = int64(0)
			now       = time.Now()
		)
		ctx.Convey("When everthing is correct", func(ctx convey.C) {
			id, err := d.TxAddSessionTag(tx, sessionID, tagID, now)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When tx.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec",
				func(_ *sql.Tx, _ string, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("tx.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.TxAddSessionTag(tx, sessionID, tagID, now)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			tx.Rollback()
		})
	})
}

func TestDaoUpdateSession(t *testing.T) {
	convey.Convey("UpdateSession", t, func(ctx convey.C) {
		var (
			c = context.Background()
			s = &model.Session{}
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			affected, err := d.UpdateSession(c, s)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.upSsn.Exec gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.upSsn), "Exec",
				func(_ *sql.Stmt, _ context.Context, _ ...interface{}) (xsql.Result, error) {
					return nil, fmt.Errorf("d.upSsn.Exec Error")
				})
			defer guard.Unpatch()
			_, err := d.UpdateSession(c, s)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagIDBySid(t *testing.T) {
	convey.Convey("TagIDBySid", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			sids = []int64{1, 2}
		)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			tsMap, err := d.TagIDBySid(c, sids)
			ctx.Convey("Then err should be nil.tsMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tsMap, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.dbMs.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.dbMs), "Query",
				func(_ *sql.DB, _ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
					return nil, fmt.Errorf("d.dbMs.Query Error")
				})
			defer guard.Unpatch()
			_, err := d.TagIDBySid(c, sids)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoplatConvert(t *testing.T) {
	convey.Convey("platConvert", t, func(ctx convey.C) {
		var (
			platform = "a,b"
		)
		ctx.Convey("When the original string is 'a,b'", func(ctx convey.C) {
			s := platConvert(platform)
			ctx.Convey("Then s should equal \"a\",\"b\".", func(ctx convey.C) {
				ctx.So(s, convey.ShouldEqual, "\"a\",\"b\"")
			})
		})
	})
}
