package feed

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/app/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_fromAids(t *testing.T) {
	Convey("fromAids", t, func() {
		is := s.fromAids(context.Background(), []int64{1, 2}, time.Now())
		So(is, ShouldNotBeNil)
	})
}
