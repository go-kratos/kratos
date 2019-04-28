package service

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/search/conf"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/search-admin-test.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
	os.Exit(m.Run())
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}
