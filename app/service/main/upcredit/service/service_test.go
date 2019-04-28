package service

import (
	"flag"
	_ "github.com/davecgh/go-spew/spew"
	"go-common/app/service/main/upcredit/conf"
	"path/filepath"
	"testing"
	"time"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/upcredit-service.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func Test_service(t *testing.T) {

}
