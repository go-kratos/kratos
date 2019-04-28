package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Ups(t *testing.T) {
	Convey("Ups", t, func() {
		_, _, err := d.Ups(context.Background(), "video", 0, 2000)
		So(err, ShouldBeNil)
	})
}
