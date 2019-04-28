package service

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestServiceaddFlowResource(t *testing.T) {
	var (
		tx     = &gorm.DB{}
		rid    = []int64{}
		flowID = int64(0)
		state  = int8(0)
		netid  = int64(0)
	)
	convey.Convey("addFlowResources", t, func(ctx convey.C) {
		err := s.addFlowResources(tx, netid, rid, flowID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicechangeFlowResource(t *testing.T) {
	var (
		tx, _     = s.gorm.BeginTx(cntx)
		netid     = int64(0)
		rid       = int64(0)
		newFlowID = int64(0)
	)
	convey.Convey("changeFlowResource", t, func(ctx convey.C) {
		err := s.updateFlowResources(cntx, tx, netid, rid, newFlowID)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
