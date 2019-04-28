package service

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/history/conf"
)

var s *Service

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/history-service-test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
