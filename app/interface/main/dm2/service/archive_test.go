package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVideoDuration(t *testing.T) {
	var (
		aid int64 = 10097265
		oid int64 = 1508
		c         = context.TODO()
	)
	Convey("", t, func() {
		_, err := svr.videoDuration(c, aid, oid)
		So(err, ShouldBeNil)
	})

}
