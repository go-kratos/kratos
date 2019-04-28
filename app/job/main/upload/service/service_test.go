package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/upload/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/bfs-upload-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestPing(t *testing.T) {
	Convey("Ping", t, func() {
		err := s.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestRetryAddRecord(t *testing.T) {

}

func TestRun(t *testing.T) {
	Convey("Run", t, func() {
		Run(s.c)
	})
}

func TestClose(t *testing.T) {
	Convey("Close", t, func() {
		s.Close()
	})
}
