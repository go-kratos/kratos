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

var s *BannerService

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = NewBannerService(conf.Conf)
}

// DEPLOY_ENV=uat go test -run TestBannerService_GetBlinkBanner
func TestBannerService_GetBlinkBanner(t *testing.T) {
	Convey("TestGetBlinkBanner", t, func() {
		res, err := s.GetBlinkBanner(context.TODO(), &rspb.GetInfoReq{
			Platform: "android",
			Build:    53100,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
