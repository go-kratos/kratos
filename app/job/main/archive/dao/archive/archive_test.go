package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Archive(t *testing.T) {
	Convey("Archive", t, func() {
		archive, err := d.Archive(context.TODO(), 1)
		So(err, ShouldBeNil)
		Println(archive)
	})
}
