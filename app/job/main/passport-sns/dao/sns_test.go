package dao

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"go-common/app/job/main/passport-sns/model"
	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestDao_SnsUserByMid(t *testing.T) {
	convey.Convey("SnsUserByMid", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQ
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SnsUserByMid(c, mid, platform)
			ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_SnsUserByUnionID(t *testing.T) {
	convey.Convey("SnsUserByUnionID", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			unionID  = ""
			platform = model.PlatformQQ
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.SnsUserByUnionID(c, unionID, platform)
			ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDao_AddSnsUser(t *testing.T) {
	convey.Convey("AddSnsUser", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQ
			expires  = int64(0)
			unionID  = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.snsDB), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.AddSnsUser(c, mid, expires, unionID, platform)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_UpdateSnsUser(t *testing.T) {
	convey.Convey("UpdateSnsUser", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQ
			expires  = int64(0)
			unionID  = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.snsDB), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.UpdateSnsUser(c, mid, expires, unionID, platform)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_UpdateSnsToken(t *testing.T) {
	convey.Convey("UpdateSnsToken", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.SnsToken{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.snsDB), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.UpdateSnsToken(c, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_TxAddSnsUser(t *testing.T) {
	convey.Convey("TxAddSnsUser", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginSnsTran(context.Background())
			a     = &model.SnsUser{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.TxAddSnsUser(tx, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_TxAddSnsOpenID(t *testing.T) {
	convey.Convey("TxAddSnsOpenID", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginSnsTran(context.Background())
			a     = &model.SnsOpenID{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.TxAddSnsOpenID(tx, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_TxAddSnsToken(t *testing.T) {
	convey.Convey("TxAddSnsToken", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginSnsTran(context.Background())
			a     = &model.SnsToken{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.TxAddSnsToken(tx, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_DelSnsUser(t *testing.T) {
	convey.Convey("DelSnsUser", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			platform = model.PlatformQQ
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.snsDB), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.DelSnsUser(c, mid, platform)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_TxUpdateSnsUserExpires(t *testing.T) {
	convey.Convey("TxUpdateSnsUserExpires", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginSnsTran(context.Background())
			a     = &model.SnsUser{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.TxUpdateSnsUserExpires(tx, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_TxUpdateSnsUser(t *testing.T) {
	convey.Convey("TxUpdateSnsUser", t, func(ctx convey.C) {
		var (
			tx, _    = d.BeginSnsTran(context.Background())
			mid      = int64(0)
			platform = model.PlatformQQ
			expires  = int64(0)
			unionID  = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.TxUpdateSnsUser(tx, mid, expires, unionID, platform)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_TxUpdateSnsToken(t *testing.T) {
	convey.Convey("TxUpdateSnsToken", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginSnsTran(context.Background())
			a     = &model.SnsToken{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(tx), "Exec", func(_ *xsql.Tx, _ string, _ ...interface{}) (res sql.Result, err error) {
				return driver.RowsAffected(1), nil
			})
			defer mock.Unpatch()
			affected, err := d.TxUpdateSnsToken(tx, a)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_openIDSuffix(t *testing.T) {
	convey.Convey("openIDSuffix", t, func(ctx convey.C) {
		res := openIDSuffix("")
		ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
