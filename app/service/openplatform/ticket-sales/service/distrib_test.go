package service

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOrder(t *testing.T) {
	Convey("GetOrder", t, func() {
		time.Sleep(time.Second)

		oi, err := svr.GetOrder(ctx, 10001)
		So(err, ShouldBeNil)
		t.Logf("order:%v", oi)
	})
}
