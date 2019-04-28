package oplog

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_QueryOpLogs(t *testing.T) {
	var (
		c          = context.TODO()
		dmid int64 = 1222
	)
	Convey("QueryOpLogs", t, func() {
		rets, err := dao.QueryOpLogs(c, dmid)
		So(err, ShouldBeNil)
		So(len(rets), ShouldBeGreaterThan, 0)
	})
}
