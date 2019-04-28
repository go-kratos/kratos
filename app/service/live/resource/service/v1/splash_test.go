package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	pb "go-common/app/service/live/resource/api/grpc/v1"
	"go-common/app/service/live/resource/conf"
)

var (
	s *SplashService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = NewSplashService(conf.Conf)
}

// go test  -test.v -test.run TestAdd
func TestGetInfo(t *testing.T) {
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
