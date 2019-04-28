package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_VideoAuditNote(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub string
	)

	Convey("VideoAuditNote", t, func() {
		sub, err = d.VideoAuditNote(c, 2333)
		So(err, ShouldNotBeNil)
		So(sub, ShouldBeEmpty)
	})
}
