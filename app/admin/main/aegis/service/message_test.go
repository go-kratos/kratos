package service

import (
	"testing"

	"context"
	"github.com/smartystreets/goconvey/convey"
)

func TestServicesendCreateTaskMsg(t *testing.T) {
	var (
		rid           = int64(1)
		flowID        = int64(1)
		dispatchLimit = int64(1)
		bizid         = int64(1)
	)

	convey.Convey("sendCreateTaskMsg", t, func(ctx convey.C) {
		err := s.sendCreateTaskMsg(context.TODO(), rid, flowID, dispatchLimit, bizid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
