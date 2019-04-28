package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/web-feed/conf"
	"go-common/app/interface/main/web-feed/http"
	"go-common/app/interface/main/web-feed/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	log.Info("web-feed start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("web-feed get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("web-feed exit")
			srv.Close()
			time.Sleep(time.Second * 2)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
