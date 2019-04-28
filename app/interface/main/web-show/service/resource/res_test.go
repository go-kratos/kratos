package resource

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/web-show/conf"
	rsmdl "go-common/app/interface/main/web-show/model/resource"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_ip = "172.0.0.1"
)

var svf *Service

func init() {
	dir, _ := filepath.Abs("../../cmd/web-show-test.toml")
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

func TestService_Resource(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		arg := &rsmdl.ArgRes{
			Mid:   1,
			ID:    0,
			Pf:    0,
			Sid:   "test",
			IP:    _ip,
			Buvid: "test",
		}
		res, count, err := svf.Resource(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		So(count, ShouldBeGreaterThan, 0)
	}))
}

func TestService_Resources(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		arg := &rsmdl.ArgRess{
			Mid:   1,
			Ids:   []int64{1, 2, 3},
			Pf:    0,
			Sid:   "test",
			IP:    _ip,
			Buvid: "test",
		}
		res, count, err := svf.Resources(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		So(count, ShouldBeGreaterThan, 0)
	}))
}
