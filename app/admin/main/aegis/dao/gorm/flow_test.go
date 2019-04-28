package gorm

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model/net"
)

func TestDaoFlowByID(t *testing.T) {
	convey.Convey("FlowByID", t, func(ctx convey.C) {
		_, err := d.FlowByID(cntx, 1)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFlowList(t *testing.T) {
	var (
		pm = &net.ListNetElementParam{
			NetID: 1,
			State: net.StateAvailable,
			Ps:    5,
			Pn:    2,
			ID:    []int64{1},
			Name:  "name",
		}
	)
	convey.Convey("FlowList", t, func(ctx convey.C) {
		result, err := d.FlowList(cntx, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoFlowByUnique(t *testing.T) {
	convey.Convey("FlowByUnique", t, func(ctx convey.C) {
		_, err := d.FlowByUnique(cntx, 0, "")
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFlowsByNet(t *testing.T) {
	convey.Convey("FlowsByNet", t, func(ctx convey.C) {
		_, err := d.FlowsByNet(cntx, []int64{})
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFlows(t *testing.T) {
	convey.Convey("Flows", t, func(ctx convey.C) {
		_, err := d.Flows(cntx, []int64{})
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoFlowIDByNet(t *testing.T) {
	convey.Convey("FlowIDByNet", t, func(ctx convey.C) {
		_, err := d.FlowIDByNet(cntx, []int64{})
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
