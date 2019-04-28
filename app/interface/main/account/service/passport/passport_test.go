package passport

import (
	"context"
	"flag"
	"testing"

	"go-common/app/interface/main/account/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	flag.Set("conf", "../../cmd/account-interface-example.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

func Test_TestUserName(t *testing.T) {
	Convey("TestUserName", func() {
		err := s.TestUserName(context.TODO(), "testname", 1)
		So(err, ShouldBeNil)
	})
}
