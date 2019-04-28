package operation

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/web-show/conf"
	rsmdl "go-common/app/interface/main/web-show/model/operation"

	. "github.com/smartystreets/goconvey/convey"
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
		arg := &rsmdl.ArgPromote{
			Tp:    "test",
			Count: 1,
			Rank:  1,
		}
		res, err := svf.Promote(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
