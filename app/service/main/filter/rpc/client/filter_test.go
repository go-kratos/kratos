package filter

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	filterconf "go-common/app/service/main/filter/conf"
	rpcmodel "go-common/app/service/main/filter/model/rpc"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx    = context.TODO()
	client *Service
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/filter-service-test.toml")
	var err error
	if err = filterconf.Init(); err != nil {
		panic(err)
	}
	client = New(nil)
	time.Sleep(time.Second)

	os.Exit(m.Run())
}

func TestFilter(t *testing.T) {
	Convey("Filter", t, func() {
		arg := &rpcmodel.ArgFilter{Area: "common", Message: "640"}
		res, err := client.Filter(ctx, arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		arg = &rpcmodel.ArgFilter{Area: "article", Message: "640"}
		res, err = client.Filter(ctx, arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})

	Convey("MFilter", t, func() {
		arg := &rpcmodel.ArgMfilter{
			Area: "common",
			Message: map[string]string{
				"1": "毛泽东",
				"2": "64",
			},
		}
		res, err := client.MFilter(ctx, arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		arg = &rpcmodel.ArgMfilter{
			Area: "article",
			Message: map[string]string{
				"1": "建裆萎业",
				"2": "64",
			},
		}
		res, err = client.MFilter(ctx, arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
