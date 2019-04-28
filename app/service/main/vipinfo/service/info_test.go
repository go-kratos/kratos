package service

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/main/vipinfo/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c = context.TODO()
	s *Service
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestInfo
func TestInfo(t *testing.T) {
	Convey("TestInfo ", t, func() {
		info, err := s.Info(c, 27515795)
		fmt.Println("info:", info)
		fmt.Println("err:", err)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestInfos
func TestInfos(t *testing.T) {
	Convey("TestInfos ", t, func() {
		info, err := s.Infos(c, []int64{27515795, 1})
		fmt.Println("info:", info)
		fmt.Println("err:", err)
		So(err, ShouldBeNil)
	})
}
