package report

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-view-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(5 * time.Second)
}

func Test_CopyWriter(t *testing.T) {
	Convey("CopyWriter", t, func() {
		s.CopyWriter(context.TODO(), 1, 0, "")
	})
}

func Test_AddReport(t *testing.T) {
	Convey("AddReport", t, func() {
		s.AddReport(context.TODO(), 1684013, 1, 0, "", "", "")
	})
}

func Test_Upload(t *testing.T) {
	Convey("Upload", t, func() {
		s.Upload(context.TODO(), "", []byte{})
	})
}
