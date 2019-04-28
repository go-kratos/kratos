package service

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/job/main/dm2/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/dm2-job.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	Convey("", t, func() {
		err := svr.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}
