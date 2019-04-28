package dao

import (
	"context"
	"go-common/app/admin/main/reply/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyIdx(t *testing.T) {
	convey.Convey("keyIdx", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			tp   = int32(0)
			sort = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyIdx(oid, tp, sort)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyNewRootIdx(t *testing.T) {
	convey.Convey("keyNewRootIdx", t, func(ctx convey.C) {
		var (
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyNewRootIdx(rpID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyAuditIdx(t *testing.T) {
	convey.Convey("keyAuditIdx", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := keyAuditIdx(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireIndex(t *testing.T) {
	convey.Convey("ExpireIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			typ  = int32(0)
			sort = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := _d.ExpireIndex(c, oid, typ, sort)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireNewChildIndex(t *testing.T) {
	convey.Convey("ExpireNewChildIndex", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			root = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ok, err := _d.ExpireNewChildIndex(c, root)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountReplies(t *testing.T) {
	convey.Convey("CountReplies", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int32(0)
			sort = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := _d.CountReplies(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMinScore(t *testing.T) {
	convey.Convey("MinScore", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int32(0)
			sort = int32(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			score, err := _d.MinScore(c, oid, tp, sort)
			ctx.Convey("Then err should be nil.score should not be nil.", func(ctx convey.C) {
				if err != nil {
					ctx.So(err, convey.ShouldNotBeNil)
				}
				ctx.So(score, convey.ShouldEqual, 0)
			})
		})
	})
}

func TestDaoAddFloorIndex(t *testing.T) {
	convey.Convey("AddFloorIndex", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.AddFloorIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCountIndex(t *testing.T) {
	convey.Convey("AddCountIndex", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.AddCountIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddLikeIndex(t *testing.T) {
	convey.Convey("AddLikeIndex", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rp  = &model.Reply{}
			rpt = &model.Report{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.AddLikeIndex(c, rp, rpt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelIndexBySort(t *testing.T) {
	convey.Convey("DelIndexBySort", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rp   = &model.Reply{}
			sort = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.DelIndexBySort(c, rp, sort)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelReplyIndex(t *testing.T) {
	convey.Convey("DelReplyIndex", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.DelReplyIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddNewChildIndex(t *testing.T) {
	convey.Convey("AddNewChildIndex", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.AddNewChildIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelAuditIndex(t *testing.T) {
	convey.Convey("DelAuditIndex", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			rp = &model.Reply{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.DelAuditIndex(c, rp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
