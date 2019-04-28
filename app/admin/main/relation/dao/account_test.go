package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRPCInfos(t *testing.T) {
	convey.Convey("RPCInfos", t, func() {
		infos, err := d.RPCInfos(context.TODO(), []int64{1, 2, 3})
		convey.So(err, convey.ShouldBeNil)
		convey.So(infos, convey.ShouldNotBeNil)
	})
}
