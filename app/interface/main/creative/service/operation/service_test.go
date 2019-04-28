package operation

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/interface/main/creative/conf"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/service"
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
		Reset(func() {
		})
		f(s)
	}
}
