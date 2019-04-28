package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/point/conf"
	"go-common/app/job/main/point/http"
	"go-common/app/job/main/point/service"
	"go-common/library/log"
)

var (
	svc *service.Service
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
	log.Info("point start")
	// service init
	svc = service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("point get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("point exit")
			if err := svc.Close(); err != nil {
				log.Error("srv close consumer error(%v)", err)
			}
			time.Sleep(2 * time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
