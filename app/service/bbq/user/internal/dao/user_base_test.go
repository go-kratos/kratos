package dao

import (
	"context"
	"go-common/app/service/bbq/user/api"
	"go-common/library/log"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyUserBase(t *testing.T) {
	convey.Convey("keyUserBase", t, func(ctx convey.C) {
		var (
			mid = int64(88895104)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyUserBase(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTmpUserBase(t *testing.T) {
	convey.Convey("JustGetUserBase", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{88895104}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UserBase(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawUserBase(t *testing.T) {
	convey.Convey("RawUserBase", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{88895104}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawUserBase(c, mids)
			log.Infow(c, "log", "xxxxxxxxxxx", "res", res)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCacheUserBase(t *testing.T) {
	convey.Convey("CacheUserBase", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheUserBase(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheUserBase(t *testing.T) {
	convey.Convey("AddCacheUserBase", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			userBases map[int64]*api.UserBase
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheUserBase(c, userBases)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheUserBase(t *testing.T) {
	convey.Convey("DelCacheUserBase", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.DelCacheUserBase(c, mid)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

//
//func TestDaoTxAddUserBase(t *testing.T) {
//	convey.Convey("TxAddUserBase", t, func(ctx convey.C) {
//		var (
//			c        = context.Background()
//			tx       = &sql.Tx{}
//			userBase = &api.UserBase{}
//		)
//		ctx.Convey("When everything goes positive", func(ctx convey.C) {
//			num, err := d.TxAddUserBase(c, tx, userBase)
//			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(num, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestDaoAddUserBaseUname(t *testing.T) {
//	convey.Convey("AddUserBaseUname", t, func(ctx convey.C) {
//		var (
//			c     = context.Background()
//			mid   = int64(88895104)
//			uname = "dfsfwe"
//		)
//		ctx.Convey("When everything goes positive", func(ctx convey.C) {
//			num, err := d.AddUserBaseUname(c, mid, uname)
//			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(num, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

//
//func TestDaoAddUserBase(t *testing.T) {
//	convey.Convey("AddUserBase", t, func(ctx convey.C) {
//		var (
//			c        = context.Background()
//			userBase = &api.UserBase{}
//		)
//		ctx.Convey("When everything goes positive", func(ctx convey.C) {
//			num, err := d.AddUserBase(c, userBase)
//			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(num, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestDaoUpdateUserBaseUname(t *testing.T) {
	convey.Convey("UpdateUserBaseUname", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			uname = "sdfewxc"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.UpdateUserBaseUname(c, mid, uname)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserBase(t *testing.T) {
	convey.Convey("UpdateUserBase", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(88895104)
			userBase = &api.UserBase{Uname: "sdfweo", Sex: 1, Face: "htttttt"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.UpdateUserBase(c, mid, userBase)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateUserField(t *testing.T) {
	convey.Convey("UpdateUserField", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.UpdateUserField(c, tx, mid, "uname", "unnnnnn")
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
	convey.Convey("UpdateUserField", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.UpdateUserField(c, tx, mid, "face", "faaaaaa")
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
	convey.Convey("UpdateUserField", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.UpdateUserField(c, tx, mid, "region", 11111)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
	convey.Convey("update cms_tag", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			//tx.Rollback()
			tx.Commit()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.UpdateUserField(c, tx, mid, "cms_tag", 1)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCheckUname(t *testing.T) {
	convey.Convey("CheckUname", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(88895104)
			uname = "lkwejroiw"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.CheckUname(c, mid, uname)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
