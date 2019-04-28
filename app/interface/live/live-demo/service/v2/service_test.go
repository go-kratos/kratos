package v2

import (
	"context"
	"flag"
	"testing"

	foo "go-common/app/interface/live/live-demo/api/http/v2"
	"go-common/app/interface/live/live-demo/conf"
	"go-common/app/interface/live/live-demo/dao"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *FooService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = NewFooService(conf.Conf)
	dao.InitAPI()
}

// go test  -test.v -test.run TestService_GetInfo
func TestService_GetInfo(t *testing.T) {
	Convey("TestService_GetInfo", t, func() {
		res, err := s.GetInfo(context.TODO(), &foo.GetInfoReq{RoomId: 10002})
		t.Logf("%v msg", res)
		So(err, ShouldBeNil)
	})
}
