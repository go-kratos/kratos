package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDMSCache(t *testing.T) {
	Convey("should return dms and nil", t, func() {
		dms, err := svr.dms(context.TODO(), 1, 1221, 1, 0)
		So(err, ShouldBeNil)
		Convey("dms shoule not be empty", func() {
			So(len(dms), ShouldNotBeEmpty)
		})
	})
}
