package data

import (
	"flag"
	"go-common/app/admin/main/up/conf"
	"os"
	"testing"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/up-admin.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	m.Run()
	os.Exit(m.Run())
}
