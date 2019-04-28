package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/ugcpay-rank/internal/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetCacheElecUPRank(t *testing.T) {
	convey.Convey("SetCacheElecUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
			val = &model.RankElecUPProto{
				CountUPTotalElec: 233,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetCacheElecUPRank(c, id, ver, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheElecUPRank(t *testing.T) {
	convey.Convey("CacheElecUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheElecUPRank(c, id, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)

			res, err = d.CacheElecUPRank(c, 0, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCacheElecUPRank(t *testing.T) {
	convey.Convey("DelCacheElecUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheElecUPRank(c, id, ver)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetCacheElecAVRank(t *testing.T) {
	convey.Convey("SetCacheElecAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			val = &model.RankElecAVProto{
				CountUPTotalElec: 233,
			}
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetCacheElecAVRank(c, id, ver, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheElecAVRank(t *testing.T) {
	convey.Convey("CacheElecAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheElecAVRank(c, id, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)

			res, err = d.CacheElecAVRank(c, 0, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCacheElecAVRank(t *testing.T) {
	convey.Convey("DelCacheElecAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheElecAVRank(c, id, ver)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetCacheElecPrepUPRank(t *testing.T) {
	convey.Convey("SetCacheElecPrepUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			val = &model.RankElecPrepUPProto{
				Count: 233,
			}
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetCacheElecPrepUPRank(c, id, ver, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheElecPrepUPRank(t *testing.T) {
	convey.Convey("CacheElecPrepUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, item, err := d.CacheElecPrepUPRank(c, id, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(item, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)

			res, item, err = d.CacheElecPrepUPRank(c, 0, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(item, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCacheElecPrepUPRank(t *testing.T) {
	convey.Convey("DelCacheElecPrepUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheElecPrepUPRank(c, id, ver)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetCacheElecPrepAVRank(t *testing.T) {
	convey.Convey("SetCacheElecPrepAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			val = &model.RankElecPrepAVProto{
				AVID: 233,
			}
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetCacheElecPrepAVRank(c, id, ver, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheElecPrepAVRank(t *testing.T) {
	convey.Convey("CacheElecPrepAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, item, err := d.CacheElecPrepAVRank(c, id, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(item, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)

			res, item, err = d.CacheElecPrepAVRank(c, 0, ver)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(item, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCacheElecPrepAVRank(t *testing.T) {
	convey.Convey("DelCacheElecPrepAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheElecPrepAVRank(c, id, ver)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheElecPrepAVRank(t *testing.T) {
	convey.Convey("AddCacheElecPrepAVRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			val = &model.RankElecPrepAVProto{
				AVID: 233,
			}
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.AddCacheElecPrepAVRank(c, id, ver, val)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheElecPrepUPRank(t *testing.T) {
	convey.Convey("AddCacheElecPrepUPRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(233)
			val = &model.RankElecPrepUPProto{
				UPMID: 233,
			}
			ver = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := d.AddCacheElecPrepUPRank(c, id, ver, val)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}
