package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTrackVideoup(t *testing.T) {
	Convey("", t, func() {
		err := svr.trackVideoup(context.TODO(), 10114205)
		So(err, ShouldBeNil)
	})
}
