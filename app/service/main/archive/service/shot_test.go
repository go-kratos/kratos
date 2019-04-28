package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Videoshot(t *testing.T) {
	Convey("Videoshot", t, func() {
		_, err := s.Videoshot(context.TODO(), 14761597, 24056839)
		So(err, ShouldNotBeNil)
	})
}
