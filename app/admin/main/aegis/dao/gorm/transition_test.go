package gorm

import (
	"go-common/app/admin/main/aegis/model/net"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var tt = &net.Transition{
	ID:          1,
	NetID:       1,
	Trigger:     net.TriggerManual,
	Name:        "first",
	ChName:      "第一次审核",
	Description: "新建变迁",
	UID:         421,
}

func TestDaoTransitionByID(t *testing.T) {
	convey.Convey("TransitionByID", t, func(ctx convey.C) {
		n, err := d.TransitionByID(cntx, tt.ID)
		ctx.Convey("Then err should be nil.n should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(n, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTransitionList(t *testing.T) {
	var (
		pm = &net.ListNetElementParam{
			NetID: 1,
			Ps:    20,
			ID:    []int64{1},
			Name:  "1",
		}
	)
	convey.Convey("TransitionList", t, func(ctx convey.C) {
		result, err := d.TransitionList(cntx, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTransitions(t *testing.T) {
	convey.Convey("TransitionList", t, func(ctx convey.C) {
		_, err := d.Transitions(cntx, []int64{})
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTransitionByUnique(t *testing.T) {
	convey.Convey("TransitionByUnique", t, func(ctx convey.C) {
		_, err := d.TransitionByUnique(cntx, 0, "")
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTransitionIDByNet(t *testing.T) {
	convey.Convey("TransitionIDByNet", t, func(ctx convey.C) {
		a, err := d.TransitionIDByNet(cntx, []int64{1}, true, true)
		t.Logf("a(%+v)", a)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTranByNet(t *testing.T) {
	convey.Convey("TranByNet", t, func(ctx convey.C) {
		_, err := d.TranByNet(cntx, 0, true)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
