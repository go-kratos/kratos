package service

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/dm2/conf"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/dm2-test.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	os.Exit(m.Run())
}
