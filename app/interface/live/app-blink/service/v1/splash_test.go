package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/interface/live/app-blink/conf"
	rspb "go-common/app/service/live/resource/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

var sp *SplashService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	sp = NewSplashService(conf.Conf)
}

// DEPLOY_ENV=uat go test -run TestNewSplashService
func TestNewSplashService(t *testing.T) {
	Convey("TestNewSplashService", t, func() {
		res, err := sp.GetInfo(context.TODO(), &rspb.GetInfoReq{
			Platform: "android",
			Build:    53100,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
