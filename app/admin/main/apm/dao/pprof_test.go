package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInstance(t *testing.T) {
	convey.Convey("Instance", t, func() {
		var appName = "main.app-svr.app-view"
		ins, err := d.Instances(context.Background(), appName)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ins, convey.ShouldNotEqual, nil)

	})
}
