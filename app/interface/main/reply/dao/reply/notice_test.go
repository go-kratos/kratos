package reply

import (
	"context"
	"go-common/library/database/sql"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewNoticeDao(t *testing.T) {
	convey.Convey("NewNoticeDao", t, func(ctx convey.C) {
		var (
			db = &sql.DB{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewNoticeDao(db)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyReplyNotice(t *testing.T) {
	convey.Convey("ReplyNotice", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nts, err := d.Notice.ReplyNotice(c)
			ctx.Convey("Then err should be nil.nts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nts, convey.ShouldNotBeNil)
			})
		})
	})
}
