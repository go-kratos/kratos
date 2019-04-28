package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/search/conf"
	"go-common/app/admin/main/search/http"
	"go-common/app/admin/main/search/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	log.Info("search-admin start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("search-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("search-admin exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
