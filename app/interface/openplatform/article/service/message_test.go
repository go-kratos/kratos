package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_shouldNofify(t *testing.T) {
	Convey("get data", t, func() {
		So(shouldNofify(0), ShouldBeFalse)
		So(shouldNofify(1), ShouldBeTrue)
		So(shouldNofify(5), ShouldBeTrue)
		So(shouldNofify(10), ShouldBeTrue)
		So(shouldNofify(15), ShouldBeFalse)
		So(shouldNofify(50), ShouldBeTrue)
		So(shouldNofify(55), ShouldBeFalse)
		So(shouldNofify(201), ShouldBeFalse)
		So(shouldNofify(300), ShouldBeTrue)
		So(shouldNofify(1010), ShouldBeFalse)
		So(shouldNofify(10000), ShouldBeTrue)
		So(shouldNofify(20000), ShouldBeTrue)
		So(shouldNofify(21000), ShouldBeFalse)
	})
}
