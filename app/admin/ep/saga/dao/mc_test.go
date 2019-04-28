package dao

import (
	"context"
	"testing"

	"go-common/app/admin/ep/saga/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("The err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcData(t *testing.T) {
	convey.Convey("test mc data", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "111"

			dataSet = make(map[string]*model.TeamDataResp)
			dataGet = make(map[string]*model.TeamDataResp)
		)

		teamData := &model.TeamDataResp{
			Department: "live",
			Business:   "ios",
			QueryDes:   "description",
			Total:      10,
		}
		dataSet["zhangsan"] = teamData

		ctx.Convey("set data", func(ctx convey.C) {
			err := d.SetData(c, key, dataSet)
			ctx.Convey("set err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("get data", func(ctx convey.C) {
			err := d.GetData(c, key, &dataGet)
			_, ok := dataGet["zhangsan"]
			ctx.Convey("get err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldEqual, true)
				ctx.So(dataGet["zhangsan"].Department, convey.ShouldEqual, "live")
				ctx.So(dataGet["zhangsan"].Total, convey.ShouldEqual, 10)
			})
		})
		ctx.Convey("delete data", func(ctx convey.C) {
			err := d.DeleteData(c, key)
			_, ok := dataGet["zhangsan"]
			ctx.Convey("get err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldEqual, false)
			})
		})
	})
}

func TestMcPipeline(t *testing.T) {
	convey.Convey("test mc Pipeline", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "111"
		)

		pipelineData := &model.PipelineDataResp{
			Department: "openplatform",
			Business:   "android",
			QueryDes:   "description",
			Total:      11,
		}

		ctx.Convey("set pipeline data", func(ctx convey.C) {
			err := d.SetPipeline(c, key, pipelineData)
			ctx.Convey("set err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("get pipeline data", func(ctx convey.C) {
			pipeline, err := d.GetPipeline(c, key)
			ctx.Convey("get err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pipeline.Business, convey.ShouldEqual, "android")
				ctx.So(pipeline.Total, convey.ShouldEqual, 11)
			})
		})
	})
}
