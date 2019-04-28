package service

import (
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/admin/main/reply/conf"
	"go-common/library/log"

	_ "github.com/go-sql-driver/mysql"
)

var svr *Service

func TestMain(m *testing.M) {
	flag.Parse()
	conf.ConfPath = "../cmd/reply-admin-test.toml"
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	svr = New(conf.Conf)
	time.Sleep(time.Millisecond * 300)
	m.Run()
	//flush db and log
	time.Sleep(time.Millisecond * 300)
	os.Exit(0)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}
