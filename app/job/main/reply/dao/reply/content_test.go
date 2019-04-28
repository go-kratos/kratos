package reply

import (
	"context"
	"go-common/app/job/main/reply/model/reply"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyTxInsert1(t *testing.T) {
	convey.Convey("TxInsert", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			rc    = &reply.Content{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Content.TxInsert(tx, oid, rc)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyGet1(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rc, err := d.Content.Get(context.Background(), oid, rpID)
			ctx.Convey("Then err should be nil.rc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rc, convey.ShouldBeNil)
			})
		})
	})
}
