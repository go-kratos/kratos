package service

import (
	"context"
	"go-common/app/admin/main/aegis/model/net"
	"testing"

	"encoding/json"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
)

func TestServiceShowDirection(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(1)
	)
	convey.Convey("ShowDirection", t, func(ctx convey.C) {
		r, err := s.ShowDirection(c, id)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
		d := r.Direction
		t.Logf("%+v", d)

		oper := &model.NetConfOper{
			OID:    d.ID,
			Action: model.LogNetActionUpdate,
			UID:    d.UID,
			NetID:  d.NetID,
			FlowID: d.FlowID,
			TranID: d.TransitionID,
			Diff: []string{
				model.LogFieldTemp(model.LogFieldDirection, net.DirDirectionDesc[d.Direction], "", false),
				model.LogFieldTemp(model.LogFieldOrder, net.DirOrderDesc[d.Order], "", false),
				model.LogFieldTemp(model.LogFieldGuard, d.Guard, "", false),
				model.LogFieldTemp(model.LogFieldOutput, d.Output, "", false),
			},
		}
		i, _ := json.Marshal(oper)
		t.Logf("%s", i)
		s.sendNetConfLog(c, model.LogTypeDirConf, oper)
		//{"diff":"[指向]为[从节点指向变化]","[顺序]为[]","[输出规则]为[]"}
	})
}

func TestServiceGetDirectionList(t *testing.T) {
	var (
		c  = context.TODO()
		pm = &net.ListDirectionParam{NetID: 1}
	)
	convey.Convey("GetDirectionList", t, func(ctx convey.C) {
		result, err := s.GetDirectionList(c, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestServicecheckDirectionUnique(t *testing.T) {
	var (
		netID        = int64(0)
		FlowID       = int64(0)
		transitionID = int64(0)
		direction    = int8(0)
	)
	convey.Convey("checkDirectionUnique", t, func(ctx convey.C) {
		err, msg := s.checkDirectionUnique(cntx, netID, FlowID, transitionID, direction)
		ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicecheckDirectionBindAvailable(t *testing.T) {
	var (
		flowID       = int64(0)
		transitionID = int64(0)
	)
	convey.Convey("checkDirectionBindAvailable", t, func(ctx convey.C) {
		err := s.checkDirectionBindAvailable(cntx, flowID, transitionID)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceisDirEnable(t *testing.T) {
	var (
		dir = &net.Direction{}
	)
	convey.Convey("isDirEnable", t, func(ctx convey.C) {
		enable := s.isDirEnable(dir)
		ctx.Convey("Then err should be nil.enable should not be nil.", func(ctx convey.C) {
			ctx.So(enable, convey.ShouldNotBeNil)
		})
	})
}

func TestServicegetInDirTransitionID(t *testing.T) {
	var (
		flowID = int64(1)
	)
	convey.Convey("fetchFlowNextEnableDirs", t, func(ctx convey.C) {
		dirs, err := s.fetchFlowNextEnableDirs(cntx, flowID)
		ctx.Convey("Then err should be nil.tids,dirList should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(dirs, convey.ShouldNotBeNil)
		})
	})
}
