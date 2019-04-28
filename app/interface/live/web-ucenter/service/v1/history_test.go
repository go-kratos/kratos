package v1

import (
	"context"
	"flag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	api "go-common/app/interface/live/web-ucenter/api/http/v1"
	"go-common/app/interface/live/web-ucenter/conf"
	bm "go-common/library/net/http/blademaster"
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

// go test  -test.v -test.run TestServiceAllowanceList
func TestGetHistoryByUid(t *testing.T) {
	Convey("TestGetHistoryByUid", t, func() {
		res, err := s.GetHistoryByUid(&bm.Context{Context: context.TODO()}, &api.GetHistoryReq{})
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}

func TestDelHistory(t *testing.T) {
	Convey("TestDelHistory", t, func() {
		res, err := s.DelHistory(&bm.Context{Context: context.TODO()}, &api.DelHistoryReq{})
		t.Logf("%v", res)
		So(err, ShouldBeNil)
	})
}
