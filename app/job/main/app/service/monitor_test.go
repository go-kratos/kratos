package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWatchSideBar(t *testing.T) {
	convey.Convey("Archive3", t, func(ctx convey.C) {
		s.WatchSideBar()
	})
}
