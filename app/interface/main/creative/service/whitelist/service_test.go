package whitelist

import (
	"flag"
	"go-common/app/interface/main/creative/conf"
	accmdl "go-common/app/interface/main/creative/model/account"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}
func Test_whitelist(t *testing.T) {
	mf := &accmdl.MyInfo{}
	var (
		uploadinfo map[string]interface{}
		white      int
	)
	Convey("UploadInfoForMainApp", t, WithService(func(s *Service) {
		uploadinfo, white = s.UploadInfoForMainApp(mf, "ios", 111)
		So(uploadinfo, ShouldNotBeNil)
		So(white, ShouldNotBeNil)
	}))
	Convey("UploadInfoForCreator", t, WithService(func(s *Service) {
		uploadinfo = s.UploadInfoForCreator(mf, 111)
		So(s.Creator, ShouldNotBeNil)
	}))
	Convey("Viewinfo", t, WithService(func(s *Service) {
		s.Viewinfo(mf)
		So(uploadinfo, ShouldNotBeNil)
	}))
}
