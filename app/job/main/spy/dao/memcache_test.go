package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMCInfo(t *testing.T) {
	Convey("TestMCInfo", t, func() {
		var (
			partition int32 = -2233
			offset1   int64 = 778899
			offset2   int64
			err       error
		)
		err = d.SetOffsetCache(c, "test", partition, offset1)
		So(err, ShouldBeNil)
		offset2, err = d.OffsetCache(c, "test", partition)
		So(offset1, ShouldEqual, offset2)
	})
}
