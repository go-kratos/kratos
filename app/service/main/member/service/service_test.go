package service

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/member/conf"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/member-service-example.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	os.Exit(m.Run())
}
