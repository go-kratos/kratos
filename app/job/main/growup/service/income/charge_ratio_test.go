package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvChargeRatio(t *testing.T) {
	Convey("AvChargeRatio", t, func() {
		_, err := s.ratio.ArchiveChargeRatio(context.Background(), 10)
		So(err, ShouldBeNil)
	})
}

func Test_UpChargeRatio(t *testing.T) {
	Convey("UpChargeRatio", t, func() {
		_, err := s.ratio.UpChargeRatio(context.Background(), 10)
		So(err, ShouldBeNil)
	})
}
