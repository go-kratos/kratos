package dao

import (
	"context"
	"go-common/app/admin/main/reply/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTxCountFoldedReplies(t *testing.T) {
	convey.Convey("TxCountFoldedReplies", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			oid   = int64(0)
			tp    = int32(0)
			root  = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.TxCountFoldedReplies(tx, oid, tp, root)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoFoldedReplies(t *testing.T) {
	convey.Convey("FoldedReplies", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int32(0)
			root = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rps, err := d.FoldedReplies(c, oid, tp, root)
			ctx.Convey("Then err should be nil.rps should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rps, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemRdsByFold(t *testing.T) {
	convey.Convey("RemRdsByFold", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			roots    = []int64{}
			childMap map[int64][]int64
			sub      = &model.Subject{}
			rpMap    map[int64]*model.Reply
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.RemRdsByFold(c, roots, childMap, sub, rpMap)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoAddRdsByFold(t *testing.T) {
	convey.Convey("AddRdsByFold", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			roots    = []int64{}
			childMap map[int64][]int64
			sub      = &model.Subject{}
			rpMap    map[int64]*model.Reply
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.AddRdsByFold(c, roots, childMap, sub, rpMap)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
