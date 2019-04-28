package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_OnlineCount(t *testing.T) {
	convey.Convey("test online num", t, func(ctx convey.C) {
		aid := int64(17592153)
		cid := int64(28727160)
		data, err := d.OnlineCount(context.Background(), aid, cid)
		ctx.So(err, convey.ShouldBeNil)
		ctx.Println(data)
	})
}
