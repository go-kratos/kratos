package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStats(t *testing.T) {
	var (
		oid = int64(1)
		typ = int32(1)
		cnt = int32(1)
		c   = context.Background()
	)
	Convey("test stats count", t, WithDao(func(d *Dao) {
		err := d.SendStats(c, typ, oid, cnt)
		So(err, ShouldBeNil)
	}))
}
