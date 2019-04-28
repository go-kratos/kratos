package service

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestServiceWeightLog(t *testing.T) {
	convey.Convey("WeightLog", t, func(ctx convey.C) {
		result, cnt, err := s.WeightLog(cntx, 49, 1, 4)
		t.Logf("cnt(%d) error(%v)", cnt, err)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}
