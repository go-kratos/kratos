package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/live/push-live/conf"
	"go-common/app/interface/live/push-live/http"
	"go-common/app/interface/live/push-live/service"
	"go-common/library/log"
	"go-common/library/net/trace"
	"time"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("push-live start")
	// service init
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("push-live get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			log.Info("push-live exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
