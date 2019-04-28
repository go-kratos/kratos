package income

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/service/ctrl"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s      *Service
	charge *AvChargeSvr
)

func init() {
	dir, _ := filepath.Abs("../../cmd/growup-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	if s == nil {
		s = New(conf.Conf, ctrl.NewUnboundedExecutor())
	}
	charge = s.avCharge
}

func TestPing(t *testing.T) {
	Convey("Test_Ping", t, func() {
		err := s.Ping(context.Background())
		So(err, ShouldBeNil)
	})
}
