package upcrm

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmQueryUpScoreHistory(t *testing.T) {
	convey.Convey("QueryUpScoreHistory", t, func(ctx convey.C) {
		var (
			mid       = int64(0)
			scoreType = []int{}
			fromdate  = time.Now()
			todate    = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.QueryUpScoreHistory(mid, scoreType, fromdate, todate)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetLatestUpScoreDate(t *testing.T) {
	convey.Convey("GetLatestUpScoreDate", t, func(ctx convey.C) {
		var (
			mid       = int64(0)
			scoreType = int(0)
			todate    = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			date, err := d.GetLatestUpScoreDate(mid, scoreType, todate)
			ctx.Convey("Then err should be nil.date should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, gorm.ErrRecordNotFound)
				ctx.So(date, convey.ShouldNotBeNil)
			})
		})
	})
}
