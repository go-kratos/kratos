package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s1 *BannerService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s1 = NewBannerService(conf.Conf)
}

// go test  -test.v -test.run TestGetBannerInfo
func TestGetBannerInfo(t *testing.T) {
	Convey("TestGetInfo", t, func() {
		res, err := s.GetInfo(context.TODO(), &pb.GetInfoReq{
			Platform: "android",
			Build:    53100,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestGetBanner
func TestGetBanner(t *testing.T) {
	Convey("TestGetInfo", t, func() {
		res, err := s1.GetBanner(context.TODO(), &pb.GetBannerReq{
			Platform: "android",
			Build:    53100,
			Type:     "xxx",
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		So(err, ShouldBeNil)
	})
}
