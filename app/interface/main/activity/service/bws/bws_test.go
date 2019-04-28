package bws

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/activity/conf"

	"context"
	"go-common/app/interface/main/activity/model/bws"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func WithService(f func(s *Service)) func() {
	return func() {
		dir, _ := filepath.Abs("../../cmd/activity-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if svf == nil {
			svf = New(conf.Conf)
		}
		time.Sleep(2 * time.Second)
		f(svf)
	}
}

func TestService_Binding(t *testing.T) {
	Convey("test binding", t, WithService(func(s *Service) {
		logMid := int64(1)
		p := &bws.ParamBinding{
			Bid: 1,
			Key: "",
		}
		err := s.Binding(context.Background(), logMid, p)
		So(err, ShouldBeNil)
	}))
}
