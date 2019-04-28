package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/player/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_BlockTime(t *testing.T) {
	convey.Convey("error should be nill", t, func(ctx convey.C) {
		res, err := d.BlockTime(context.Background(), 88889069)
		ctx.So(err, convey.ShouldBeNil)
		ctx.So(res, convey.ShouldHaveSameTypeAs, &model.BlockTime{})
	})
}
