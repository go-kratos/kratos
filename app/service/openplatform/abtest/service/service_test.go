package service

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/service/openplatform/abtest/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func init() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	svr = New(conf.Conf)
}

func TestPing(t *testing.T) {
	Convey("TestPing: ", t, func() {
		err := svr.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestEnd(t *testing.T) {
	svr.Close()
}
