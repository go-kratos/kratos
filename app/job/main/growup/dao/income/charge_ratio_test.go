package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvChargeRatio(t *testing.T) {
	Convey("AvChargeRatio", t, func() {
		_, _, err := d.ArchiveChargeRatio(context.Background(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_UpChargeRatio(t *testing.T) {
	Convey("UpChargeRatio", t, func() {
		_, _, err := d.UpChargeRatio(context.Background(), 0, 2000)
		So(err, ShouldBeNil)
	})
}
