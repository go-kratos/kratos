package reply

import (
	"context"
	"go-common/library/database/sql"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewConfigDao(t *testing.T) {
	convey.Convey("NewConfigDao", t, func(ctx convey.C) {
		var (
			db = &sql.DB{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewConfigDao(db)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyLoadConfig(t *testing.T) {
	convey.Convey("LoadConfig", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			tp       = int8(0)
			category = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			m, err := d.Config.LoadConfig(c, oid, tp, category)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(m, convey.ShouldBeNil)
			})
		})
	})
}
