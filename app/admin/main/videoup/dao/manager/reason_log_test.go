package manager

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_AddReasonLog(t *testing.T) {
	convey.Convey("AddReasonLog", t, WithDao(func(d *Dao) {
		_, err := d.AddReasonLog(context.TODO(), 0, 0, 0, 0, 0, 0, time.Now(), time.Now())
		convey.So(err, convey.ShouldBeNil)
	}))
}
