package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyUserFilter(t *testing.T) {
	convey.Convey("keyUserFilter", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyUserFilter(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyUpFilter(t *testing.T) {
	convey.Convey("keyUpFilter", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyUpFilter(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyGlobalFilter(t *testing.T) {
	convey.Convey("keyGlobalFilter", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyGlobalFilter()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddUserFilterCache(t *testing.T) {
	convey.Convey("AddUserFilterCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			data = []*model.UserFilter{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddUserFilterCache(c, mid, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelUserFilterCache(t *testing.T) {
	convey.Convey("DelUserFilterCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelUserFilterCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUserFilterCache(t *testing.T) {
	convey.Convey("UserFilterCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UserFilterCache(c, mid)
		})
	})
}

func TestDaoAddUpFilterCache(t *testing.T) {
	convey.Convey("AddUpFilterCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			data = []*model.UpFilter{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddUpFilterCache(c, mid, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelUpFilterCache(t *testing.T) {
	convey.Convey("DelUpFilterCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelUpFilterCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpFilterCache(t *testing.T) {
	convey.Convey("UpFilterCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := testDao.UpFilterCache(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddGlobalFilterCache(t *testing.T) {
	convey.Convey("AddGlobalFilterCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = []*model.GlobalFilter{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.AddGlobalFilterCache(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelGlobalFilterCache(t *testing.T) {
	convey.Convey("DelGlobalFilterCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelGlobalFilterCache(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGlobalFilterCache(t *testing.T) {
	convey.Convey("GlobalFilterCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.GlobalFilterCache(c)
		})
	})
}
