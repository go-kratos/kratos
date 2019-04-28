package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBatchQueryAccount(t *testing.T) {
	convey.Convey("BatchQueryAccount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(0)
			limit = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchQueryAccount(c, start, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryAccountByMid(t *testing.T) {
	convey.Convey("QueryAccountByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryAccountByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryAccountByTel(t *testing.T) {
	convey.Convey("QueryAccountByTel", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tel = "13122111111"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryAccountByTel(c, tel)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryAccountByMail(t *testing.T) {
	convey.Convey("QueryAccountByMail", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mail = "598717394@qq.com"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryAccountByMail(c, mail)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchQueryAccountByTime(t *testing.T) {
	convey.Convey("BatchQueryAccountByTime", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = time.Now()
			end   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchQueryAccountByTime(c, start, end)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchQueryAccountInfo(t *testing.T) {
	convey.Convey("BatchQueryAccountInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			start  = int64(0)
			limit  = int64(0)
			suffix = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchQueryAccountInfo(c, start, limit, suffix)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryAccountInfoByMid(t *testing.T) {
	convey.Convey("QueryAccountInfoByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryAccountInfoByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchQueryAccountInfoByTime(t *testing.T) {
	convey.Convey("BatchQueryAccountInfoByTime", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			start  = time.Now()
			end    = time.Now()
			suffix = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchQueryAccountInfoByTime(c, start, end, suffix)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchQueryAccountSns(t *testing.T) {
	convey.Convey("BatchQueryAccountSns", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(0)
			limit = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchQueryAccountSns(c, start, limit)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryAccountSnsByMid(t *testing.T) {
	convey.Convey("QueryAccountSnsByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryAccountSnsByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryTelBindLog(t *testing.T) {
	convey.Convey("QueryTelBindLog", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryTelBindLog(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryAccountRegByMid(t *testing.T) {
	convey.Convey("QueryAccountRegByMid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(22222222)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.QueryAccountRegByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchQueryAccountRegByTime(t *testing.T) {
	convey.Convey("BatchQueryAccountRegByTime", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			start  = time.Now()
			end    = time.Now()
			suffix = int(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.BatchQueryAccountRegByTime(c, start, end, suffix)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}
