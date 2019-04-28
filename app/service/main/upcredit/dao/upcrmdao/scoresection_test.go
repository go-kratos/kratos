package upcrmdao

import (
	"go-common/app/service/main/upcredit/model/calculator"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmdaoInsertScoreSection(t *testing.T) {
	var (
		statis    calculator.OverAllStatistic
		scoreType = int(0)
		date      = time.Now()
	)
	convey.Convey("InsertScoreSection", t, func(ctx convey.C) {
		err := d.InsertScoreSection(statis, scoreType, date)
		err = IgnoreErr(err)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
