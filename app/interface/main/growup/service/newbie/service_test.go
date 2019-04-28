package newbie

import (
	"flag"
	"go-common/app/interface/main/growup/conf"
	"os"
	"testing"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/growup-interface.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	os.Exit(m.Run())
}
