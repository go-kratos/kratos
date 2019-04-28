package dao

import (
	"context"
	"go-common/app/service/main/ugcpay/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddCacheOrderUser(t *testing.T) {
	convey.Convey("AddCacheOrderUser", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = "test"
			val = &model.Order{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheOrderUser(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheOrderUser(t *testing.T) {
	convey.Convey("CacheOrderUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheOrderUser(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheOrderUser(t *testing.T) {
	convey.Convey("DelCacheOrderUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheOrderUser(c, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheAsset(t *testing.T) {
	convey.Convey("AddCacheAsset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(233333)
			otype    = "archive"
			currency = "bp"
			value    = &model.Asset{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheAsset(c, id, otype, currency, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheAsset(t *testing.T) {
	convey.Convey("CacheAsset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(233333)
			otype    = "archive"
			currency = "bp"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheAsset(c, id, otype, currency)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCacheAsset(t *testing.T) {
	convey.Convey("DelCacheAsset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(233333)
			otype    = "archive"
			currency = "bp"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheAsset(c, id, otype, currency)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

// func TestDaoCacheAggrIncomeUser(t *testing.T) {
// 	convey.Convey("CacheAggrIncomeUser", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			id  = int64(46333)
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			res, err := d.CacheAggrIncomeUser(c, id, cur)
// 			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 				ctx.So(res, convey.ShouldNotBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoAddCacheAggrIncomeUser(t *testing.T) {
// 	convey.Convey("AddCacheAggrIncomeUser", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			id  = int64(46333)
// 			val = &model.AggrIncomeUser{}
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			err := d.AddCacheAggrIncomeUser(c, id, cur, val)
// 			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoDelCacheAggrIncomeUser(t *testing.T) {
// 	convey.Convey("DelCacheAggrIncomeUser", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			id  = int64(46333)
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			err := d.DelCacheAggrIncomeUser(c, id, cur)
// 			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoCacheAggrIncomeUserMonthly(t *testing.T) {
// 	convey.Convey("CacheAggrIncomeUserMonthly", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			id  = int64(46333)
// 			ver = int64(123)
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			res, err := d.CacheAggrIncomeUserAssetList(c, id, cur, ver)
// 			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 				ctx.So(res, convey.ShouldNotBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoAddCacheAggrIncomeUserMonthly(t *testing.T) {
// 	convey.Convey("AddCacheAggrIncomeUserMonthly", t, func(ctx convey.C) {
// 		var (
// 			c     = context.Background()
// 			id    = int64(46333)
// 			ver   = int64(123)
// 			cur   = "bp"
// 			value = []*model.AggrIncomeUserAsset{}
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			err := d.AddCacheAggrIncomeUserAssetList(c, id, cur, ver, value)
// 			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoDelCacheAggrIncomeUserMonthly(t *testing.T) {
// 	convey.Convey("DelCacheAggrIncomeUserMonthly", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			id  = int64(46333)
// 			ver = int64(123)
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			err := d.DelCacheAggrIncomeUserAssetList(c, id, cur, ver)
// 			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 			})
// 		})
// 	})
// }
