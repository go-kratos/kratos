package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_CheckBlock(t *testing.T) {
	var (
		c = context.TODO()
	)
	mid := int64(304)
	Convey("should return err be nil", t, func() {
		ok, err := s.CheckBlock(c, mid)
		So(err, ShouldBeNil)
		So(ok, ShouldNotBeNil)
	})
}
