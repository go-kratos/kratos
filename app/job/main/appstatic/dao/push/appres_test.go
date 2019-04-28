package push

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/app-resource/api/v1"

	"github.com/smartystreets/goconvey/convey"
)

func TestPusherrlog(t *testing.T) {
	var (
		step = ""
		err  error
	)
	convey.Convey("errlog", t, func(ctx convey.C) {
		errlog(step, err)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestPushgrpcClient(t *testing.T) {
	var (
		grpcAddrs []string
		err       error
		p1        v1.AppResourceClient
	)
	convey.Convey("pickAddrs", t, func(ctx convey.C) {
		grpcAddrs, err = d.pickAddrs()
		ctx.Convey("Then err should be nil.grpcAddrs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(grpcAddrs, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("grpcClient", t, func(ctx convey.C) {
		fmt.Println("Call ", grpcAddrs[0])
		p1, err = d.grpcClient(grpcAddrs[0])
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("callRefresh", t, func(ctx convey.C) {
		err = d.CallRefresh(context.Background())
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})

}
