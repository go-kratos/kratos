package feed

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Archives(t *testing.T) {
	Convey("Archives", t, func() {
		a, err := d.Archives(context.TODO(), []int64{1}, "")
		So(err, ShouldBeNil)
		Println(a)
	})
}
