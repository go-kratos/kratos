package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddLog(t *testing.T) {
	Convey("log", t, func() {
		d.AddLog(1, time.Now().Unix(), 1, 2, "test", "127.0.0.1", "", 0, 0)
	})
}

func TestLogs(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("work", t, func() {
		_, err := d.CoinLog(c, 88888929)
		So(err, ShouldBeNil)
	})
}
