package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Videos2(t *testing.T) {
	Convey("Videos2", t, func() {
		_, err := d.Videos2(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}
