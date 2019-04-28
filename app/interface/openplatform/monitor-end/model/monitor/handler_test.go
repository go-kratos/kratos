package monitor

import (
	"context"
	"testing"

	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	m *MonitorHandler
	c = context.Background()
)

// TestNewMonitor .
func TestNewMonitor(t *testing.T) {
	Convey("new monitor", t, func() {
		m = NewMonitor(nil)
		So(m, ShouldNotBeNil)
	})
}

// TestLog .
func TestLog(t *testing.T) {
	var err error
	Convey("info", t, func() {
		m.Info(c, "test", []log.D{}...)
		So(err, ShouldBeNil)
	})
	Convey("warn", t, func() {
		m.Warn(c, "test", []log.D{}...)
		So(err, ShouldBeNil)
	})
	Convey("error", t, func() {
		m.Error(c, "test", []log.D{}...)
		So(err, ShouldBeNil)
	})
}
