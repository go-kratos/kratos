package service

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/captcha/conf"
)

var svr *Service

func TestMain(m *testing.M) {
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/captcha-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
}
