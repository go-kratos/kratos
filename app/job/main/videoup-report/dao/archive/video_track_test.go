package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AddVideoTrack(t *testing.T) {
	Convey("AddVideoTrack", t, func() {
		_, err := d.AddVideoTrack(context.TODO(), 1, "Test_InVideoHis", 0, 0, "", "", "")
		So(err, ShouldBeNil)
	})
}
