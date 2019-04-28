package data

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_MonitorNotify(t *testing.T) {
	Convey("MonitorNotify", t, func() {
		_, err := d.MonitorNotify(context.TODO())
		So(err, ShouldBeNil)
	})
}
