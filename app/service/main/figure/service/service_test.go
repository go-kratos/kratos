package service

import (
	"context"
	"flag"
	"testing"

	"go-common/app/service/main/figure/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
	ctx       = context.TODO()
	mid int64 = 7593623
)

func init() {
	flag.Set("conf", "../cmd/figure-service-test.toml")
	err := conf.Init()
	if err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
}

func TestService(t *testing.T) {
	Convey("TestService", t, func() {
		err := svr.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestService_FigureInfo(t *testing.T) {
	Convey("TestService_FigureInfo", t, func() {
		f, err := svr.BatchFigureWithRank(ctx, []int64{mid})
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
	})
}

func TestService_FigureWithRank(t *testing.T) {
	Convey("TestService_FigureWithRank", t, func() {
		f, err := svr.FigureWithRank(ctx, mid)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
	})
}
