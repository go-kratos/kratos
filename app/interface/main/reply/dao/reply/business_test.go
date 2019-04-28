package reply

import (
	"context"
	"go-common/library/database/sql"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewBusinessDao(t *testing.T) {
	convey.Convey("NewBusinessDao", t, func(ctx convey.C) {
		var (
			db = &sql.DB{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewBusinessDao(db)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyListBusiness(t *testing.T) {
	convey.Convey("ListBusiness", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			business, err := d.Business.ListBusiness(c)
			ctx.Convey("Then err should be nil.business should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(business, convey.ShouldNotBeNil)
			})
		})
	})
}
