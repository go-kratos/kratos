package point

import (
	"context"
	"flag"
	"testing"

	"go-common/app/interface/main/account/conf"
	"go-common/app/service/main/point/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	flag.Set("conf", "../../cmd/account-interface-example.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

// go test  -test.v -test.run TestServicePointInfo
func TestServicePointInfo(t *testing.T) {
	Convey("TestServicePointInfo", t, func() {
		res, err := s.PointInfo(context.TODO(), 110000262)
		t.Logf("info: %v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServicePointPage
func TestServicePointPage(t *testing.T) {
	Convey("TestServicePointPage", t, func() {
		res, err := s.PointPage(context.TODO(), &model.ArgRPCPointHistory{Mid: 2089809, PN: 1, PS: 10})
		t.Logf("res: %v", res)
		So(err, ShouldBeNil)
	})
}
