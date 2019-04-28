package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/coupon/conf"
	"go-common/app/job/main/coupon/http"
	"go-common/app/job/main/coupon/service"
	"go-common/library/log"
)

var (
	srv *service.Service
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
	log.Info("coupon start")
	// service init
	srv = service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("coupon get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("coupon exit")
			srv.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
