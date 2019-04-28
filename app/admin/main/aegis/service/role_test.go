package service

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestServiceGetRoleBiz(t *testing.T) {
	convey.Convey("GetRoleBiz", t, func(ctx convey.C) {
		result, err := s.GetRoleBiz(cntx, 421, "leader", true)
		t.Logf("result(%+v)", result)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})

}
