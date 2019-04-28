package newcomer

import (
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/service"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/creative.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
	m.Run()
	os.Exit(m.Run())
}
