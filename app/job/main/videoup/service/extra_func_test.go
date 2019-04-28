package service

import (
	"context"
	"flag"
	"testing"

	"go-common/app/job/main/videoup/conf"

	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/videoup-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func TestArchiveVideo(t *testing.T) {
	Convey("test archive video", t, func() {
		v, a, err := s.archiveVideo(context.Background(), "j171212at3l9zsodcv33n5b7y2ihz0m0")
		So(err, ShouldBeNil)
		So(a, ShouldNotBeNil)
		So(v, ShouldNotBeNil)
		t.Logf("resp: %v", a)
	})

}

func TestService_IsUpperFirstPass(t *testing.T) {
	Convey("", t, func() {
		is, err := s.IsUpperFirstPass(context.Background(), 17515232, 10101454)
		So(err, ShouldBeNil)
		So(is, ShouldNotBeNil)
	})
}
