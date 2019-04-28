package offer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPushFail(t *testing.T) {
	Convey("PushFail", t, func() {
		err := d.PushFail(ctx(), "")
		So(err, ShouldBeNil)
	})
}
