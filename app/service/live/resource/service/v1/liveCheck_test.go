package v1

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testLiveCheckService *LiveCheckService
)

func init() {
	dir, _ := filepath.Abs("../../cmd/test.toml")
	flag.Set("conf", dir)
	conf.Init()
	testLiveCheckService = NewLiveCheckService(conf.Conf)
}

func Test_LiveCheck(t *testing.T) {
	Convey("LiveCheck", t, func() {
		_, err := testLiveCheckService.LiveCheck(context.TODO(), &v1.LiveCheckReq{
			Platform: "ios",
			System:   "12.0",
			Mobile:   "iPhone 8",
		})
		So(err, ShouldBeNil)
	})
}

func Test_GetLiveCheckList(t *testing.T) {
	Convey("GetLiveCheckList", t, func() {
		_, err := testLiveCheckService.GetLiveCheckList(context.TODO(), &v1.GetLiveCheckListReq{})
		So(err, ShouldBeNil)
	})
}
func Test_AddLiveCheck(t *testing.T) {
	Convey("AddLiveCheck", t, func() {
		_, err := testLiveCheckService.AddLiveCheck(context.TODO(), &v1.AddLiveCheckReq{
			LiveCheck: `{"android":[],"ios":[{"system":"9.0","mobile":["iPhone 2G","iPhone 3G","iPhone 3GS","iPhone 4","iPhone 4S","iPhone 5","iPhone 5c"]}]}`,
		})
		So(err, ShouldBeNil)
	})
}
