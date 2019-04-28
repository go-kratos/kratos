package service

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/service/main/msm/conf"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	flag.Parse()
	dir, _ := filepath.Abs("../cmd/msm-service-example.toml")
	fmt.Println(dir)
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
}
