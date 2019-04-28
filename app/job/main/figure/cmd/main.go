package main

import (
	"flag"
	"os"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/http"
	"go-common/app/job/main/figure/service"
	"go-common/library/log"
	"go-common/library/os/signal"
	"go-common/library/syscall"
)

var (
	svr *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	// service init
	svr = service.New(conf.Conf)
	http.Init(svr)
	log.Info("figure-service start")
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Info("figure-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("figure-service exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
