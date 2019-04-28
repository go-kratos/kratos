package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_AddNetSafeMd5(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("AddNetSafeMd5", t, func(ctx C) {
		_, err := d.AddNetSafeMd5(c, 23333, "ssadasdasdasd")
		So(err, ShouldBeNil)
	})
}
