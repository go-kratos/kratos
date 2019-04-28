package mcndao

import (
	"testing"
	"time"

	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaoGetMcnDataSummary(t *testing.T) {
	convey.Convey("GetMcnDataSummary", t, func(ctx convey.C) {
		var (
			selec = "*"
			query = "1=?"
			args  = "1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetMcnDataSummary(selec, query, args)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoGetMcnDataSummaryWithDiff(t *testing.T) {
	convey.Convey("GetMcnDataSummaryWithDiff", t, func(ctx convey.C) {
		var (
			signID       = int64(0)
			dataTYpe     mcnmodel.McnDataType
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetMcnDataSummaryWithDiff(signID, dataTYpe, generateDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoGetDataUpLatestDate(t *testing.T) {
	convey.Convey("GetDataUpLatestDate", t, func(ctx convey.C) {
		var (
			dataType = mcnmodel.DataType(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			generateDate, err := d.GetDataUpLatestDate(dataType, 0)
			ctx.Convey("Then err should be nil.generateDate should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, gorm.ErrRecordNotFound)
				ctx.So(generateDate, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoGetAllUpData(t *testing.T) {
	convey.Convey("GetAllUpData", t, func(ctx convey.C) {
		var (
			signID       = int64(0)
			upmid        = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetAllUpData(signID, upmid, generateDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoGetAllUpDataTemp(t *testing.T) {
	convey.Convey("GetAllUpDataTemp", t, func(ctx convey.C) {
		var (
			signID       = int64(0)
			upmid        = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetAllUpDataTemp(signID, upmid, generateDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
