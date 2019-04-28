package upcrm

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmStartTask(t *testing.T) {
	convey.Convey("StartTask", t, func(ctx convey.C) {
		var (
			taskType = int(0)
			now      = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.StartTask(taskType, now)
			err = nil
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmFinishTask(t *testing.T) {
	convey.Convey("FinishTask", t, func(ctx convey.C) {
		var (
			taskType = int(0)
			now      = time.Now()
			state    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.FinishTask(taskType, now, state)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}
