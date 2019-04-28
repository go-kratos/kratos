package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_InVideoHis(t *testing.T) {
	Convey("InVideoHis", t, func() {
		_, err := d.InVideoHis(context.TODO(), 1, "Test_InVideoHis", 0, 0, "", "", "")
		So(err, ShouldBeNil)
	})
}
