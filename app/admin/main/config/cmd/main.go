package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/config/conf"
	"go-common/app/admin/main/config/http"
	"go-common/app/admin/main/config/service"
	"go-common/library/log"

	// register config lint
	_ "go-common/app/admin/main/config/pkg/lint/json"
	_ "go-common/app/admin/main/config/pkg/lint/toml"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("config-admin start")
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("config-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svr.Close()
			log.Info("config-admin exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
