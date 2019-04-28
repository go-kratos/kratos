package dao

import (
	"context"
	"testing"

	"go-common/app/infra/canal/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_TiDBPosition(t *testing.T) {
	info := &model.TiDBInfo{
		Name:      "test",
		ClusterID: "1",
		Offset:    2,
		CommitTS:  403845808070328359,
	}
	convey.Convey("add position", t, func(ctx convey.C) {
		err := d.UpdateTiDBPosition(context.Background(), info)
		ctx.So(err, convey.ShouldBeNil)
		ctx.Convey("get position", func(ctx convey.C) {
			gotRes, err := d.TiDBPosition(context.Background(), info.Name)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(gotRes, convey.ShouldResemble, info)
		})
	})
}
