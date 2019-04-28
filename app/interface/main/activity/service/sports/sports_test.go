package sports

import (
	"context"
	"flag"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/model/sports"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func init() {
	dir, _ := filepath.Abs("../../cmd/activity-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(err)
	}
	if svf == nil {
		svf = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svf)
	}
}

func TestService_Qq(t *testing.T) {
	Convey("test service qq", t, WithService(func(svf *Service) {
		var (
			params url.Values
		)
		res, err := svf.Qq(context.Background(), params, &sports.ParamQq{Tp: 1})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestService_News(t *testing.T) {
	Convey("test service qq news", t, WithService(func(svf *Service) {
		var (
			params url.Values
		)
		res, err := svf.News(context.Background(), params, &sports.ParamNews{Route: "getQQNewsIndexAndItemsVerify"})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
