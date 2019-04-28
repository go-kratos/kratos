package upcrm

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmQueryUpBaseInfo(t *testing.T) {
	convey.Convey("QueryUpBaseInfo", t, func(ctx convey.C) {
		var (
			mid    = int64(0)
			fields = "*"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryUpBaseInfo(mid, fields)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, gorm.ErrRecordNotFound)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryUpBaseInfoBatchByMid(t *testing.T) {
	convey.Convey("QueryUpBaseInfoBatchByMid", t, func(ctx convey.C) {
		var (
			fields = "*"
			mid    = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryUpBaseInfoBatchByMid(fields, mid)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmQueryUpBaseInfoBatchByID(t *testing.T) {
	convey.Convey("QueryUpBaseInfoBatchByID", t, func(ctx convey.C) {
		var (
			fields = "*"
			id     = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryUpBaseInfoBatchByID(fields, id)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
